package envisalink

type ZoneState string

const (
	ZoneStateAlarm         ZoneState = "Alarm"
	ZoneStateAlarmRestore  ZoneState = "Alarm Restore"
	ZoneStateTemper        ZoneState = "Temper"
	ZoneStateTemperRestore ZoneState = "Temper Restore"
	ZoneStateFault         ZoneState = "Fault"
	ZoneStateFaultRestore  ZoneState = "Fault Restore"
	ZoneStateOpen          ZoneState = "Open"
	ZoneStateRestore       ZoneState = "Restore"
)

type Zone struct {
	ID    string
	State ZoneState
}

func newZone(id string) Zone {
	return Zone{
		ID: id,
	}
}
