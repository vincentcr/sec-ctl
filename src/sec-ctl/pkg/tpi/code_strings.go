// Code generated by "stringer -type=ClientCode,ServerCode -output=code_strings.go"; DO NOT EDIT.

package tpi

import "fmt"

const _ClientCode_name = "ClientCodePollClientCodeStatusReportClientCodeNetworkLoginClientCodeDumpZoneTimersClientCodeSetTimeAndDateClientCodeCommandOutputControlClientCodePartitionArmControlAwayClientCodePartitionArmControlStayArmClientCodePartitionArmControlZeroEntryDelayClientCodePartitionArmControlWithCodeClientCodePartitionDisarmControlClientCodeTimeStampControlClientCodeTimeBroadcastControlClientCodeTemperatureBroadcastControlClientCodeTriggerPanicAlarmClientCodeSendKeystrokeStringClientCodeEnterUserCodeProgrammingClientCodeEnterUserProgrammingClientCodeKeepAliveClientCodeCodeSend"

var _ClientCode_map = map[ClientCode]string{
	0:   _ClientCode_name[0:14],
	1:   _ClientCode_name[14:36],
	5:   _ClientCode_name[36:58],
	8:   _ClientCode_name[58:82],
	10:  _ClientCode_name[82:106],
	20:  _ClientCode_name[106:136],
	30:  _ClientCode_name[136:169],
	31:  _ClientCode_name[169:205],
	32:  _ClientCode_name[205:248],
	33:  _ClientCode_name[248:285],
	40:  _ClientCode_name[285:317],
	55:  _ClientCode_name[317:343],
	56:  _ClientCode_name[343:373],
	57:  _ClientCode_name[373:410],
	60:  _ClientCode_name[410:437],
	71:  _ClientCode_name[437:466],
	72:  _ClientCode_name[466:500],
	73:  _ClientCode_name[500:530],
	74:  _ClientCode_name[530:549],
	200: _ClientCode_name[549:567],
}

func (i ClientCode) String() string {
	if str, ok := _ClientCode_map[i]; ok {
		return str
	}
	return fmt.Sprintf("ClientCode(%d)", i)
}

const _ServerCode_name = "ServerCodeAckServerCodeCmdErrServerCodeSysErrServerCodeLoginResServerCodeKeypadLedStateServerCodeKeypadLedFlashStateServerCodeSystemTimeServerCodeRingDetectServerCodeIndoorTemperatureServerCodeOutdoorTemperatureServerCodeZoneAlarmServerCodeZoneAlarmRestoreServerCodeZoneTemperServerCodeZoneTemperRestoreServerCodeZoneFaultServerCodeZoneFaultRestoreServerCodeZoneOpenServerCodeZoneRestoreServerCodeZoneTimerTickServerCodeDuressAlarmServerCodeFireAlarmServerCodeFireAlarmRestoreServerCodeAuxillaryAlarmServerCodeAuxillaryAlarmRestoreServerCodePanicAlarmServerCodePanicAlarmRestoreServerCodeSmokeOrAuxAlarmServerCodeSmokeOrAuxAlarmRestoreServerCodePartitionReadyServerCodePartitionNotReadyServerCodePartitionArmedServerCodePartitionReadyForceArmingEnabledServerCodePartitionInAlarmServerCodePartitionDisarmedServerCodeExitDelayInProgressServerCodeEntryDelayInProgressServerCodeKeypadLockOutServerCodePartitionArmingFailedServerCodePGMOutputInProgressServerCodeChimeEnabledServerCodeChimeDisabledServerCodeInvalidAccessCodeServerCodeFunctionNotAvailableServerCodeArmingFailedServerCodePartitionBusyServerCodeSystemArmingInProgressServerCodeSystemInInstallersModeServerCodeUserClosingServerCodeSpecialClosingServerCodePartialClosingServerCodeUserOpeningServerCodeSpecialOpeningServerCodePanelBatteryTroubleServerCodePanelBatteryTroubleRestoreServerCodePanelACTroubleServerCodePanelACRestoreServerCodeSystemBellTroubleServerCodeSystemBellTroubleRestoralServerCodeFTCTroubleServerCodeBufferNearFullServerCodeGeneralSystemTamperServerCodeGeneralSystemTamperRestoreServerCodeTroubleLEDOnServerCodeTroubleLEDOffServerCodeFireTroubleAlarmServerCodeFireTroubleAlarmRestoreServerCodeVerboseTroubleStatusServerCodeCodeRequiredServerCodeCommandOutputPressedServerCodeMasterCodeRequiredServerCodeInstallersCodeRequired"

var _ServerCode_map = map[ServerCode]string{
	500: _ServerCode_name[0:13],
	501: _ServerCode_name[13:29],
	502: _ServerCode_name[29:45],
	505: _ServerCode_name[45:63],
	510: _ServerCode_name[63:87],
	511: _ServerCode_name[87:116],
	550: _ServerCode_name[116:136],
	560: _ServerCode_name[136:156],
	561: _ServerCode_name[156:183],
	562: _ServerCode_name[183:211],
	601: _ServerCode_name[211:230],
	602: _ServerCode_name[230:256],
	603: _ServerCode_name[256:276],
	604: _ServerCode_name[276:303],
	605: _ServerCode_name[303:322],
	606: _ServerCode_name[322:348],
	609: _ServerCode_name[348:366],
	610: _ServerCode_name[366:387],
	615: _ServerCode_name[387:410],
	620: _ServerCode_name[410:431],
	621: _ServerCode_name[431:450],
	622: _ServerCode_name[450:476],
	623: _ServerCode_name[476:500],
	624: _ServerCode_name[500:531],
	625: _ServerCode_name[531:551],
	626: _ServerCode_name[551:578],
	631: _ServerCode_name[578:603],
	632: _ServerCode_name[603:635],
	650: _ServerCode_name[635:659],
	651: _ServerCode_name[659:686],
	652: _ServerCode_name[686:710],
	653: _ServerCode_name[710:752],
	654: _ServerCode_name[752:778],
	655: _ServerCode_name[778:805],
	656: _ServerCode_name[805:834],
	657: _ServerCode_name[834:864],
	658: _ServerCode_name[864:887],
	659: _ServerCode_name[887:918],
	660: _ServerCode_name[918:947],
	663: _ServerCode_name[947:969],
	664: _ServerCode_name[969:992],
	670: _ServerCode_name[992:1019],
	671: _ServerCode_name[1019:1049],
	672: _ServerCode_name[1049:1071],
	673: _ServerCode_name[1071:1094],
	674: _ServerCode_name[1094:1126],
	680: _ServerCode_name[1126:1158],
	700: _ServerCode_name[1158:1179],
	701: _ServerCode_name[1179:1203],
	702: _ServerCode_name[1203:1227],
	750: _ServerCode_name[1227:1248],
	751: _ServerCode_name[1248:1272],
	800: _ServerCode_name[1272:1301],
	801: _ServerCode_name[1301:1337],
	802: _ServerCode_name[1337:1361],
	803: _ServerCode_name[1361:1385],
	806: _ServerCode_name[1385:1412],
	807: _ServerCode_name[1412:1447],
	814: _ServerCode_name[1447:1467],
	816: _ServerCode_name[1467:1491],
	829: _ServerCode_name[1491:1520],
	830: _ServerCode_name[1520:1556],
	840: _ServerCode_name[1556:1578],
	841: _ServerCode_name[1578:1601],
	842: _ServerCode_name[1601:1627],
	843: _ServerCode_name[1627:1660],
	849: _ServerCode_name[1660:1690],
	900: _ServerCode_name[1690:1712],
	912: _ServerCode_name[1712:1742],
	921: _ServerCode_name[1742:1770],
	922: _ServerCode_name[1770:1802],
}

func (i ServerCode) String() string {
	if str, ok := _ServerCode_map[i]; ok {
		return str
	}
	return fmt.Sprintf("ServerCode(%d)", i)
}
