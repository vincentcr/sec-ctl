package tpi

import (
	"fmt"
)

type PartitionState string

const (
	PartitionStateReady    PartitionState = "Ready"
	PartitionStateNotReady PartitionState = "Not Ready"
	PartitionStateArmed    PartitionState = "Armed"
	PartitionStateInAlarm  PartitionState = "In Alarm"
	PartitionStateDisarmed PartitionState = "Disarmed"
	PartitionStateBusy     PartitionState = "Busy"
)

type Partition struct {
	ID                  string
	State               PartitionState
	TroubleStateLED     bool
	KeypadLEDFlashState KeypadLEDFlashState
	KeypadLEDState      KeypadLEDState
}

func newPartition(id string) *Partition {
	return &Partition{
		ID: id,
	}
}

func (p *Partition) String() string {
	return fmt.Sprintf("Partition{ID:%v, State:%v, TroubleStateLED:%v}", p.ID, p.State, p.TroubleStateLED)
}
