package main

import (
	"encoding/hex"
	"fmt"
	"time"
	"tpi-mon/pkg/site"
	"tpi-mon/pkg/tpi"
)

const keepAliveDelay = 30 * time.Second
const stateRefreshDelay = 300 * time.Second
const maxPendingMessages = 4

type localClient struct {
	loggedIn bool
	password string
	id       string

	conn           *localClientConnector
	partitions     map[string]*site.Partition
	zones          map[string]*site.Zone
	eventChs       []chan site.Event
	stateChangeChs []chan site.StateChange

	errorCh             chan error
	systemTroubleStatus site.SystemTroubleStatus
}

// NewLocalClient creates a new local client, from the supplied local server info
func newLocalClient(hostname string, port uint16, password string) site.Client {

	c := &localClient{
		id:             "1",
		password:       password,
		partitions:     map[string]*site.Partition{},
		zones:          map[string]*site.Zone{},
		eventChs:       make([]chan site.Event, 0),
		stateChangeChs: make([]chan site.StateChange, 0),
	}

	c.conn = newLocalConnection(hostname, port, c.processMessage)
	c.startTimersLoop()

	return c
}

func (c *localClient) SubscribeToEvents() chan site.Event {
	ch := make(chan site.Event)
	c.eventChs = append(c.eventChs, ch)
	return ch
}

func (c *localClient) SubscribeToStateChange() chan site.StateChange {
	ch := make(chan site.StateChange)
	c.stateChangeChs = append(c.stateChangeChs, ch)
	return ch
}

func (c *localClient) GetState() site.SystemState {
	return site.SystemState{
		ID:            c.id,
		Partitions:    c.getPartitions(),
		Zones:         c.getZones(),
		TroubleStatus: c.systemTroubleStatus,
	}
}

func (c *localClient) getPartitions() []site.Partition {
	parts := make([]site.Partition, 0, len(c.partitions))
	for _, p := range c.partitions {
		parts = append(parts, *p)
	}
	return parts
}

func (c *localClient) getZones() []site.Zone {
	zones := make([]site.Zone, 0, len(c.zones))
	for _, z := range c.zones {
		zones = append(zones, *z)
	}
	return zones
}

func (c *localClient) Exec(cmd site.UserCommand) error {

	var msg tpi.ClientMessage
	err := cmd.Validate()
	if err != nil {
		return err
	}

	switch cmd.Code {
	case site.CmdArmAway:
		msg = tpi.ClientMessage{Code: tpi.ClientCodePartitionArmControlAway, Data: []byte(cmd.PartitionID)}
	case site.CmdArmStay:
		msg = tpi.ClientMessage{Code: tpi.ClientCodePartitionArmControlStayArm, Data: []byte(cmd.PartitionID)}
	case site.CmdArmWithZeroEntryDelay:
		msg = tpi.ClientMessage{Code: tpi.ClientCodePartitionArmControlZeroEntryDelay, Data: []byte(cmd.PartitionID)}
	case site.CmdArmWithPIN:
		data := append([]byte(cmd.PartitionID), []byte(cmd.PIN)...)
		msg = tpi.ClientMessage{Code: tpi.ClientCodePartitionArmControlWithCode, Data: data}
	case site.CmdDisarm:
		data := append([]byte(cmd.PartitionID), []byte(cmd.PIN)...)
		msg = tpi.ClientMessage{Code: tpi.ClientCodePartitionDisarmControl, Data: data}
	case site.CmdPanic:
		msg = tpi.ClientMessage{Code: tpi.ClientCodeTriggerPanicAlarm, Data: []byte(cmd.PanicTarget)}
	default:
		panic(fmt.Errorf("Unhandled user command %#v", cmd))
	}

	c.enqueueMessage(msg)
	return nil
}

func (c *localClient) enqueueMessage(msg tpi.ClientMessage) {
	c.conn.enqueueMessage(msg)
}

func (c *localClient) startTimersLoop() {
	go func() {
		tickKeepAlive := time.Tick(keepAliveDelay)
		tickStateRefreshDelay := time.Tick(stateRefreshDelay)

		for {
			select {
			case <-tickKeepAlive:
				c.poll()
			case <-tickStateRefreshDelay:
				c.requestStateRefresh()
			}
		}
	}()
}

func (c *localClient) poll() {
	if c.loggedIn {
		c.enqueueMessage(tpi.ClientMessage{Code: tpi.ClientCodePoll})
	}
}

func (c *localClient) requestStateRefresh() {
	if c.loggedIn {
		c.enqueueMessage(tpi.ClientMessage{Code: tpi.ClientCodeStatusReport})
	}
}

func (c *localClient) processMessage(i interface{}) error {

	msg := i.(tpi.ServerMessage)

	switch msg.Code {

	case tpi.ServerCodeLoginRes:
		c.processLoginResult(msg)

	case tpi.ServerCodeAck: // ignore

	case tpi.ServerCodeSysErr:
		c.processSystemError(msg)

	case tpi.ServerCodeKeypadLedState:
		c.processKeypadLEDState(msg)
	case tpi.ServerCodeKeypadLedFlashState:
		c.processKeypadLEDFlashState(msg)

	case tpi.ServerCodePartitionReady:
		c.processPartitionState(msg, site.PartitionStateReady)
	case tpi.ServerCodePartitionNotReady:
		c.processPartitionState(msg, site.PartitionStateNotReady)
	case tpi.ServerCodePartitionArmed:
		c.processPartitionState(msg, site.PartitionStateArmed)
	case tpi.ServerCodePartitionInAlarm:
		c.processPartitionState(msg, site.PartitionStateInAlarm)
	case tpi.ServerCodePartitionDisarmed:
		c.processPartitionState(msg, site.PartitionStateDisarmed)
	case tpi.ServerCodePartitionBusy:
		c.processPartitionState(msg, site.PartitionStateBusy)

	case tpi.ServerCodeZoneAlarm:
		c.processZoneState(msg, site.ZoneStateAlarm)
	case tpi.ServerCodeZoneAlarmRestore:
		c.processZoneState(msg, site.ZoneStateAlarmRestore)
	case tpi.ServerCodeZoneTemper:
		c.processZoneState(msg, site.ZoneStateTemper)
	case tpi.ServerCodeZoneTemperRestore:
		c.processZoneState(msg, site.ZoneStateTemperRestore)
	case tpi.ServerCodeZoneFault:
		c.processZoneState(msg, site.ZoneStateFault)
	case tpi.ServerCodeZoneFaultRestore:
		c.processZoneState(msg, site.ZoneStateFaultRestore)
	case tpi.ServerCodeZoneOpen:
		c.processZoneState(msg, site.ZoneStateOpen)
	case tpi.ServerCodeZoneRestore:
		c.processZoneState(msg, site.ZoneStateRestore)

	case tpi.ServerCodeTroubleLEDOff, tpi.ServerCodeTroubleLEDOn:
		c.processTroubleLED(msg)

	case tpi.ServerCodeExitDelayInProgress, tpi.ServerCodeEntryDelayInProgress,
		tpi.ServerCodeKeypadLockOut, tpi.ServerCodePartitionArmingFailed,
		tpi.ServerCodePGMOutputInProgress, tpi.ServerCodeChimeEnabled, tpi.ServerCodeChimeDisabled,
		tpi.ServerCodeSystemArmingInProgress, tpi.ServerCodePartialClosing,
		tpi.ServerCodeSpecialClosing, tpi.ServerCodeSpecialOpening:
		c.processPartitionEvent(site.LevelInfo, msg)

	case tpi.ServerCodeInvalidAccessCode:
		c.processPartitionEvent(site.LevelWarn, msg)

	case tpi.ServerCodeUserClosing, tpi.ServerCodeUserOpening:
		c.processUserEvent(msg)

	case tpi.ServerCodeVerboseTroubleStatus:
		return c.updateVerboseTroubleStatus(msg)

	case tpi.ServerCodePanelBatteryTrouble, tpi.ServerCodePanelACTrouble,
		tpi.ServerCodeSystemBellTrouble, tpi.ServerCodeFTCTrouble,
		tpi.ServerCodeBufferNearFull, tpi.ServerCodeGeneralSystemTamper:
		c.publishEvent(newServerEvent(site.LevelTrouble, msg.Code))

	case tpi.ServerCodeDuressAlarm, tpi.ServerCodeFireAlarm, tpi.ServerCodeAuxillaryAlarm,
		tpi.ServerCodeSmokeOrAuxAlarm, tpi.ServerCodeFireTroubleAlarm, tpi.ServerCodePanicAlarm:
		c.publishEvent(newServerEvent(site.LevelAlarm, msg.Code))

	default:
		c.publishEvent(newServerEvent(site.LevelInfo, msg.Code))
	}

	return nil
}

func (c *localClient) publishEvent(e *site.Event) {
	go func() { // async so that blocked consumers do not block caller
		for _, ch := range c.eventChs {
			ch <- *e
		}
	}()
}

func (c *localClient) publishStateChange(chgType site.StateChangeType, data interface{}) {
	chg := site.StateChange{Type: chgType, Data: data}
	go func() { // async so that blocked consumers do not block caller
		for _, ch := range c.stateChangeChs {
			ch <- chg
		}
	}()
}

func (c *localClient) processPartitionEvent(level site.EventLevel, msg tpi.ServerMessage) {
	partID := string(msg.Data)
	c.publishEvent(newServerEvent(level, msg.Code).SetPartitionID(partID))
}

func (c *localClient) processUserEvent(msg tpi.ServerMessage) {
	partID := string(msg.Data[0])
	userID := string(msg.Data[1:])
	c.publishEvent(newServerEvent(site.LevelInfo, msg.Code).SetPartitionID(partID).SetUserID(userID))
}

func (c *localClient) processLoginResult(msg tpi.ServerMessage) {
	loginRes := tpi.LoginRes(msg.Data)
	if loginRes == tpi.LoginResSuccess { // login success
		c.loggedIn = true
		c.requestStateRefresh()
	} else if loginRes == tpi.LoginResFailure {
		c.errorCh <- fmt.Errorf("Login attempt rejected")
		c.loggedIn = false
	} else {
		loginMsg := tpi.ClientMessage{
			Code: tpi.ClientCodeNetworkLogin,
			Data: []byte(c.password),
		}
		c.enqueueMessage(loginMsg)
	}
}

func (c *localClient) processPartitionState(msg tpi.ServerMessage, newState site.PartitionState) {

	partID := string(msg.Data)

	p := c.getPartition(partID)

	if p.State != newState {
		p.State = newState
		level := site.LevelInfo
		if newState == site.PartitionStateInAlarm {
			level = site.LevelAlarm
		}

		c.publishStateChange(site.StateChangePartition, p)
		c.publishEvent(newServerEvent(level, msg.Code).SetPartitionID(partID))
	}
}

func (c *localClient) getPartition(partID string) *site.Partition {
	p, ok := c.partitions[partID]
	if !ok {
		p = site.NewPartition(partID)
		c.partitions[partID] = p
	}
	return p
}

func (c *localClient) processTroubleLED(msg tpi.ServerMessage) {
	partID := string(msg.Data)
	state := msg.Code == tpi.ServerCodeTroubleLEDOn
	p := c.getPartition(partID)
	if p.TroubleStateLED != state {
		p.TroubleStateLED = state
		level := site.LevelInfo
		if state {
			level = site.LevelTrouble
		}
		c.publishStateChange(site.StateChangePartition, p)

		c.publishEvent(newServerEvent(level, msg.Code).SetPartitionID(partID))
	}
}

func (c *localClient) processKeypadLEDState(msg tpi.ServerMessage) error {
	bitset, err := decodeHexByte(msg.Data)
	if err != nil {
		return err
	}

	state := site.KeypadLEDState(bitset)

	p := c.getPartition("1")
	if p.KeypadLEDState != state {
		p.KeypadLEDState = state
		c.publishStateChange(site.StateChangePartition, p)
		c.publishEvent(newServerEvent(site.LevelInfo, msg.Code).SetPartitionID("1").SetData("state", state))
	}
	return nil
}

func (c *localClient) processKeypadLEDFlashState(msg tpi.ServerMessage) error {
	bitset, err := decodeHexByte(msg.Data)
	if err != nil {
		return err
	}

	state := site.KeypadLEDFlashState(bitset)

	p := c.getPartition("1")
	if p.KeypadLEDFlashState != state {
		p.KeypadLEDFlashState = state
		c.publishStateChange(site.StateChangePartition, p)
		c.publishEvent(newServerEvent(site.LevelInfo, msg.Code).SetPartitionID("1").SetData("state", state))
	}
	return nil
}

func (c *localClient) updateVerboseTroubleStatus(msg tpi.ServerMessage) error {
	bitset, err := decodeHexByte(msg.Data)
	if err != nil {
		return err
	}
	status := site.SystemTroubleStatus(bitset)
	if c.systemTroubleStatus != status {
		c.systemTroubleStatus = status
		level := site.LevelInfo
		if status != 0 {
			level = site.LevelTrouble
		}
		c.publishStateChange(site.StateChangeSystemTroubleStatus, status)

		c.publishEvent(newServerEvent(level, msg.Code).SetData("status", status))
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

func (c *localClient) processZoneState(msg tpi.ServerMessage, newState site.ZoneState) {
	var partID string
	var zoneID string
	if newState == site.ZoneStateFault || newState == site.ZoneStateFaultRestore || newState == site.ZoneStateOpen || newState == site.ZoneStateRestore {
		zoneID = string(msg.Data)
	} else {
		partID = string(msg.Data[:1])
		zoneID = string(msg.Data[1:])
	}

	z := c.getZone(zoneID)
	if z.State != newState {
		z.State = newState
		level := site.LevelInfo
		if newState == site.ZoneStateAlarm {
			level = site.LevelAlarm
		} else if newState == site.ZoneStateFault || newState == site.ZoneStateTemper {
			level = site.LevelTrouble
		}

		c.publishStateChange(site.StateChangeZone, z)
		c.publishEvent(newServerEvent(level, msg.Code).SetPartitionID(partID).SetZoneID(zoneID))
	}
}

func (c *localClient) getZone(zoneID string) *site.Zone {
	z, ok := c.zones[zoneID]
	if !ok {
		z = site.NewZone(zoneID)
		c.zones[zoneID] = z
	}
	return z
}

func (c *localClient) processSystemError(msg tpi.ServerMessage) error {
	errCode, err := tpi.DecodeIntCode(msg.Data)
	errDesc := tpi.GetErrorCodeDescription(errCode)

	c.publishEvent(newServerEvent(site.LevelError, msg.Code).SetData("error", errDesc))

	if err != nil {
		return err
	}
	return nil
}

func newServerEvent(level site.EventLevel, code tpi.ServerCode) *site.Event {
	e := site.NewEvent(level, code.Name())
	desc := tpi.GetServerCodeDescription(code)
	e.SetDescription(desc)
	return e
}
