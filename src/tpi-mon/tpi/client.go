package tpi

import (
	"encoding/hex"
	"fmt"
	"net"
	"sync"
	"time"
	"tpi-mon/tpi/proto"
)

const keepAliveDelay = 30 * time.Second
const stateRefreshDelay = 300 * time.Second
const maxPendingMessages = 4

type Client interface {
	GetEventCh() chan *Event
	GetErrorCh() chan error
	GetSystemTroubleStatus() SystemTroubleStatus
	GetZones() []*Zone
	GetPartitions() []*Partition

	AwayArm(partID string)
	StayArm(partID string)
	ZeroEntryDelayArm(partID string)
	ArmWithCode(partID string, pin string)
	Disarm(partID string, pin string)
	PanicAlarm(target string)
}

type LocalClient struct {
	loggedIn bool
	password string

	conn                *net.TCPConn
	readCh              chan ServerMessage
	writeCh             chan ClientMessage
	writeCond           *sync.Cond
	msgsPendingAck      []ClientCode
	partitions          map[string]*Partition
	zones               map[string]*Zone
	eventCh             chan *Event
	errorCh             chan error
	systemTroubleStatus SystemTroubleStatus
}

func NewLocalClient(hostname string, port uint16, password string) (Client, error) {
	c := &LocalClient{
		partitions:     map[string]*Partition{},
		zones:          map[string]*Zone{},
		readCh:         make(chan ServerMessage),
		writeCh:        make(chan ClientMessage),
		writeCond:      sync.NewCond(&sync.Mutex{}),
		eventCh:        make(chan *Event, 512),
		errorCh:        make(chan error),
		msgsPendingAck: make([]ClientCode, 0, maxPendingMessages),
	}
	if err := c.connect(hostname, port, password); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *LocalClient) connect(hostname string, port uint16, password string) error {
	c.password = password
	servAddr := fmt.Sprintf("%s:%d", hostname, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	c.conn = conn
	c.startReadLoop()
	c.startWriteLoop()
	c.startProcessingLoop()

	return nil
}

func (c *LocalClient) GetEventCh() chan *Event {
	return c.eventCh
}

func (c *LocalClient) GetErrorCh() chan error {
	return c.errorCh
}

func (c *LocalClient) GetPartitions() []*Partition {
	parts := make([]*Partition, len(c.partitions))
	idx := 0
	for _, p := range c.partitions {
		parts[idx] = p
		idx++
	}
	return parts
}

func (c *LocalClient) GetZones() []*Zone {
	zones := make([]*Zone, len(c.zones))
	idx := 0
	for _, p := range c.zones {
		zones[idx] = p
		idx++
	}
	return zones
}
func (c *LocalClient) GetSystemTroubleStatus() SystemTroubleStatus {
	return c.systemTroubleStatus
}

func (c *LocalClient) AwayArm(partID string) {
	msg := ClientMessage{Code: ClientCodePartitionArmControlAway, Data: []byte(partID)}
	c.sendMessage(msg)
}

func (c *LocalClient) StayArm(partID string) {
	msg := ClientMessage{Code: ClientCodePartitionArmControlStayArm, Data: []byte(partID)}
	c.sendMessage(msg)
}

func (c *LocalClient) ZeroEntryDelayArm(partID string) {
	msg := ClientMessage{Code: ClientCodePartitionArmControlZeroEntryDelay, Data: []byte(partID)}
	c.sendMessage(msg)
}

func (c *LocalClient) ArmWithCode(partID string, pin string) {
	data := append([]byte(partID), []byte(pin)...)
	msg := ClientMessage{Code: ClientCodePartitionDisarmControl, Data: data}
	c.sendMessage(msg)
}

func (c *LocalClient) Disarm(partID string, pin string) {
	data := append([]byte(partID), []byte(pin)...)
	msg := ClientMessage{Code: ClientCodePartitionDisarmControl, Data: data}
	c.sendMessage(msg)
}

func (c *LocalClient) PanicAlarm(target string) {
	data := []byte(target)
	msg := ClientMessage{Code: ClientCodePartitionDisarmControl, Data: data}
	c.sendMessage(msg)
}

func (c *LocalClient) startWriteLoop() {
	go func() {
		for {
			select {
			case msg := <-c.writeCh:
				err := writeCientMessage(msg, c.conn)
				if err != nil {
					c.errorCh <- fmt.Errorf("write error: %v", err.Error())
				}
			}
		}
	}()
}

func (c *LocalClient) startReadLoop() {
	go func() {
		for {
			msgs, err := readAvailableServerMessages(c.conn)
			for _, msg := range msgs {
				c.readCh <- msg
			}
			if err != nil {
				c.errorCh <- fmt.Errorf("read error: %v", err.Error())
			}
		}
	}()
}

func (c *LocalClient) startProcessingLoop() {
	go func() {
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
	}()
}

func (c *LocalClient) poll() {
	if c.loggedIn {
		c.sendMessage(ClientMessage{Code: ClientCodePoll})
	}
}

func (c *LocalClient) requestStateRefresh() {
	if c.loggedIn {
		c.sendMessage(ClientMessage{Code: ClientCodeStatusReport})
	}
}

func (c *LocalClient) sendMessage(msg ClientMessage) {
	c.writeCond.L.Lock()
	for len(c.msgsPendingAck) == maxPendingMessages {
		c.writeCond.Wait()
	}
	c.msgsPendingAck = append(c.msgsPendingAck, msg.Code)
	c.writeCond.L.Unlock()
	c.writeCh <- msg
}

func (c *LocalClient) processAck(msg ServerMessage) {

	codeInt, err := proto.DecodeIntCode(msg.Data)
	if err != nil {
		panic(fmt.Errorf("failed to decode code for msg %v: %v", msg, err))
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

func (c *LocalClient) processServerMessage(msg ServerMessage) error {
	switch msg.Code {

	case ServerCodeAck:
		c.processAck(msg)

	case ServerCodeLoginRes:
		c.processLoginResult(msg)

	case ServerCodeSysErr:
		c.processSystemError(msg)

	case ServerCodeKeypadLedState:
		c.processKeypadLEDState(msg)
	case ServerCodeKeypadLedFlashState:
		c.processKeypadLEDFlashState(msg)

	case ServerCodePartitionReady:
		c.processPartitionState(msg, PartitionStateReady)
	case ServerCodePartitionNotReady:
		c.processPartitionState(msg, PartitionStateNotReady)
	case ServerCodePartitionArmed:
		c.processPartitionState(msg, PartitionStateArmed)
	case ServerCodePartitionInAlarm:
		c.processPartitionState(msg, PartitionStateInAlarm)
	case ServerCodePartitionDisarmed:
		c.processPartitionState(msg, PartitionStateDisarmed)
	case ServerCodePartitionBusy:
		c.processPartitionState(msg, PartitionStateBusy)

	case ServerCodeZoneAlarm:
		c.processZoneState(msg, ZoneStateAlarm)
	case ServerCodeZoneAlarmRestore:
		c.processZoneState(msg, ZoneStateAlarmRestore)
	case ServerCodeZoneTemper:
		c.processZoneState(msg, ZoneStateTemper)
	case ServerCodeZoneTemperRestore:
		c.processZoneState(msg, ZoneStateTemperRestore)
	case ServerCodeZoneFault:
		c.processZoneState(msg, ZoneStateFault)
	case ServerCodeZoneFaultRestore:
		c.processZoneState(msg, ZoneStateFaultRestore)
	case ServerCodeZoneOpen:
		c.processZoneState(msg, ZoneStateOpen)
	case ServerCodeZoneRestore:
		c.processZoneState(msg, ZoneStateRestore)

	case ServerCodeTroubleLEDOff, ServerCodeTroubleLEDOn:
		c.processTroubleLED(msg)

	case ServerCodeExitDelayInProgress, ServerCodeEntryDelayInProgress,
		ServerCodeKeypadLockOut, ServerCodePartitionArmingFailed,
		ServerCodePGMOutputInProgress, ServerCodeChimeEnabled, ServerCodeChimeDisabled,
		ServerCodeSystemArmingInProgress, ServerCodePartialClosing,
		ServerCodeSpecialClosing, ServerCodeSpecialOpening:
		c.processPartitionEvent(LevelInfo, msg)

	case ServerCodeInvalidAccessCode:
		c.processPartitionEvent(LevelWarn, msg)

	case ServerCodeUserClosing, ServerCodeUserOpening:
		c.processUserEvent(msg)

	case ServerCodeVerboseTroubleStatus:
		return c.updateVerboseTroubleStatus(msg)

	case ServerCodePanelBatteryTrouble, ServerCodePanelACTrouble,
		ServerCodeSystemBellTrouble, ServerCodeFTCTrouble,
		ServerCodeBufferNearFull, ServerCodeGeneralSystemTamper:
		c.logEvent(LevelTrouble, msg)

	case ServerCodeDuressAlarm, ServerCodeFireAlarm, ServerCodeAuxillaryAlarm,
		ServerCodeSmokeOrAuxAlarm, ServerCodeFireTroubleAlarm, ServerCodePanicAlarm:
		c.logEvent(LevelAlarm, msg)

	default:
		c.logEvent(LevelInfo, msg)
	}

	return nil
}

func (c *LocalClient) processPartitionEvent(level EventLevel, msg ServerMessage) {
	partID := string(msg.Data)
	c.eventCh <- newServerEvent(level, msg.Code).setPartition(partID)
}

func (c *LocalClient) processUserEvent(msg ServerMessage) {
	partID := string(msg.Data[0])
	userID := string(msg.Data[1:])
	c.eventCh <- newServerEvent(LevelInfo, msg.Code).setPartition(partID).setUserID(userID)
}

func (c *LocalClient) processLoginResult(msg ServerMessage) {
	loginRes := LoginRes(msg.Data)
	if loginRes == LoginResSuccess { // login success
		c.loggedIn = true
		c.requestStateRefresh()
	} else if loginRes == LoginResFailure {
		c.errorCh <- fmt.Errorf("Login attempt rejected")
		c.loggedIn = false
	} else {
		loginMsg := ClientMessage{
			Code: ClientCodeNetworkLogin,
			Data: []byte(c.password),
		}
		c.sendMessage(loginMsg)
	}
}

func (c *LocalClient) processPartitionState(msg ServerMessage, newState PartitionState) {

	partID := string(msg.Data)

	p := c.getPartition(partID)

	if p.State != newState {
		p.State = newState
		level := LevelInfo
		if newState == PartitionStateInAlarm {
			level = LevelAlarm
		}

		c.eventCh <- newServerEvent(level, msg.Code).setPartition(partID)
	}
}

func (c *LocalClient) getPartition(partID string) *Partition {
	p, ok := c.partitions[partID]
	if !ok {
		p = newPartition(partID)
		c.partitions[partID] = p
	}
	return p
}

func (c *LocalClient) processTroubleLED(msg ServerMessage) {
	partID := string(msg.Data)
	state := msg.Code == ServerCodeTroubleLEDOn
	p := c.getPartition(partID)
	if p.TroubleStateLED != state {
		p.TroubleStateLED = state
		level := LevelInfo
		if state {
			level = LevelTrouble
		}

		c.eventCh <- newServerEvent(level, msg.Code).setPartition(partID)
	}
}

func (c *LocalClient) processKeypadLEDState(msg ServerMessage) error {
	bitset, err := decodeHexByte(msg.Data)
	if err != nil {
		return err
	}

	state := KeypadLEDState(bitset)

	p := c.getPartition("1")
	if p.KeypadLEDState != state {
		p.KeypadLEDState = state
		c.eventCh <- newServerEvent(LevelInfo, msg.Code).setPartition("1").setData("state", state)
	}
	return nil
}

func (c *LocalClient) processKeypadLEDFlashState(msg ServerMessage) error {
	bitset, err := decodeHexByte(msg.Data)
	if err != nil {
		return err
	}

	state := KeypadLEDFlashState(bitset)

	p := c.getPartition("1")
	if p.KeypadLEDFlashState != state {
		p.KeypadLEDFlashState = state
		c.eventCh <- newServerEvent(LevelInfo, msg.Code).setPartition("1").setData("state", state)
	}
	return nil
}

func (c *LocalClient) updateVerboseTroubleStatus(msg ServerMessage) error {
	bitset, err := decodeHexByte(msg.Data)
	if err != nil {
		return err
	}
	status := SystemTroubleStatus(bitset)
	if c.systemTroubleStatus != status {
		c.systemTroubleStatus = status
		level := LevelInfo
		if status != 0 {
			level = LevelTrouble
		}
		c.eventCh <- newServerEvent(level, msg.Code).setData("status", status)
	}

	return nil
}

func decodeHexByte(data []byte) (byte, error) {
	arr := make([]byte, 1)
	if _, err := hex.Decode(arr, data); err != nil {
		return 0, err
	}
	return arr[0], nil
}

func (c *LocalClient) processZoneState(msg ServerMessage, newState ZoneState) {
	var partID string
	var zoneID string
	if newState == ZoneStateFault || newState == ZoneStateFaultRestore || newState == ZoneStateOpen || newState == ZoneStateRestore {
		zoneID = string(msg.Data)
	} else {
		partID = string(msg.Data[:1])
		zoneID = string(msg.Data[1:])
	}

	z := c.getZone(zoneID)
	if z.State != newState {
		z.State = newState
		level := LevelInfo
		if newState == ZoneStateAlarm {
			level = LevelAlarm
		} else if newState == ZoneStateFault || newState == ZoneStateTemper {
			level = LevelTrouble
		}

		c.eventCh <- newServerEvent(level, msg.Code).setPartition(partID).setZoneID(zoneID)
	}
}

func (c *LocalClient) getZone(zoneID string) *Zone {
	z, ok := c.zones[zoneID]
	if !ok {
		z = newZone(zoneID)
		c.zones[zoneID] = z
	}
	return z
}

func (c *LocalClient) processSystemError(msg ServerMessage) error {
	errCode, err := proto.DecodeIntCode(msg.Data)
	errDesc := tpiErrors[errCode]

	c.eventCh <- newServerEvent(LevelError, msg.Code).setData("error", errDesc)

	if err != nil {
		return err
	}
	return nil
}

func (c *LocalClient) logEvent(lvl EventLevel, msg ServerMessage) {
	c.eventCh <- newServerEvent(lvl, msg.Code)
}
