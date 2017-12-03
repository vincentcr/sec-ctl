package sites

import "fmt"

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

func NewZone(id string) *Zone {
	return &Zone{
		ID: id,
	}
}

func (z *Zone) String() string {
	return fmt.Sprintf("Zone{ID:%v, State:%v}", z.ID, z.State)
}
