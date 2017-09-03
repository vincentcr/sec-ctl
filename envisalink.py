#!/usr/bin/env python


import logging
import uuid
import socket
import errno
import time
import threading
import select
from collections import OrderedDict


from Enum import Enum

DefaultPort = 4025
MaxReadRetryCount = 10
CRLF = b'\r\n'
verbose = True


# def main():
#     logging.getLogger().setLevel(logging.DEBUG)
#     tpi = Envisalink(hostname='10.0.0.11', password=b'Q4m1gh')
#     tpi.send_cmd(ClientCodes.StatusReport)
#     while True:
#         time.sleep(5)

#     return


class Envisalink(object):

    def __init__(self, hostname, password):

        self._connect(hostname)

        self._stopped = False
        self._read_th = threading.Thread(target=self.read_msgs_loop)
        self._read_th.daemon = True
        self._read_th.start()
        self._read_lock = threading.RLock()
        self._data_event = threading.Event()
        self._msgs = OrderedDict()
        self._msg_listeners = {}
        self._write_lock = threading.Lock()
        self._keepalive_th = None

        self._authenticate(password)
        self._start_keepalive_loop()

    def _connect(self, hostname):
        address = socket.gethostbyname(hostname)
        self.host = address
        self._socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        print(address, DefaultPort)
        self._socket.connect((address, DefaultPort))
        self._socket.setblocking(0)

    def __del__(self):
        if self._socket:
            self._socket.close()

    def _authenticate(self, password):
        time.sleep(0)
        self._wait_next_msg()
        self.send_cmd(
            ClientCodes.NetworkLogin,
            data=password.encode('utf-8'),
            expected_reply_code=ServerCodes.LoginRes
        )

    def _start_keepalive_loop(self):
        if self._keepalive_th:
            raise Exception("keepalive thread already running")
        self._keepalive_th = threading.Thread(target=self._keepalive_loop)
        self._keepalive_th.daemon = True
        self._keepalive_th.start()

    def _keepalive_loop(self):
        while not self._stopped:
            self.send_cmd(
                ClientCodes.Poll,
                expected_reply_code=ServerCodes.Ack
            )
            time.sleep(5)

    def on(self, code, listener):
        listeners = self._msg_listeners.get(code)
        if not listeners:
            listeners = []
            self._msg_listeners[code] = listeners
        listeners.append(listener)

    def once(self, code, callback):
        def listener(msg):
            self.remove_listener(listener)
            listener(msg)
        self.on(code, callback)

    def remove_listener(self, code, listener):
        listeners = self._msg_listeners.get(code)
        listeners.remove(listener)

    def _consume_msg_loop(self):
        while not self._stopped:
            self._data_event.wait()
            if self._msg_listeners:
                self._try_consume_msgs()

    def _try_consume_msgs(self):
        msgs = self._msgs.items()
        to_rm = []
        for msg in msgs:
            listeners = self._msg_listeners.get(msg.code) or []
            for listener in listeners:
                listener(msg)
                to_rm.append(msg)

        self._remove_msgs(to_rm)

    def _try_consume_msg(self, listener, msg):
        if listener.listener.expected_code \
                    and listener.expected_code == msg.code:
            listener(msg)
            return True

    def _remove_msgs(self, to_rm):
        self._read_lock.acquire()
        try:
            for msg in to_rm:
                del self._msgs[msg.uid]
        finally:
            self._read_lock.release()

    def read_msgs_loop(self):
        while not self._stopped:
            msgs = self._read_msgs()
            logging.debug('_read_msgs_loop: %s', msgs)
            self._push_msgs(msgs)

    def _push_msgs(self, msgs):
        self._read_lock.acquire()
        for msg in msgs:
            self._msgs[msg.uid] = msg
        # logging.debug("_push_msgs")
        self._data_event.set()
        self._read_lock.release()

    def _pop_msg(self):
        self._read_lock.acquire()
        if len(self._msgs) > 0:
            msg = self._msgs.popitem()
            if len(self._msgs) == 0:
                self._data_event.clear()
        else:
            msg = None

        logging.debug("_pop_msg, len = %s: %s" % (len(self._msgs), msg))

        self._read_lock.release()
        return msg

    def _wait_next_msg(self):
        self._data_event.wait()
        return self._pop_msg()

    def _wait_msg(self, code):
        event = threading.Event()
        received_msg = None

        def cb(msg):
            nonlocal received_msg
            received_msg = msg
            event.set()
        self.once(code, cb)
        event.wait()
        return received_msg

    def _read_msgs(self):
        # logging.debug('_read_msgs')

        buf = bytearray()
        eol = False
        while not eol:  # loop until command ends with CRLF
            chunk = self._read_next_chunk()
            if chunk is not None:
                eol = CRLF in chunk
                buf.extend(chunk)

        packets = buf.split(CRLF)
        # logging.debug('read packets: %s', packets)

        msgs = [
            Message.decode(bytes(packet))
            for packet in packets
            if len(packet)
        ]
        return msgs

    def _read_next_chunk(self, retry_count=0):
        try:
            select.select([self._socket], [], [])
            chunk = self._socket.recv(1024)
            chunk_size = len(chunk)
            if chunk_size:
                logging.debug('_read_next_chunk: got %s bytes', chunk_size)
        except socket.error as ex:

            if errno.EAGAIN == ex.errno and retry_count < MaxReadRetryCount:
                time.sleep(1 * retry_count)
                return self._read_next_chunk(retry_count + 1)
            else:
                raise

        return chunk

    def send_cmd(self, code, data=None, expected_reply_code=None):
        msg = Message(code=code, data=data).encode()
        logging.info('send_cmd: code=%s, data=%s, encoded=%s', code, data, msg)
        self._write(msg)
        time.sleep(0)
        if expected_reply_code is not None:
            self.once(code=expected_reply_code, callback=lambda: True)
        reply = self._wait_next_msg()
        return reply

    def _write(self, data):
        # logging.debug('sending: "%s"', [b for b in data])
        self._write_lock.acquire()
        try:
            self._socket.send(data)
        finally:
            self._write_lock.release()


class Message(object):

    def __init__(self, code, data):
        self.code = code
        self.data = data
        self.uid = uuidgen()

    def encode(self):
        cmd = self._encode_cmd()
        checksum = self._compute_checksum(cmd)
        encoded = cmd + checksum + CRLF
        return encoded

    def _encode_cmd(self):
        code_str = b'%03d' % self.code.val
        if self.data is not None:
            data_str = bytes(self.data)
        else:
            data_str = b''

        cmd = code_str + data_str
        return cmd

    def _compute_checksum(self, cmd):
        cmd_chars = [c for c in cmd]
        checksum_int = sum(cmd_chars) & 0xff
        checksum = b'%02x' % checksum_int
        return checksum

    @classmethod
    def decode(cls, packet):
        try:
            data_end = -2
            data_start = 3

            code_int = int(packet[:data_start])
            data = packet[data_start:data_end]
            checksum = packet[data_end:]

            # logging.debug('packet=%s, data=%s, code_int=%s, checksum=%s',
            #               packet, data, code_int, checksum)
            code = ServerCodes.by_val(code_int)
            msg = Message(code=code, data=data)

            msg._verify_checksum(checksum)

            return msg
        except e:
            raise Exception("Failed to decode: {}".format(packet)) from e

    def _verify_checksum(self, checksum_in):
        cmd = self._encode_cmd()
        checksum = self._compute_checksum(cmd)
        assert checksum_in.lower() == checksum.lower(), \
            "Bad checksum for %s: computed: %s; received: %s" % (
                self,
                checksum.lower(),
                checksum_in.lower()
            )

    def __repr__(self):
        return 'Message[code={code};data={data}]'.format(
            code=self.code, data=self.data
            )

    def __str__(self):
        sep_plus_data = ':%s' % self.data if len(self.data) else ''
        return '{code}{data}'.format(
            code=self.code, data=sep_plus_data
            )


def uuidgen():
    return str(uuid.uuid4()).lower().replace('-', '')


ClientCodes = Enum(**{
    'Poll': 0,
    'StatusReport': 1,
    'DumpZoneTimers': 8,
    'NetworkLogin': 5,
    'SetTimeAndDate': 10,
    'CommandOutputControl': 20,
    'PartitionArmControl_Away': 30,
    'PartitionArmControl_StayArm': 31,
    'PartitionArmControl_ZeroEntryDelay': 32,
    'PartitionArmControl_WithCode': 33,
    'PartitionDisarmControl': 40,
    'TimeStampControl': 55,
    'TimeBroadcastControl': 56,
    'TemperatureBroadcastControl': 57,
    'TriggerPanicAlarm': 60,
    # deprecated 'Single Keystroke_Partition1': 70,
    'SendKeystrokeString': 71,
    'EnterUserCodeProgramming': 72,
    'EnterUserProgramming': 73,
    'KeepAlive': 74,
    'CodeSend': 200,
})

'''
     other codes: 229 241 237
'''
ServerCodes = Enum(**{

    'Ack': 500,
    'CmdErr': 501,
    'SysErr': 502,
    'LoginRes': 505,
    'KeypadLedState': 510,
    'KeypadLedFlashState': 511,
    'SystemTime': 550,
    'RingDetect': 560,
    'IndoorTemperature': 561,
    'OutdoorTemperature': 562,

    'ZoneAlarm': 601,
    'ZoneAlarmRestore': 602,
    'ZoneTemper': 603,
    'ZoneTemperRestore': 604,
    'ZoneFault': 605,
    'ZoneFaultRestore': 606,
    'ZoneOpen': 609,
    'ZoneRestore': 610,

    'ZoneTimerTick': 615,

    'DuressAlarm': 620,
    'FireAlarm': 621,
    'FireAlarmRestore': 622,
    'AuxillaryAlarm': 623,
    'AuxillaryAlarmRestore': 624,
    'PanicAlarm': 625,
    'PanicAlarmRestore': 626,
    'SmokeOrAuxAlarm': 631,
    'SmokeOrAuxAlarmRestore': 632,

    'PartitionReady': 650,
    'PartitionNotReady': 651,
    'PartitionArmed': 652,
    'PartitionReady_ForceArmingEnabled': 653,
    'PartitionInAlarm': 654,
    'PartitionDisarmed': 655,
    'ExitDelayInProgress': 656,
    'EntryDelayInProgress': 657,

    'KeypadLockOut': 658,
    'PartitionArmingFailed': 659,
    'PGMOutputInProgress': 660,
    'ChimeEnabled': 663,
    'ChimeDisabled': 664,
    'InvalidAccessCode': 670,
    'FunctionNotAvailable': 671,
    'FunctionNotAvailable': 671,
    'ArmingFailed': 672,
    'PartitionBusy': 673,
    'SystemArmingInProgress': 674,
    'SystemArmingInProgress': 674,
    'SystemInInstallersMode': 680,

    'UserClosing': 700,
    'SpecialClosing': 701,
    'PartialClosing': 702,
    'UserOpening': 750,
    'SpecialOpening': 751,

    'PanelBatteryTrouble': 800,
    'PanelBatteryTroubleRestore': 801,
    'PanelACTrouble': 802,
    'PanelACRestore': 803,
    'SystemBellTrouble': 806,
    'SystemBellTroubleRestoral': 807,

    'FTCTrouble': 814,
    'BufferNearFull': 816,
    'GeneralSystemTamper': 829,
    'GeneralSystemTamperRestore': 830,
    'TroubleLEDOn': 840,
    'TroubleLEDOff': 841,
    'FireTroubleAlarm': 842,
    'FireTroubleAlarmRestore': 843,
    'VerboseTroubleStatus': 849,

    'CodeRequired': 900,
    'CommandOutputPressed': 912,
    'MasterCodeRequired': 921,
    'InstallersCodeRequired': 922,

})

# if __name__ == "__main__":
#     main()
