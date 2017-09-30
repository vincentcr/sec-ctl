package site

import (
	"fmt"
)

// PartitionState represents the state of the partition
type PartitionState string

const (
	// PartitionStateReady represents the "Ready" state of the partition
	PartitionStateReady PartitionState = "Ready"
	// PartitionStateNotReady represents the "Not Ready" state of the partition
	PartitionStateNotReady PartitionState = "Not Ready"
	// PartitionStateArmed represents the "Armed" state of the partition
	PartitionStateArmed PartitionState = "Armed"
	// PartitionStateInAlarm represents the "In Alarm" state of the partition
	PartitionStateInAlarm PartitionState = "In Alarm"
	// PartitionStateDisarmed represents the "Disarmed" state of the partition
	PartitionStateDisarmed PartitionState = "Disarmed"
	// PartitionStateBusy represents the "Busy" state of the partition
	PartitionStateBusy PartitionState = "Busy"
)

// Partition represents a partition in the alarm system
type Partition struct {
	ID                  string
	State               PartitionState
	TroubleStateLED     bool
	KeypadLEDFlashState KeypadLEDFlashState
	KeypadLEDState      KeypadLEDState
}

func NewPartition(id string) *Partition {
	return &Partition{
		ID: id,
	}
}

func (p *Partition) String() string {
	return fmt.Sprintf("Partition{ID:%v, State:%v, TroubleStateLED:%v}", p.ID, p.State, p.TroubleStateLED)
}

// KeypadLEDState represents the state of the keypad LED
type KeypadLEDState byte

const (
	// KeypadLEDStateReady is the "Ready" state of the LED
	KeypadLEDStateReady KeypadLEDState = 1 << iota
	// KeypadLEDStateArmed is the "Armed" state of the LED
	KeypadLEDStateArmed
	// KeypadLEDStateMemory is the "Memory" state of the LED
	KeypadLEDStateMemory
	// KeypadLEDStateBypass is the "Bypass" state of the LED
	KeypadLEDStateBypass
	// KeypadLEDStateTrouble is the "Trouble" state of the LED
	KeypadLEDStateTrouble
	// KeypadLEDStateProgram is the "Program" state of the LED
	KeypadLEDStateProgram
	// KeypadLEDStateFire is the "Fire" state of the LED
	KeypadLEDStateFire
	// KeypadLEDStateBacklight is the "Backlight" state of the LED
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

// KeypadLEDFlashState represents the flash state of the keypad LED
type KeypadLEDFlashState byte

const (
	// KeypadLEDFlashStateReady is the "Ready" state of the LED
	KeypadLEDFlashStateReady KeypadLEDFlashState = 1 << iota
	// KeypadLEDFlashStateArmed is the "Armed" state of the LED
	KeypadLEDFlashStateArmed
	// KeypadLEDFlashStateMemory is the "Memory" state of the LED
	KeypadLEDFlashStateMemory
	// KeypadLEDFlashStateBypass is the "Bypass" state of the LED
	KeypadLEDFlashStateBypass
	// KeypadLEDFlashStateTrouble is the "Trouble" state of the LED
	KeypadLEDFlashStateTrouble
	// KeypadLEDFlashStateProgram is the "Program" state of the LED
	KeypadLEDFlashStateProgram
	// KeypadLEDFlashStateFire is the "Fire" state of the LED
	KeypadLEDFlashStateFire
	// KeypadLEDFlashStateBacklight is the "Backlight" state of the LED
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
