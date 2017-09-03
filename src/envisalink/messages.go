package envisalink

import (
	"fmt"
	"strconv"
	"strings"
)

var crlf = []byte("\r\n")

type ClientCode int

//go:generate stringer -type=ClientCode

const (
	Poll                              ClientCode = 0
	StatusReport                      ClientCode = 1
	DumpZoneTimers                    ClientCode = 8
	NetworkLogin                      ClientCode = 5
	SetTimeAndDate                    ClientCode = 10
	CommandOutputControl              ClientCode = 20
	PartitionArmControlAway           ClientCode = 30
	PartitionArmControlStayArm        ClientCode = 31
	PartitionArmControlZeroEntryDelay ClientCode = 32
	PartitionArmControlWithCode       ClientCode = 33
	PartitionDisarmControl            ClientCode = 40
	TimeStampControl                  ClientCode = 55
	TimeBroadcastControl              ClientCode = 56
	TemperatureBroadcastControl       ClientCode = 57
	TriggerPanicAlarm                 ClientCode = 60
	SendKeystrokeString               ClientCode = 71
	EnterUserCodeProgramming          ClientCode = 72
	EnterUserProgramming              ClientCode = 73
	KeepAlive                         ClientCode = 74
	CodeSend                          ClientCode = 200
)

type ServerCode int

const (
	Ack                              ServerCode = 500
	CmdErr                           ServerCode = 501
	SysErr                           ServerCode = 502
	LoginRes                         ServerCode = 505
	KeypadLedState                   ServerCode = 510
	KeypadLedFlashState              ServerCode = 511
	SystemTime                       ServerCode = 550
	RingDetect                       ServerCode = 560
	IndoorTemperature                ServerCode = 561
	OutdoorTemperature               ServerCode = 562
	ZoneAlarm                        ServerCode = 601
	ZoneAlarmRestore                 ServerCode = 602
	ZoneTemper                       ServerCode = 603
	ZoneTemperRestore                ServerCode = 604
	ZoneFault                        ServerCode = 605
	ZoneFaultRestore                 ServerCode = 606
	ZoneOpen                         ServerCode = 609
	ZoneRestore                      ServerCode = 610
	ZoneTimerTick                    ServerCode = 615
	DuressAlarm                      ServerCode = 620
	FireAlarm                        ServerCode = 621
	FireAlarmRestore                 ServerCode = 622
	AuxillaryAlarm                   ServerCode = 623
	AuxillaryAlarmRestore            ServerCode = 624
	PanicAlarm                       ServerCode = 625
	PanicAlarmRestore                ServerCode = 626
	SmokeOrAuxAlarm                  ServerCode = 631
	SmokeOrAuxAlarmRestore           ServerCode = 632
	PartitionReady                   ServerCode = 650
	PartitionNotReady                ServerCode = 651
	PartitionArmed                   ServerCode = 652
	PartitionReadyForceArmingEnabled ServerCode = 653
	PartitionInAlarm                 ServerCode = 654
	PartitionDisarmed                ServerCode = 655
	ExitDelayInProgress              ServerCode = 656
	EntryDelayInProgress             ServerCode = 657
	KeypadLockOut                    ServerCode = 658
	PartitionArmingFailed            ServerCode = 659
	PGMOutputInProgress              ServerCode = 660
	ChimeEnabled                     ServerCode = 663
	ChimeDisabled                    ServerCode = 664
	InvalidAccessCode                ServerCode = 670
	FunctionNotAvailable             ServerCode = 671
	ArmingFailed                     ServerCode = 672
	PartitionBusy                    ServerCode = 673
	SystemArmingInProgress           ServerCode = 674
	SystemInInstallersMode           ServerCode = 680
	UserClosing                      ServerCode = 700
	SpecialClosing                   ServerCode = 701
	PartialClosing                   ServerCode = 702
	UserOpening                      ServerCode = 750
	SpecialOpening                   ServerCode = 751
	PanelBatteryTrouble              ServerCode = 800
	PanelBatteryTroubleRestore       ServerCode = 801
	PanelACTrouble                   ServerCode = 802
	PanelACRestore                   ServerCode = 803
	SystemBellTrouble                ServerCode = 806
	SystemBellTroubleRestoral        ServerCode = 807
	FTCTrouble                       ServerCode = 814
	BufferNearFull                   ServerCode = 816
	GeneralSystemTamper              ServerCode = 829
	GeneralSystemTamperRestore       ServerCode = 830
	TroubleLEDOn                     ServerCode = 840
	TroubleLEDOff                    ServerCode = 841
	FireTroubleAlarm                 ServerCode = 842
	FireTroubleAlarmRestore          ServerCode = 843
	VerboseTroubleStatus             ServerCode = 849
	CodeRequired                     ServerCode = 900
	CommandOutputPressed             ServerCode = 912
	MasterCodeRequired               ServerCode = 921
	InstallersCodeRequired           ServerCode = 922
)

type ClientMessage struct {
	Code ClientCode
	Data []byte
}

type ServerMessage struct {
	Code ServerCode
	Data []byte
}

/*
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

*/
func msgDecode(msgBytes []byte) (ServerMessage, error) {
	if len(msgBytes) < 5 {
		return ServerMessage{}, fmt.Errorf("Got %d bytes, need at least 5", len(msgBytes))
	}

	dataStart := 3
	dataEnd := len(msgBytes) - 2
	codeBytes := msgBytes[:dataStart]
	data := msgBytes[dataStart:dataEnd]
	expectedChecksum := string(msgBytes[dataEnd:])
	actualChecksum := msgChecksum(msgBytes[:dataEnd])
	if strings.ToLower(expectedChecksum) != strings.ToLower(actualChecksum) {
		return ServerMessage{}, fmt.Errorf("failed to decode message %v: data %v, expected checksum %v, actual %v",
			msgBytes, data, expectedChecksum, actualChecksum)
	}

	codeInt, err := strconv.Atoi(string(codeBytes))
	if err != nil {
		return ServerMessage{}, fmt.Errorf("failed to decode message %v: invalid code %s", msgBytes, codeBytes)
	}

	msg := ServerMessage{
		Code: ServerCode(codeInt),
		Data: data,
	}

	return msg, nil

}

func msgEncode(msg ClientMessage) []byte {
	encoded := []byte(fmt.Sprintf("%03d", msg.Code))

	if msg.Data != nil {
		encoded = append(encoded, msg.Data...)
	}

	checksum := msgChecksum(encoded)

	return append(append(encoded, []byte(checksum)...), crlf...)
}

func msgChecksum(bytes []byte) string {
	var sum int
	for _, b := range bytes {
		sum += int(b)
	}

	sum = sum & 0xff
	checksum := fmt.Sprintf("%02X", sum)
	return checksum
}
