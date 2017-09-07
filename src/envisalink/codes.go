package envisalink

type ClientCode int

//go:generate stringer -type=ClientCode

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
