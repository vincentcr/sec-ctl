package tpi

import (
	"fmt"
	"reflect"
	"time"
)

type EventLevel string

const (
	LevelDebug   EventLevel = "DEBUG"
	LevelInfo    EventLevel = "INFO"
	LevelWarn    EventLevel = "WARN"
	LevelError   EventLevel = "ERROR"
	LevelAlarm   EventLevel = "ALARM"
	LevelTrouble EventLevel = "TROUBLE"
)

//Event represents a TPI event
type Event struct {
	Level       EventLevel
	Code        string
	Time        time.Time
	IsAlarm     bool
	Description string
	PartitionID string
	ZoneID      string
	UserID      string
	Data        map[string]interface{}
}

func newEvent(level EventLevel, code string) *Event {
	evt := &Event{
		Level: level,
		Code:  code,
		Time:  time.Now(),
		Data:  map[string]interface{}{},
	}

	return evt
}

func newServerEvent(level EventLevel, code ServerCode) *Event {
	e := newEvent(level, code.stringHuman())
	desc := serverCodeDescriptions[code]
	e.setDescription(desc)
	return e
}

func (e *Event) setDescription(desc string) *Event {
	e.Description = desc
	return e
}

func (e *Event) setPartition(partID string) *Event {
	e.PartitionID = partID
	return e
}

func (e *Event) setZoneID(zoneID string) *Event {
	e.ZoneID = zoneID
	return e
}

func (e *Event) setUserID(userID string) *Event {
	e.UserID = userID
	return e
}

func (e *Event) setData(k string, v interface{}) *Event {
	e.Data[k] = v
	return e
}

func (e Event) String() string {

	desc := e.Description
	if e.UserID != "" {
		desc += " by User " + e.UserID
	}

	dataDesc := ""

	if e.PartitionID != "" {
		dataDesc += "partitionID: " + e.PartitionID
	}

	if e.ZoneID != "" {
		if len(dataDesc) > 0 {
			dataDesc += "; "
		}
		dataDesc += "zoneID: " + e.ZoneID
	}

	if e.UserID != "" {
		if len(dataDesc) > 0 {
			dataDesc += "; "
		}
		dataDesc += "userID: " + e.UserID
	}

	for k, v := range e.Data {
		if !isZeroOfUnderlyingType(v) {
			if len(dataDesc) > 0 {
				dataDesc += "; "
			}
			dataDesc += fmt.Sprintf("%s: %v", k, v)
		}
	}

	if len(dataDesc) > 0 {
		dataDesc = " [" + dataDesc + "]"
	}

	timeStr := e.Time.Format("2006-01-02 15:04:05")

	return fmt.Sprintf("%v: %v %v %s%v", timeStr, e.Level, e.Code, desc, dataDesc)
}

func isZeroOfUnderlyingType(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}
