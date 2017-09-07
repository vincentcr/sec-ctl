package envisalink

import "fmt"

type EventCode string

const (
	EventCodePartitionUpdate EventCode = "Partition Update"
	EventCodeZoneUpdate      EventCode = "Zone Update"
	EventCodeSystemError     EventCode = "System Error"
)

type Event struct {
	PartitionID string
	ZoneID      string
	Code        EventCode
	EventData   interface{}
}

func (e Event) String() string {
	return fmt.Sprintf("Event{Code: %v, PartitionID: %v, ZoneID: %v, Data: %v}", e.Code, e.PartitionID, e.ZoneID, e.EventData)
}

func logger(eventCh chan Event) {
	for {
		select {
		case e := <-eventCh:

		}
	}
}
go func() {
	for {
		select {
		case msg := <-c.writeCh:
			err := msgWrite(c.conn, msg)
			if err != nil {
				panic(err)
			}
		}
	}
}()
