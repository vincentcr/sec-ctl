package envisalink

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
	ID    string
	State PartitionState
}

func newPartition(id string) Partition {
	return Partition{
		ID: id,
	}
}
