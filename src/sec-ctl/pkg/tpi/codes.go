package tpi

import "fmt"

//go:generate stringer -type=ClientCode,ServerCode -output=code_strings.go

// ClientCode are the client command codes supported by the envisalink
type ClientCode int

const (
	// ClientCodePoll is the code for the Poll command
	ClientCodePoll ClientCode = 0
	// ClientCodeStatusReport is the code for the StatusReport command
	ClientCodeStatusReport ClientCode = 1
	// ClientCodeDumpZoneTimers is the code for the DumpZoneTimers command
	ClientCodeDumpZoneTimers ClientCode = 8
	// ClientCodeNetworkLogin is the code for the NetworkLogin command
	ClientCodeNetworkLogin ClientCode = 5
	// ClientCodeSetTimeAndDate is the code for the SetTimeAndDate command
	ClientCodeSetTimeAndDate ClientCode = 10
	// ClientCodeCommandOutputControl is the code for the CommandOutputControl command
	ClientCodeCommandOutputControl ClientCode = 20
	// ClientCodePartitionArmControlAway is the code for the PartitionArmControlAway command
	ClientCodePartitionArmControlAway ClientCode = 30
	// ClientCodePartitionArmControlStayArm is the code for the PartitionArmControlStayArm command
	ClientCodePartitionArmControlStayArm ClientCode = 31
	// ClientCodePartitionArmControlZeroEntryDelay is the code for the PartitionArmControlZeroEntryDelay command
	ClientCodePartitionArmControlZeroEntryDelay ClientCode = 32
	// ClientCodePartitionArmControlWithCode is the code for the PartitionArmControlWithCode command
	ClientCodePartitionArmControlWithCode ClientCode = 33
	// ClientCodePartitionDisarmControl is the code for the PartitionDisarmControl command
	ClientCodePartitionDisarmControl ClientCode = 40
	// ClientCodeTimeStampControl is the code for the TimeStampControl command
	ClientCodeTimeStampControl ClientCode = 55
	// ClientCodeTimeBroadcastControl is the code for the TimeBroadcastControl command
	ClientCodeTimeBroadcastControl ClientCode = 56
	// ClientCodeTemperatureBroadcastControl is the code for the TemperatureBroadcastControl command
	ClientCodeTemperatureBroadcastControl ClientCode = 57
	// ClientCodeTriggerPanicAlarm is the code for the TriggerPanicAlarm command
	ClientCodeTriggerPanicAlarm ClientCode = 60
	// ClientCodeSendKeystrokeString is the code for the SendKeystrokeString command
	ClientCodeSendKeystrokeString ClientCode = 71
	// ClientCodeEnterUserCodeProgramming is the code for the EnterUserCodeProgramming command
	ClientCodeEnterUserCodeProgramming ClientCode = 72
	// ClientCodeEnterUserProgramming is the code for the EnterUserProgramming command
	ClientCodeEnterUserProgramming ClientCode = 73
	// ClientCodeKeepAlive is the code for the KeepAlive command
	ClientCodeKeepAlive ClientCode = 74
	// ClientCodeCodeSend is the code for the CodeSend command
	ClientCodeCodeSend ClientCode = 200
)

// ServerCode are the event codes sent by the server
type ServerCode int

const (
	// ServerCodeAck is the server code for Ack
	ServerCodeAck ServerCode = 500
	// ServerCodeCmdErr is the server code for CmdErr
	ServerCodeCmdErr ServerCode = 501
	// ServerCodeSysErr is the server code for SysErr
	ServerCodeSysErr ServerCode = 502
	// ServerCodeLoginRes is the server code for LoginRes
	ServerCodeLoginRes ServerCode = 505
	// ServerCodeKeypadLedState is the server code for KeypadLedState
	ServerCodeKeypadLedState ServerCode = 510
	// ServerCodeKeypadLedFlashState is the server code for KeypadLedFlashState
	ServerCodeKeypadLedFlashState ServerCode = 511
	// ServerCodeSystemTime is the server code for SystemTime
	ServerCodeSystemTime ServerCode = 550
	// ServerCodeRingDetect is the server code for RingDetect
	ServerCodeRingDetect ServerCode = 560
	// ServerCodeIndoorTemperature is the server code for IndoorTemperature
	ServerCodeIndoorTemperature ServerCode = 561
	// ServerCodeOutdoorTemperature is the server code for OutdoorTemperature
	ServerCodeOutdoorTemperature ServerCode = 562
	// ServerCodeZoneAlarm is the server code for ZoneAlarm
	ServerCodeZoneAlarm ServerCode = 601
	// ServerCodeZoneAlarmRestore is the server code for ZoneAlarmRestore
	ServerCodeZoneAlarmRestore ServerCode = 602
	// ServerCodeZoneTemper is the server code for ZoneTemper
	ServerCodeZoneTemper ServerCode = 603
	// ServerCodeZoneTemperRestore is the server code for ZoneTemperRestore
	ServerCodeZoneTemperRestore ServerCode = 604
	// ServerCodeZoneFault is the server code for ZoneFault
	ServerCodeZoneFault ServerCode = 605
	// ServerCodeZoneFaultRestore is the server code for ZoneFaultRestore
	ServerCodeZoneFaultRestore ServerCode = 606
	// ServerCodeZoneOpen is the server code for ZoneOpen
	ServerCodeZoneOpen ServerCode = 609
	// ServerCodeZoneRestore is the server code for ZoneRestore
	ServerCodeZoneRestore ServerCode = 610
	// ServerCodeZoneTimerTick is the server code for ZoneTimerTick
	ServerCodeZoneTimerTick ServerCode = 615
	// ServerCodeDuressAlarm is the server code for DuressAlarm
	ServerCodeDuressAlarm ServerCode = 620
	// ServerCodeFireAlarm is the server code for FireAlarm
	ServerCodeFireAlarm ServerCode = 621
	// ServerCodeFireAlarmRestore is the server code for FireAlarmRestore
	ServerCodeFireAlarmRestore ServerCode = 622
	// ServerCodeAuxillaryAlarm is the server code for AuxillaryAlarm
	ServerCodeAuxillaryAlarm ServerCode = 623
	// ServerCodeAuxillaryAlarmRestore is the server code for AuxillaryAlarmRestore
	ServerCodeAuxillaryAlarmRestore ServerCode = 624
	// ServerCodePanicAlarm is the server code for PanicAlarm
	ServerCodePanicAlarm ServerCode = 625
	// ServerCodePanicAlarmRestore is the server code for PanicAlarmRestore
	ServerCodePanicAlarmRestore ServerCode = 626
	// ServerCodeSmokeOrAuxAlarm is the server code for SmokeOrAuxAlarm
	ServerCodeSmokeOrAuxAlarm ServerCode = 631
	// ServerCodeSmokeOrAuxAlarmRestore is the server code for SmokeOrAuxAlarmRestore
	ServerCodeSmokeOrAuxAlarmRestore ServerCode = 632
	// ServerCodePartitionReady is the server code for PartitionReady
	ServerCodePartitionReady ServerCode = 650
	// ServerCodePartitionNotReady is the server code for PartitionNotReady
	ServerCodePartitionNotReady ServerCode = 651
	// ServerCodePartitionArmed is the server code for PartitionArmed
	ServerCodePartitionArmed ServerCode = 652
	// ServerCodePartitionReadyForceArmingEnabled is the server code for PartitionReadyForceArmingEnabled
	ServerCodePartitionReadyForceArmingEnabled ServerCode = 653
	// ServerCodePartitionInAlarm is the server code for PartitionInAlarm
	ServerCodePartitionInAlarm ServerCode = 654
	// ServerCodePartitionDisarmed is the server code for PartitionDisarmed
	ServerCodePartitionDisarmed ServerCode = 655
	// ServerCodeExitDelayInProgress is the server code for ExitDelayInProgress
	ServerCodeExitDelayInProgress ServerCode = 656
	// ServerCodeEntryDelayInProgress is the server code for EntryDelayInProgress
	ServerCodeEntryDelayInProgress ServerCode = 657
	// ServerCodeKeypadLockOut is the server code for KeypadLockOut
	ServerCodeKeypadLockOut ServerCode = 658
	// ServerCodePartitionArmingFailed is the server code for PartitionArmingFailed
	ServerCodePartitionArmingFailed ServerCode = 659
	// ServerCodePGMOutputInProgress is the server code for PGMOutputInProgress
	ServerCodePGMOutputInProgress ServerCode = 660
	// ServerCodeChimeEnabled is the server code for ChimeEnabled
	ServerCodeChimeEnabled ServerCode = 663
	// ServerCodeChimeDisabled is the server code for ChimeDisabled
	ServerCodeChimeDisabled ServerCode = 664
	// ServerCodeInvalidAccessCode is the server code for InvalidAccessCode
	ServerCodeInvalidAccessCode ServerCode = 670
	// ServerCodeFunctionNotAvailable is the server code for FunctionNotAvailable
	ServerCodeFunctionNotAvailable ServerCode = 671
	// ServerCodeArmingFailed is the server code for ArmingFailed
	ServerCodeArmingFailed ServerCode = 672
	// ServerCodePartitionBusy is the server code for PartitionBusy
	ServerCodePartitionBusy ServerCode = 673
	// ServerCodeSystemArmingInProgress is the server code for SystemArmingInProgress
	ServerCodeSystemArmingInProgress ServerCode = 674
	// ServerCodeSystemInInstallersMode is the server code for SystemInInstallersMode
	ServerCodeSystemInInstallersMode ServerCode = 680
	// ServerCodeUserClosing is the server code for UserClosing
	ServerCodeUserClosing ServerCode = 700
	// ServerCodeSpecialClosing is the server code for SpecialClosing
	ServerCodeSpecialClosing ServerCode = 701
	// ServerCodePartialClosing is the server code for PartialClosing
	ServerCodePartialClosing ServerCode = 702
	// ServerCodeUserOpening is the server code for UserOpening
	ServerCodeUserOpening ServerCode = 750
	// ServerCodeSpecialOpening is the server code for SpecialOpening
	ServerCodeSpecialOpening ServerCode = 751
	// ServerCodePanelBatteryTrouble is the server code for PanelBatteryTrouble
	ServerCodePanelBatteryTrouble ServerCode = 800
	// ServerCodePanelBatteryTroubleRestore is the server code for PanelBatteryTroubleRestore
	ServerCodePanelBatteryTroubleRestore ServerCode = 801
	// ServerCodePanelACTrouble is the server code for PanelACTrouble
	ServerCodePanelACTrouble ServerCode = 802
	// ServerCodePanelACRestore is the server code for PanelACRestore
	ServerCodePanelACRestore ServerCode = 803
	// ServerCodeSystemBellTrouble is the server code for SystemBellTrouble
	ServerCodeSystemBellTrouble ServerCode = 806
	// ServerCodeSystemBellTroubleRestoral is the server code for SystemBellTroubleRestoral
	ServerCodeSystemBellTroubleRestoral ServerCode = 807
	// ServerCodeFTCTrouble is the server code for FTCTrouble
	ServerCodeFTCTrouble ServerCode = 814
	// ServerCodeBufferNearFull is the server code for BufferNearFull
	ServerCodeBufferNearFull ServerCode = 816
	// ServerCodeGeneralSystemTamper is the server code for GeneralSystemTamper
	ServerCodeGeneralSystemTamper ServerCode = 829
	// ServerCodeGeneralSystemTamperRestore is the server code for GeneralSystemTamperRestore
	ServerCodeGeneralSystemTamperRestore ServerCode = 830
	// ServerCodeTroubleLEDOn is the server code for TroubleLEDOn
	ServerCodeTroubleLEDOn ServerCode = 840
	// ServerCodeTroubleLEDOff is the server code for TroubleLEDOff
	ServerCodeTroubleLEDOff ServerCode = 841
	// ServerCodeFireTroubleAlarm is the server code for FireTroubleAlarm
	ServerCodeFireTroubleAlarm ServerCode = 842
	// ServerCodeFireTroubleAlarmRestore is the server code for FireTroubleAlarmRestore
	ServerCodeFireTroubleAlarmRestore ServerCode = 843
	// ServerCodeVerboseTroubleStatus is the server code for VerboseTroubleStatus
	ServerCodeVerboseTroubleStatus ServerCode = 849
	// ServerCodeCodeRequired is the server code for CodeRequired
	ServerCodeCodeRequired ServerCode = 900
	// ServerCodeCommandOutputPressed is the server code for CommandOutputPressed
	ServerCodeCommandOutputPressed ServerCode = 912
	// ServerCodeMasterCodeRequired is the server code for MasterCodeRequired
	ServerCodeMasterCodeRequired ServerCode = 921
	// ServerCodeInstallersCodeRequired is the server code for InstallersCodeRequired
	ServerCodeInstallersCodeRequired ServerCode = 922
)

// Name returns the code's name
func (c ServerCode) Name() string {
	s := fmt.Sprintf("%v", c)
	return s[10:] //strip "ServerCode" prefix
}

// GetServerCodeDescription returns a description of the server code
func GetServerCodeDescription(code ServerCode) string {
	return serverCodeDescriptions[code]
}

var serverCodeDescriptions = map[ServerCode]string{
	ServerCodeAck:                              "Command Acknowledge",
	ServerCodeCmdErr:                           "Command Error",
	ServerCodeSysErr:                           "System Error",
	ServerCodeLoginRes:                         "Login Interaction",
	ServerCodeKeypadLedState:                   "Keypad LED State",
	ServerCodeKeypadLedFlashState:              "Keypad LED Flash State",
	ServerCodeSystemTime:                       "System Time Broadcast",
	ServerCodeRingDetect:                       "Ring Detected",
	ServerCodeIndoorTemperature:                "Indoor Temperature",
	ServerCodeOutdoorTemperature:               "Outdoor Temperature",
	ServerCodeZoneAlarm:                        "Zone Alarm",
	ServerCodeZoneAlarmRestore:                 "Zone Alarm Restore",
	ServerCodeZoneTemper:                       "Zone Temper",
	ServerCodeZoneTemperRestore:                "Zone Temper Restore",
	ServerCodeZoneFault:                        "Zone Fault",
	ServerCodeZoneFaultRestore:                 "Zone Fault Restore",
	ServerCodeZoneOpen:                         "Zone Open",
	ServerCodeZoneRestore:                      "Zone Restore",
	ServerCodeZoneTimerTick:                    "Zone Timer Dump",
	ServerCodeDuressAlarm:                      "Duress Alarm",
	ServerCodeFireAlarm:                        "Fire Alarm",
	ServerCodeFireAlarmRestore:                 "Fire Alarm Restore",
	ServerCodeAuxillaryAlarm:                   "Auxiliary Alarm",
	ServerCodeAuxillaryAlarmRestore:            "Auxiliary Alarm Restore",
	ServerCodePanicAlarm:                       "Panic Alarm",
	ServerCodePanicAlarmRestore:                "Panic Alarm Restore",
	ServerCodeSmokeOrAuxAlarm:                  "Smoke / Aux Alarm",
	ServerCodeSmokeOrAuxAlarmRestore:           "Smoke / Aux Alarm Restore",
	ServerCodePartitionReady:                   "Partition Ready",
	ServerCodePartitionNotReady:                "Partition Not Ready",
	ServerCodePartitionReadyForceArmingEnabled: "Partition Ready - Force Arming Enabled",
	ServerCodePartitionArmed:                   "Partition Armed",
	ServerCodePartitionInAlarm:                 "Partition In Alarm",
	ServerCodePartitionDisarmed:                "Partition Disarmed",
	ServerCodeExitDelayInProgress:              "Exit Delay in Progress",
	ServerCodeEntryDelayInProgress:             "Entry Delay in Progress",
	ServerCodeKeypadLockOut:                    "Keypad Lockout",
	ServerCodePartitionArmingFailed:            "Partition Failed to Arm",
	ServerCodePartitionBusy:                    "Partition Busy",
	ServerCodePGMOutputInProgress:              "PGM Output In Progress",
	ServerCodeChimeEnabled:                     "Chime Enabled",
	ServerCodeChimeDisabled:                    "Chime Disabled",
	ServerCodeInvalidAccessCode:                "Invalid Access Code",
	ServerCodeFunctionNotAvailable:             "Function Not Available",
	ServerCodeArmingFailed:                     "Arming Failed",
	ServerCodeSystemArmingInProgress:           "System Arming in Progress",
	ServerCodeSystemInInstallersMode:           "System in Installers Mode",
	ServerCodeUserClosing:                      "User Closing",
	ServerCodeSpecialClosing:                   "Special Closing",
	ServerCodePartialClosing:                   "Partial Closing",
	ServerCodeUserOpening:                      "User Opening",
	ServerCodeSpecialOpening:                   "Special Opening",
	ServerCodePanelBatteryTrouble:              "Panel Battery Trouble",
	ServerCodePanelBatteryTroubleRestore:       "Panel Battery Trouble Restore",
	ServerCodePanelACTrouble:                   "Panel AC Trouble",
	ServerCodePanelACRestore:                   "Panel AC Restore",
	ServerCodeSystemBellTrouble:                "System Bell Trouble",
	ServerCodeSystemBellTroubleRestoral:        "System Bell Trouble Restore",
	ServerCodeFTCTrouble:                       "FTC Trouble",
	ServerCodeBufferNearFull:                   "Buffer Near Full",
	ServerCodeGeneralSystemTamper:              "General System Tamper",
	ServerCodeGeneralSystemTamperRestore:       "General System Tamper Restore",
	ServerCodeTroubleLEDOff:                    "Trouble LED Off",
	ServerCodeTroubleLEDOn:                     "Trouble LED On",
	ServerCodeFireTroubleAlarm:                 "Fire Trouble Alarm",
	ServerCodeFireTroubleAlarmRestore:          "Fire Trouble Alarm Restore",
	ServerCodeVerboseTroubleStatus:             "Verbose Trouble Status",
	ServerCodeCodeRequired:                     "Code Required",
	ServerCodeMasterCodeRequired:               "Master Code Required",
	ServerCodeCommandOutputPressed:             "Command Output Pressed",
	ServerCodeInstallersCodeRequired:           "Installers Code Required",
}

// GetErrorCodeDescription returns a description of the error code
func GetErrorCodeDescription(errCode int) string {
	return errorCodeDescriptions[errCode]
}

var errorCodeDescriptions = map[int]string{
	1:  "Receive Buffer Overrun",
	2:  "Receive Buffer Overflow",
	3:  "Transmit Buffer Overflow",
	10: "Keybus Transmit Buffer Overrun",
	11: "Keybus Transmit Time Timeout",
	12: "Keybus Transmit Mode Timeout",
	13: "Keybus Transmit Keystring Timeout",
	14: "Keybus Interface Not Functioning",
	15: "Keybus Busy",
	16: "Keybus Busy Lockout",
	17: "Keybus Busy Installers Mode",
	18: "Keybus Busy General Busy",
	20: "API Command Syntax Error",
	21: "API Command Partition Error",
	22: "API Command Not Supported",
	23: "API System Not Armed",
	24: "API System Not Ready to Arm",
	25: "API Command Invalid Length",
	26: "API User Code not Required",
	27: "API Invalid Characters in Command",
}

// LoginRes are the code values for login
type LoginRes string

const (
	// LoginResLoginRequest represents the LoginRequest code
	LoginResLoginRequest LoginRes = "3"
	// LoginResTimeout represents the Timeout code
	LoginResTimeout LoginRes = "2"
	// LoginResSuccess represents the Success code
	LoginResSuccess LoginRes = "1"
	// LoginResFailure represents the Failure code
	LoginResFailure LoginRes = "0"
)

// ArmMode are the code values for arm mode
type ArmMode byte

const (
	// ArmModeAway represents the Away code
	ArmModeAway ArmMode = 0
	// ArmModeStay represents the Stay code
	ArmModeStay ArmMode = 1
	// ArmModeZeroEntryAway represents the ZeroEntryAway code
	ArmModeZeroEntryAway ArmMode = 2
	// ArmModeZeroEntryStay represents the ZeroEntryStay code
	ArmModeZeroEntryStay ArmMode = 3
)
