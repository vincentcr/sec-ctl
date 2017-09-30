package site

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

func appendFlagDesc(desc string, flags int, flag int, flagDesc string) string {
	if (flags & flag) != 0 {
		if len(desc) != 0 {
			desc += ", "
		}
		desc += flagDesc
	}
	return desc
}
