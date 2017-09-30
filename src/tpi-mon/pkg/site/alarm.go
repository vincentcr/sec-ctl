package site

import (
	"time"
)

// AlarmType represents the different types of alarms
type AlarmType string

const (
	//AlarmTypePartition represents an alarm on a partition
	AlarmTypePartition AlarmType = "Partition"
	//AlarmTypeDuress represents a duress alarm
	AlarmTypeDuress AlarmType = "Duress"
	//AlarmTypeFire represents a "key" fire alarm
	AlarmTypeFire AlarmType = "Fire"
	//AlarmTypeAux represents a "key" aux alarm
	AlarmTypeAux AlarmType = "Aux"
	//AlarmTypePanic represents a panic alarm
	AlarmTypePanic AlarmType = "Panic"
	//AlarmTypeSmokeOrAux represents a smoke alarm
	AlarmTypeSmokeOrAux AlarmType = "Smoke/Aux"
)

// Alarm represents an alarm event
type Alarm struct {
	AlarmType   AlarmType
	Triggered   time.Time
	Restored    time.Time
	PartitionID string
	ZoneID      string
}
