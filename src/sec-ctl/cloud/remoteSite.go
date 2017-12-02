package main

import (
	"encoding/json"
	"sec-ctl/cloud/db"
	"sec-ctl/pkg/sites"
	"sec-ctl/pkg/ws"
)

type remoteSite struct {
	id                  db.UUID
	conn                *ws.Conn
	queue               *queue
	partitions          map[string]sites.Partition
	zones               map[string]sites.Zone
	systemTroubleStatus sites.SystemTroubleStatus
	alarms              []sites.Alarm
	eventChs            []chan sites.Event
	stateChangeChs      []chan sites.StateChange
}

func getSiteQueueName(id db.UUID, purpose string) string {
	return "sites:" + id.String() + ":" + purpose
}

func newRemoteSite(site db.Site, conn *ws.Conn, queue *queue) *remoteSite {
	c := &remoteSite{
		id:             site.ID,
		conn:           conn,
		queue:          queue,
		partitions:     map[string]sites.Partition{},
		zones:          map[string]sites.Zone{},
		eventChs:       make([]chan sites.Event, 0),
		stateChangeChs: make([]chan sites.StateChange, 0),
	}

	go func() {
		c.readLoop()
	}()

	go func() {
		c.send(ws.ControlMessage{Code: ws.CtrlGetState})
	}()

	go func() {
		queue.startConsumeLoop(getSiteQueueName(site.ID, "commands"), func(msg qMessage) error {
			cmd := sites.UserCommand{}
			err := json.Unmarshal(msg.data, &cmd)
			if err != nil {
				return err
			}
			err = c.conn.Write(cmd)
			if err != nil {
				c.handleConnErr(err)
			}
			return err
		})
	}()

	return c
}

func (c *remoteSite) readLoop() {
	for {
		i, err := c.conn.Read()
		if err != nil {
			c.handleConnErr(err)
			break
		}

		switch o := i.(type) {
		case sites.SystemState:
			c.processState(o)
		case sites.StateChange:
			c.processStateChange(o)
		case sites.Event:
			c.processEvent(o)
		default:
			logger.Panicf("Unexpected message: %#v", i)
		}

	}
}

func (c *remoteSite) send(obj interface{}) {
	if err := c.conn.Write(obj); err != nil {
		c.handleConnErr(err)
	}

}

func (c *remoteSite) handleConnErr(err error) {
	logger.Println("client disconnected:", err)
	c.queue.publish(queueNameSiteRemoved, []byte(c.id))
}

func (c *remoteSite) Exec(cmd sites.UserCommand) error {

	if err := cmd.Validate(); err != nil {
		return err
	}

	c.send(cmd)

	return nil
}

func (c *remoteSite) processState(st sites.SystemState) {

	for _, p := range st.Partitions {
		c.updatePartition(p)
	}
	for _, z := range st.Zones {
		c.updateZone(z)
	}

	c.alarms = st.Alarms
	c.systemTroubleStatus = st.TroubleStatus
}

func (c *remoteSite) processStateChange(chg sites.StateChange) {
	switch chg.Type {
	case sites.StateChangePartition:
		c.updatePartition(chg.Data.(sites.Partition))
	case sites.StateChangeZone:
		c.updateZone(chg.Data.(sites.Zone))
	case sites.StateChangeSystemTroubleStatus:
		c.updateSystemTroubleStatus(chg.Data.(sites.SystemTroubleStatus))
	default:
		logger.Panicf("Unhandled state change: %v", chg)
	}
}

func (c *remoteSite) updatePartition(part sites.Partition) {
	c.partitions[part.ID] = part
}

func (c *remoteSite) updateZone(zone sites.Zone) {
	c.zones[zone.ID] = zone
}

func (c *remoteSite) updateSystemTroubleStatus(status sites.SystemTroubleStatus) {
	c.systemTroubleStatus = status
}

func (c *remoteSite) processEvent(e sites.Event) {

	data, err := json.Marshal(e)
	if err != nil {
		// this should really work, if it doesn't there is a bug.
		logger.Panicf("Unable to jsonify %#v: %v", e, err)
	}

	c.queue.publish(getSiteQueueName(c.id, "events"), data)
}

func (c *remoteSite) GetState() sites.SystemState {
	return sites.SystemState{
		ID:            c.id.String(),
		Partitions:    c.getPartitions(),
		Zones:         c.getZones(),
		Alarms:        c.alarms,
		TroubleStatus: c.systemTroubleStatus,
	}
}

func (c *remoteSite) getPartitions() []sites.Partition {
	parts := make([]sites.Partition, len(c.partitions))
	idx := 0
	for _, p := range c.partitions {
		parts[idx] = p
		idx++
	}
	return parts
}

func (c *remoteSite) getZones() []sites.Zone {
	zones := make([]sites.Zone, len(c.zones))
	idx := 0
	for _, z := range c.zones {
		zones[idx] = z
		idx++
	}
	return zones
}
