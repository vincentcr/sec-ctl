package main

import (
	"fmt"
	"io"
	"tpi-mon/pkg/site"
	"tpi-mon/pkg/ws"
)

type remoteSite struct {
	id                  string
	conn                *ws.Conn
	registry            *clientRegistry
	partitions          map[string]site.Partition
	zones               map[string]site.Zone
	systemTroubleStatus site.SystemTroubleStatus
	alarms              []site.Alarm
	eventChs            []chan site.Event
	stateChangeChs      []chan site.StateChange
}

type flowCtrl struct {
	done chan struct{}
	err  chan error
}

func initRemoteSite(conn *ws.Conn, registry *clientRegistry) {
	c := &remoteSite{
		id:             "1",
		conn:           conn,
		registry:       registry,
		partitions:     map[string]site.Partition{},
		zones:          map[string]site.Zone{},
		eventChs:       make([]chan site.Event, 0),
		stateChangeChs: make([]chan site.StateChange, 0),
	}

	go func() {
		c.readLoop()
	}()

	go func() {
		c.send(ws.ControlMessage{Code: ws.CtrlGetState})
	}()

}

func (c *remoteSite) readLoop() {
	for {
		i, err := c.conn.Read()
		if err != nil {
			c.handleConnErr(err)
			break
		}

		switch o := i.(type) {
		case site.SystemState:
			c.processState(o)
		case site.StateChange:
			c.processStateChange(o)
		case site.Event:
			c.processEvent(o)
		default:
			panic(fmt.Errorf("Unexpected message: %#v", i))
		}

	}
}

func (c *remoteSite) send(obj interface{}) {
	if err := c.conn.Write(obj); err != nil {
		c.handleConnErr(err)
	}

}

func (c *remoteSite) handleConnErr(err error) {
	c.registry.removeClient(c)
	if err != io.EOF {
		panic(err) //todo: can do better than this LOL
	}
}

func (c *remoteSite) Exec(cmd site.UserCommand) error {

	if err := cmd.Validate(); err != nil {
		return err
	}

	c.send(cmd)

	return nil
}

func (c *remoteSite) SubscribeToEvents() chan site.Event {
	ch := make(chan site.Event)
	c.eventChs = append(c.eventChs, ch)
	return ch
}

func (c *remoteSite) SubscribeToStateChange() chan site.StateChange {
	ch := make(chan site.StateChange)
	c.stateChangeChs = append(c.stateChangeChs, ch)
	return ch
}

func (c *remoteSite) processState(st site.SystemState) {

	for _, p := range st.Partitions {
		c.updatePartition(p)
	}
	for _, z := range st.Zones {
		c.updateZone(z)
	}

	c.id = st.ID
	c.alarms = st.Alarms
	c.systemTroubleStatus = st.TroubleStatus

	c.registry.addClient(c)
}

func (c *remoteSite) processStateChange(chg site.StateChange) {
	switch chg.Type {
	case site.StateChangePartition:
		c.updatePartition(chg.Data.(site.Partition))
	case site.StateChangeZone:
		c.updateZone(chg.Data.(site.Zone))
	case site.StateChangeSystemTroubleStatus:
		c.updateSystemTroubleStatus(chg.Data.(site.SystemTroubleStatus))
	default:
		panic(fmt.Errorf("Unhandled state change: %v", chg))
	}
}

func (c *remoteSite) updatePartition(part site.Partition) {
	c.partitions[part.ID] = part
}

func (c *remoteSite) updateZone(zone site.Zone) {
	c.zones[zone.ID] = zone
}

func (c *remoteSite) updateSystemTroubleStatus(status site.SystemTroubleStatus) {
	c.systemTroubleStatus = status
}

func (c *remoteSite) processEvent(e site.Event) {
	for _, ch := range c.eventChs {
		ch <- e
	}
}

func (c *remoteSite) GetState() site.SystemState {
	return site.SystemState{
		ID:            c.id,
		Partitions:    c.getPartitions(),
		Zones:         c.getZones(),
		Alarms:        c.alarms,
		TroubleStatus: c.systemTroubleStatus,
	}
}

func (c *remoteSite) getPartitions() []site.Partition {
	parts := make([]site.Partition, len(c.partitions))
	idx := 0
	for _, p := range c.partitions {
		parts[idx] = p
		idx++
	}
	return parts
}

func (c *remoteSite) getZones() []site.Zone {
	zones := make([]site.Zone, len(c.zones))
	idx := 0
	for _, z := range c.zones {
		zones[idx] = z
		idx++
	}
	return zones
}
