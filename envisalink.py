#!/usr/bin/env python 

import os, sys, re, logging
import socket, errno, time, threading, select

from Enum import Enum

DefaultPort = 4025
MaxReadRetryCount = 10
CRLF = '\r\n'
verbose = True

def main():
    logging.getLogger().setLevel(logging.DEBUG)
    tpi =   Envisalink( hostname = '10.0.0.11', password = 'Q4m1gh')
    tpi.send_cmd(ClientCommands.StatusReport)

    while True:
        tpi.send_cmd(ClientCommands.KeepAlive)
        time.sleep(5)

    return

class Envisalink(object):

    def __init__(self, hostname, password):

        self._stopped = False

        self._connect(hostname)

        self._read_th = threading.Thread(target = self.read_msgs_loop)
        self._read_th.daemon = True
        self._read_th.start()
        self._read_lock = threading.RLock()
        self._data_event = threading.Event()
        self._msgs = []

        self._authenticate(password)


    def _connect(self, hostname):
        address = socket.gethostbyname(hostname)
        self.host = address
        self._socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._socket.connect((address, DefaultPort))
        self._socket.setblocking(0)

    def __del__(self):
        if self._socket:
            self._socket.close()

    def _authenticate(self, password):
        time.sleep(0)
        reply = self._wait_next_msg()
        reply = self.send_cmd(ClientCommands.NetworkLogin, password)


    def read_msgs_loop(self):
        while not self._stopped:
            msgs = list(self._read_msgs())
            logging.debug('_read_msgs_loop: %s', msgs)
            self._push_msgs(msgs)

    def _push_msgs(self, msgs):
        self._read_lock.acquire()
        self._msgs.extend(msgs)
        logging.debug("_push_msgs")
        self._data_event.set()
        self._read_lock.release()

    def _pop_msg(self):
        self._read_lock.acquire()
        if len(self._msgs) > 0:
            msg = self._msgs.pop()
            if len(self._msgs) == 0:
                self._data_event.clear()
        else:
            msg = None

        logging.debug("_pop_msg, len = %s: %s" % (len(self._msgs), msg))

        self._read_lock.release()
        return msg

    def _wait_msg(self, code):
        self._data_event.wait()




    def _wait_next_msg(self):
        self._data_event.wait()
        return self._pop_msg()


    def _read_msgs(self):
        logging.debug('_read_msgs')

        bytes = ''
        eol = False
        while not eol: # loop until command ends with CRLF
            chunk = self._read_next_chunk()
            if chunk != None:
                eol = self._has_eol(chunk)
                bytes += chunk

        packets = bytes.split(CRLF)

        for packet in packets:
            if len(packet):
                yield Message.decode(packet)

    def _has_eol(self, chunk):
        return len(chunk) >= 2 and chunk[-2] == '\r' and chunk[-1] == '\n' 

    def _read_next_chunk(self, retry_count = 0):

        try:
            select.select([self._socket], [],[])
            chunk = self._socket.recv(1024)
            chunk_size = len(chunk)
            if chunk_size:
                logging.debug('_read_next_chunk: got %s bytes: %s', chunk_size, [c for c in chunk])
        except socket.error as ex:

            if errno.EAGAIN == ex.errno and retry_count < MaxReadRetryCount:
                time.sleep(1 * retry_count)
                return self._read_next_chunk(retry_count + 1)
            else:
                raise

        return chunk

    def send_cmd(self, code, data = None):
        msg = Message(code = code, data = data).encode()
        logging.info('send_cmd: %s', msg)
        self._write(msg)
        time.sleep(0)
        reply = self._wait_next_msg()
        return reply

    def _write(self, data):
        logging.debug('sending: "%s"', [ b for b in data])
        self._socket.send(data)



class Message(object):

    def __init__(self, code, data):
        self.code = code
        self.data = data


    def encode(self):
        cmd = self._encode_cmd()
        checksum = self._compute_checksum(cmd)
        encoded = cmd + checksum + CRLF
        return encoded

    def _encode_cmd(self):
        code_str = '%03d' % self.code.val
        if self.data is not None:
            data_str = str(self.data)
        else:
            data_str = ''

        cmd = code_str + data_str
        return cmd

    def _compute_checksum(self, cmd):
        cmd_chars = [ ord(c) for c in cmd ]
        checksum_int = sum(cmd_chars) & 0xff
        checksum = '%02x' % checksum_int
        return checksum


    @classmethod
    def decode(cls, packet):
        data_end = -2
        data_start = 3

        code_int = int(packet[:data_start])
        data = packet[data_start:data_end]
        checksum = packet[data_end:]

        #log('packet={packet}, data={data}, code_int={code_int}, checksum={checksum}', packet = packet, data = data, code_int = code_int, checksum=checksum)
        code = ServerCodes.by_val(code_int)
        msg = Message(code = code, data = data)

        msg._verify_checksum(checksum)

        return msg

    def _verify_checksum(self, checksum_in):
        cmd = self._encode_cmd()
        checksum = self._compute_checksum(cmd)
        assert checksum_in.lower() == checksum, "Bad checksum for {msg}: computed: {computed}; received: {received}".format(
            msg = self, computed = checksum, received = checksum_in
            )

    def __repr__(self):
        return 'Message[code={code};data={data}]'.format(
            code = self.code, data = self.data
            )

    def __str__(self):
        sep_plus_data = ':%s' % self.data if len(self.data) else ''
        return '{code}{data}'.format(
            code = self.code, data = sep_plus_data
            )        


ClientCommands = Enum(**{
    'Pool' : 0,
    'StatusReport' : 1,
    'DumpZoneTimers' : 8,
    'NetworkLogin' : 5,
    'SetTimeAndDate' : 10,
    'CommandOutputControl' : 20,
    'PartitionArmControl_Away' : 30,
    'PartitionArmControl_StayArm' : 31,
    'PartitionArmControl_ZeroEntryDelay' : 32,
    'PartitionArmControl_WithCode' : 33,
    'PartitionDisarmControl' : 40,
    'TimeStampControl' : 55,
    'TimeBroadcastControl' : 56,
    'TemperatureBroadcastControl' : 57,
    'TriggerPanicAlarm' : 60,
    # deprecated 'Single Keystroke_Partition1' : 70, 
    'SendKeystrokeString' : 71,
    'EnterUserCodeProgramming' : 72,
    'EnterUserProgramming' : 73,
    'KeepAlive' : 74,
    'CodeSend' : 200,
})

#229 241 237 
ServerCodes = Enum(**{
    
    'Ack' : 500,
    'CmdErr' : 501,
    'SysErr' : 502,
    'LoginRes' : 505,
    'KeypadLedState' : 510,
    'KeypadLedFlashState' : 511,
    'SystemTime' : 550,
    'RingDetect' : 560,
    'IndoorTemperature' : 561,
    'OutdoorTemperature' : 562,

    'ZoneAlarm' : 601,
    'ZoneAlarmRestore' : 602,
    'ZoneTemper' : 603,
    'ZoneTemperRestore' : 604,
    'ZoneFault' : 605,
    'ZoneFaultRestore' : 606,
    'ZoneOpen' : 609,
    'ZoneRestore' : 610,

    'ZoneTimerTick' : 615,

    'DuressAlarm' : 620,
    'FireAlarm' : 621,
    'FireAlarmRestore' : 622,
    'AuxillaryAlarm' : 623,
    'AuxillaryAlarmRestore' : 624,
    'PanicAlarm' : 625,
    'PanicAlarmRestore' : 626,
    'SmokeOrAuxAlarm' : 631,
    'SmokeOrAuxAlarmRestore' : 632,

    'PartitionReady' : 650,
    'PartitionNotReady' : 651,
    'PartitionArmed' : 652,
    'PartitionReady_ForceArmingEnabled' : 653,
    'PartitionInAlarm' : 654,
    'PartitionDisarmed' : 655,
    'ExitDelayInProgress' : 656,
    'EntryDelayInProgress' : 657,

    'KeypadLockOut' : 658,
    'PartitionArmingFailed' : 659,
    'PGMOutputInProgress' : 660,
    'ChimeEnabled' : 663,
    'ChimeDisabled' : 664,
    'InvalidAccessCode' : 670,
    'FunctionNotAvailable' : 671,
    'FunctionNotAvailable' : 671,
    'ArmingFailed' : 672,
    'PartitionBusy' : 673,
    'SystemArmingInProgress' : 674,
    'SystemArmingInProgress' : 674,
    'SystemInInstallersMode' : 680,
    
    'UserClosing' : 700,
    'SpecialClosing' : 701,
    'PartialClosing' : 702,
    'UserOpening' : 750,
    'SpecialOpening' : 751,

    'PanelBatteryTrouble' : 800,
    'PanelBatteryTroubleRestore' : 801,
    'PanelACTrouble' : 802,
    'PanelACRestore' : 803,
    'SystemBellTrouble' : 806,
    'SystemBellTroubleRestoral' : 807,

    'FTCTrouble' : 814,
    'BufferNearFull' : 816,
    'GeneralSystemTamper' : 829,
    'GeneralSystemTamperRestore' : 830,
    'TroubleLEDOn' : 840,
    'TroubleLEDOff' : 841,
    'FireTroubleAlarm' : 842,
    'FireTroubleAlarmRestore' : 843,
    'VerboseTroubleStatus' : 849,

    'CodeRequired' : 900,
    'CommandOutputPressed' : 912,
    'MasterCodeRequired' : 921,
    'InstallersCodeRequired' : 922,
})

main()



