package envisalink

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const keepAliveDelay = 10 * time.Second
const stateRefreshDelay = 300 * time.Second
const maxPendingMessages = 4

//const password = "Q4m1gh" //TODO: read from config

type Client struct {
	loggedIn bool
	password string

	conn           *net.TCPConn
	readCh         chan ServerMessage
	writeCh        chan ClientMessage
	writeCond      *sync.Cond
	msgsPendingAck []ClientCode

	Partitions map[string]Partition
	Zones      map[string]Zone
	EventCh    chan Event
}

func NewClient() *Client {
	return &Client{
		Partitions: map[string]Partition{},
		Zones:      map[string]Zone{},
	}
}

func (c *Client) Connect(hostname string, password string) error {
	c.password = password
	servAddr := hostname + ":4025"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	c.conn = conn
	c.readCh = make(chan ServerMessage)
	c.writeCh = make(chan ClientMessage)
	c.writeCond = sync.NewCond(&sync.Mutex{})
	c.EventCh = make(chan Event)
	c.msgsPendingAck = make([]ClientCode, 0, maxPendingMessages)
	c.startReadLoop()
	c.startWriteLoop()
	c.startProcessingLoop()

	return nil
}

func (c *Client) startWriteLoop() {
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
}

func (c *Client) startReadLoop() {
	go func() {
		for {
			msgs, err := msgReadAvailable(c.conn)
			for _, msg := range msgs {
				c.readCh <- msg
			}
			if err != nil {
				panic(err)
			}
		}
	}()
}

func (c *Client) startProcessingLoop() {
	tickKeepAlive := time.Tick(keepAliveDelay)
	tickStateRefreshDelay := time.Tick(stateRefreshDelay)

	for {
		select {
		case <-tickKeepAlive:
			c.poll()
		case <-tickStateRefreshDelay:
			c.requestStateRefresh()
		case msg := <-c.readCh:
			c.processServerMessage(msg)
		}
	}
}

func (c *Client) poll() {
	if c.loggedIn {
		c.sendMessage(ClientMessage{Code: ClientCodePoll})
	}
}

func (c *Client) requestStateRefresh() {
	if c.loggedIn {
		c.sendMessage(ClientMessage{Code: ClientCodeStatusReport})
	}
}

func (c *Client) sendMessage(msg ClientMessage) {
	c.writeCond.L.Lock()
	for len(c.msgsPendingAck) == maxPendingMessages {
		c.writeCond.Wait()
	}
	c.msgsPendingAck = append(c.msgsPendingAck, msg.Code)
	c.writeCond.L.Unlock()
	c.writeCh <- msg
}

func (c *Client) msgAcked(msg ServerMessage) {

	codeInt, err := msgCodeDecode(msg.Data)
	if err != nil {
		panic(err)
	}

	code := ClientCode(codeInt)

	for i, expectedCode := range c.msgsPendingAck {
		if expectedCode == code {
			c.writeCond.L.Lock()
			c.msgsPendingAck = append(c.msgsPendingAck[:i], c.msgsPendingAck[i+1:]...)
			c.writeCond.L.Unlock()
			c.writeCond.Signal()
			return
		}
	}

	panic(fmt.Errorf("Unexpected ack for %v. pending: %v", msg, c.msgsPendingAck))
}

func (c *Client) processServerMessage(msg ServerMessage) {
	switch msg.Code {
	case ServerCodeLoginRes:
		c.handleLogin(msg)

	case ServerCodeAck:
		c.msgAcked(msg)

	case ServerCodePartitionReady:
		c.updatePartition(string(msg.Data), PartitionStateReady)
	case ServerCodePartitionNotReady:
		c.updatePartition(string(msg.Data), PartitionStateNotReady)
	case ServerCodePartitionArmed:
		c.updatePartition(string(msg.Data), PartitionStateArmed)
	case ServerCodePartitionInAlarm:
		c.updatePartition(string(msg.Data), PartitionStateInAlarm)
	case ServerCodePartitionDisarmed:
		c.updatePartition(string(msg.Data), PartitionStateDisarmed)
	case ServerCodePartitionBusy:
		c.updatePartition(string(msg.Data), PartitionStateBusy)

	case ServerCodeZoneAlarm:
		c.updateZone(string(msg.Data), ZoneStateAlarm)
	case ServerCodeZoneAlarmRestore:
		c.updateZone(string(msg.Data), ZoneStateAlarmRestore)
	case ServerCodeZoneTemper:
		c.updateZone(string(msg.Data), ZoneStateTemper)
	case ServerCodeZoneTemperRestore:
		c.updateZone(string(msg.Data), ZoneStateTemperRestore)
	case ServerCodeZoneFault:
		c.updateZone(string(msg.Data), ZoneStateFault)
	case ServerCodeZoneFaultRestore:
		c.updateZone(string(msg.Data), ZoneStateFaultRestore)
	case ServerCodeZoneOpen:
		c.updateZone(string(msg.Data), ZoneStateOpen)
	case ServerCodeZoneRestore:
		c.updateZone(string(msg.Data), ZoneStateRestore)
	}
}

func (c *Client) handleLogin(msg ServerMessage) {
	if string(msg.Data) == "1" { // login success
		c.loggedIn = true
		c.requestStateRefresh()
	} else {
		loginMsg := ClientMessage{
			Code: ClientCodeNetworkLogin,
			Data: []byte(c.password),
		}
		c.sendMessage(loginMsg)
	}
}

func (c *Client) updatePartition(partitionID string, newState PartitionState) {
	part, ok := c.Partitions[partitionID]
	if !ok {
		part = newPartition(partitionID)
		c.Partitions[partitionID] = part
	}

	if part.State != newState {
		part.State = newState
		c.EventCh <- Event{
			Code:        EventCodePartitionUpdate,
			PartitionID: partitionID,
			EventData:   newState,
		}
	}

}

func (c *Client) updateZone(zoneID string, newState ZoneState) {

	zone, ok := c.Zones[zoneID]
	if !ok {
		zone = newZone(zoneID)
		c.Zones[zoneID] = zone
	}

	if zone.State != newState {
		zone.State = newState
		c.EventCh <- Event{
			Code:      EventCodeZoneUpdate,
			ZoneID:    zoneID,
			EventData: newState,
		}
	}
}

func (c *Client) AwayArm(partID string) {

	msg := ClientMessage{Code: ClientCodePartitionArmControlAway, Data: []byte(partID)}
	c.sendMessage(msg)
}
