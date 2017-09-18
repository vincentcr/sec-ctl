package tpi

//go:generate stringer -type=ClientCode,ServerCode

// ClientCode are the client command codes supported by the envisalink
type ClientCode int

const (
	ClientCodePoll                              ClientCode = 0
	ClientCodeStatusReport                      ClientCode = 1
	ClientCodeDumpZoneTimers                    ClientCode = 8
	ClientCodeNetworkLogin                      ClientCode = 5
	ClientCodeSetTimeAndDate                    ClientCode = 10
	ClientCodeCommandOutputControl              ClientCode = 20
	ClientCodePartitionArmControlAway           ClientCode = 30
	ClientCodePartitionArmControlStayArm        ClientCode = 31
	ClientCodePartitionArmControlZeroEntryDelay ClientCode = 32
	ClientCodePartitionArmControlWithCode       ClientCode = 33
	ClientCodePartitionDisarmControl            ClientCode = 40
	ClientCodeTimeStampControl                  ClientCode = 55
	ClientCodeTimeBroadcastControl              ClientCode = 56
	ClientCodeTemperatureBroadcastControl       ClientCode = 57
	ClientCodeTriggerPanicAlarm                 ClientCode = 60
	ClientCodeSendKeystrokeString               ClientCode = 71
	ClientCodeEnterUserCodeProgramming          ClientCode = 72
	ClientCodeEnterUserProgramming              ClientCode = 73
	ClientCodeKeepAlive                         ClientCode = 74
	ClientCodeCodeSend                          ClientCode = 200
)

// ServerCode are the event codes sent by the server
type ServerCode int

const (
	ServerCodeAck                              ServerCode = 500
	ServerCodeCmdErr                           ServerCode = 501
	ServerCodeSysErr                           ServerCode = 502
	ServerCodeLoginRes                         ServerCode = 505
	ServerCodeKeypadLedState                   ServerCode = 510
	ServerCodeKeypadLedFlashState              ServerCode = 511
	ServerCodeSystemTime                       ServerCode = 550
	ServerCodeRingDetect                       ServerCode = 560
	ServerCodeIndoorTemperature                ServerCode = 561
	ServerCodeOutdoorTemperature               ServerCode = 562
	ServerCodeZoneAlarm                        ServerCode = 601
	ServerCodeZoneAlarmRestore                 ServerCode = 602
	ServerCodeZoneTemper                       ServerCode = 603
	ServerCodeZoneTemperRestore                ServerCode = 604
	ServerCodeZoneFault                        ServerCode = 605
	ServerCodeZoneFaultRestore                 ServerCode = 606
	ServerCodeZoneOpen                         ServerCode = 609
	ServerCodeZoneRestore                      ServerCode = 610
	ServerCodeZoneTimerTick                    ServerCode = 615
	ServerCodeDuressAlarm                      ServerCode = 620
	ServerCodeFireAlarm                        ServerCode = 621
	ServerCodeFireAlarmRestore                 ServerCode = 622
	ServerCodeAuxillaryAlarm                   ServerCode = 623
	ServerCodeAuxillaryAlarmRestore            ServerCode = 624
	ServerCodePanicAlarm                       ServerCode = 625
	ServerCodePanicAlarmRestore                ServerCode = 626
	ServerCodeSmokeOrAuxAlarm                  ServerCode = 631
	ServerCodeSmokeOrAuxAlarmRestore           ServerCode = 632
	ServerCodePartitionReady                   ServerCode = 650
	ServerCodePartitionNotReady                ServerCode = 651
	ServerCodePartitionArmed                   ServerCode = 652
	ServerCodePartitionReadyForceArmingEnabled ServerCode = 653
	ServerCodePartitionInAlarm                 ServerCode = 654
	ServerCodePartitionDisarmed                ServerCode = 655
	ServerCodeExitDelayInProgress              ServerCode = 656
	ServerCodeEntryDelayInProgress             ServerCode = 657
	ServerCodeKeypadLockOut                    ServerCode = 658
	ServerCodePartitionArmingFailed            ServerCode = 659
	ServerCodePGMOutputInProgress              ServerCode = 660
	ServerCodeChimeEnabled                     ServerCode = 663
	ServerCodeChimeDisabled                    ServerCode = 664
	ServerCodeInvalidAccessCode                ServerCode = 670
	ServerCodeFunctionNotAvailable             ServerCode = 671
	ServerCodeArmingFailed                     ServerCode = 672
	ServerCodePartitionBusy                    ServerCode = 673
	ServerCodeSystemArmingInProgress           ServerCode = 674
	ServerCodeSystemInInstallersMode           ServerCode = 680
	ServerCodeUserClosing                      ServerCode = 700
	ServerCodeSpecialClosing                   ServerCode = 701
	ServerCodePartialClosing                   ServerCode = 702
	ServerCodeUserOpening                      ServerCode = 750
	ServerCodeSpecialOpening                   ServerCode = 751
	ServerCodePanelBatteryTrouble              ServerCode = 800
	ServerCodePanelBatteryTroubleRestore       ServerCode = 801
	ServerCodePanelACTrouble                   ServerCode = 802
	ServerCodePanelACRestore                   ServerCode = 803
	ServerCodeSystemBellTrouble                ServerCode = 806
	ServerCodeSystemBellTroubleRestoral        ServerCode = 807
	ServerCodeFTCTrouble                       ServerCode = 814
	ServerCodeBufferNearFull                   ServerCode = 816
	ServerCodeGeneralSystemTamper              ServerCode = 829
	ServerCodeGeneralSystemTamperRestore       ServerCode = 830
	ServerCodeTroubleLEDOn                     ServerCode = 840
	ServerCodeTroubleLEDOff                    ServerCode = 841
	ServerCodeFireTroubleAlarm                 ServerCode = 842
	ServerCodeFireTroubleAlarmRestore          ServerCode = 843
	ServerCodeVerboseTroubleStatus             ServerCode = 849
	ServerCodeCodeRequired                     ServerCode = 900
	ServerCodeCommandOutputPressed             ServerCode = 912
	ServerCodeMasterCodeRequired               ServerCode = 921
	ServerCodeInstallersCodeRequired           ServerCode = 922
)

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

func (c ServerCode) stringHuman() string {
	return c.String()[10:] //strip "ServerCode" prefix
}

type KeypadLEDState byte

const (
	KeypadLEDStateReady KeypadLEDState = 1 << iota
	KeypadLEDStateArmed
	KeypadLEDStateMemory
	KeypadLEDStateBypass
	KeypadLEDStateTrouble
	KeypadLEDStateProgram
	KeypadLEDStateFire
	KeypadLEDStateBacklight
)

func (s KeypadLEDState) String() string {
	desc := ""
	desc = appendFlagDesc(desc, int(s), int(KeypadLEDFlashStateReady), "Ready")
	desc = appendFlagDesc(desc, int(s), int(KeypadLEDFlashStateArmed), "Armed")
	desc = appendFlagDesc(desc, int(s), int(KeypadLEDFlashStateMemory), "Memory")
	desc = appendFlagDesc(desc, int(s), int(KeypadLEDFlashStateBypass), "Bypass")
	desc = appendFlagDesc(desc, int(s), int(KeypadLEDFlashStateProgram), "Program")
	desc = appendFlagDesc(desc, int(s), int(KeypadLEDFlashStateFire), "Fire")
	desc = appendFlagDesc(desc, int(s), int(KeypadLEDFlashStateBacklight), "Backlight")
	return desc
}

func appendFlagDesc(desc string, flags int, flag int, flagDesc string) string {
	if (flags & flag) != 0 {
		if len(desc) != 0 {
			desc += ", "
		}
		desc += flagDesc
	}
	return desc
}

type KeypadLEDFlashState byte

const (
	KeypadLEDFlashStateReady KeypadLEDFlashState = 1 << iota
	KeypadLEDFlashStateArmed
	KeypadLEDFlashStateMemory
	KeypadLEDFlashStateBypass
	KeypadLEDFlashStateTrouble
	KeypadLEDFlashStateProgram
	KeypadLEDFlashStateFire
	KeypadLEDFlashStateBacklight
)

func (s KeypadLEDFlashState) String() string {
	desc := ""
	desc = appendFlagDesc(desc, int(s), int(KeypadLEDStateReady), "Ready")
	desc = appendFlagDesc(desc, int(s), int(KeypadLEDStateArmed), "Armed")
	desc = appendFlagDesc(desc, int(s), int(KeypadLEDStateMemory), "Memory")
	desc = appendFlagDesc(desc, int(s), int(KeypadLEDStateBypass), "Bypass")
	desc = appendFlagDesc(desc, int(s), int(KeypadLEDStateProgram), "Program")
	desc = appendFlagDesc(desc, int(s), int(KeypadLEDStateFire), "Fire")
	desc = appendFlagDesc(desc, int(s), int(KeypadLEDStateBacklight), "Backlight")
	return desc
}

// SystemTroubleStatus is a flag set of system trouble codes
type SystemTroubleStatus byte

const (
	// SystemTroubleStatusServiceRequired indicates Service is required
	SystemTroubleStatusServiceRequired SystemTroubleStatus = 1 << iota
	// SystemTroubleStatusACPowerLost indicates AC Power was lost
	SystemTroubleStatusACPowerLost
	// SystemTroubleStatusTelephoneLineFault indicates a telephone line fault
	SystemTroubleStatusTelephoneLineFault
	// SystemTroubleStatusFailureToCommunicate indicates a failure to communicate
	SystemTroubleStatusFailureToCommunicate
	// SystemTroubleStatusSensorOrZoneFault indicates a sensor or zone fault
	SystemTroubleStatusSensorOrZoneFault
	// SystemTroubleStatusSensorOrZoneTemper indicates a sensor or zone temper
	SystemTroubleStatusSensorOrZoneTemper
	// SystemTroubleStatusSensorOrZoneLowBattery indicates a sensor or zone is on low battery
	SystemTroubleStatusSensorOrZoneLowBattery
	// SystemTroubleStatusLossOfTime indicates a loss of time
	SystemTroubleStatusLossOfTime
)

func (s SystemTroubleStatus) String() string {
	desc := ""
	desc = appendFlagDesc(desc, int(s), int(SystemTroubleStatusServiceRequired), "Service Required")
	desc = appendFlagDesc(desc, int(s), int(SystemTroubleStatusACPowerLost), "AC Power Lost")
	desc = appendFlagDesc(desc, int(s), int(SystemTroubleStatusTelephoneLineFault), "Telephone Line Fault")
	desc = appendFlagDesc(desc, int(s), int(SystemTroubleStatusFailureToCommunicate), "Failure to Communicate")
	desc = appendFlagDesc(desc, int(s), int(SystemTroubleStatusSensorOrZoneFault), "Sensor/Zone Fault")
	desc = appendFlagDesc(desc, int(s), int(SystemTroubleStatusSensorOrZoneTemper), "Sensor/Zone Temper")
	desc = appendFlagDesc(desc, int(s), int(SystemTroubleStatusSensorOrZoneLowBattery), "Sensor/Zone Low Battery")
	desc = appendFlagDesc(desc, int(s), int(SystemTroubleStatusLossOfTime), "Loss Of Time")
	return desc
}

var tpiErrors = map[int]string{
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

type LoginRes string

const (
	LoginResLoginRequest LoginRes = "3"
	LoginResTimeout      LoginRes = "2"
	LoginResSuccess      LoginRes = "1"
	LoginResFailure      LoginRes = "0"
)

type ArmMode byte

const (
	ArmModeAway          ArmMode = 0
	ArmModeStay          ArmMode = 1
	ArmModeZeroEntryAway ArmMode = 2
	ArmModeZeroEntryStay ArmMode = 3
)
