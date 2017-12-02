package main

import (
	"encoding/hex"
	"time"
	"sec-ctl/pkg/sites"
	"sec-ctl/pkg/tpi"
)

const keepAliveDelay = 30 * time.Second
const stateRefreshDelay = 300 * time.Second
const maxPendingMessages = 4

type localSite struct {
	loggedIn bool
	password string
	id       string

	conn           *localSiteConnector
	partitions     map[string]*sites.Partition
	zones          map[string]*sites.Zone
	eventChs       []chan sites.Event
	stateChangeChs []chan sites.StateChange

	systemTroubleStatus sites.SystemTroubleStatus
}

// NewLocalClient creates a new local client, from the supplied local server info
func newLocalSite(hostname string, port uint16, password string, id string) sites.Site {

	c := &localSite{
		id:             id,
		password:       password,
		partitions:     map[string]*sites.Partition{},
		zones:          map[string]*sites.Zone{},
		eventChs:       make([]chan sites.Event, 0),
		stateChangeChs: make([]chan sites.StateChange, 0),
	}

	c.conn = newLocalSiteConnector(hostname, port, c.processMessage)
	c.startTimersLoop()

	return c
}

func (c *localSite) SubscribeToEvents() chan sites.Event {
	ch := make(chan sites.Event)
	c.eventChs = append(c.eventChs, ch)
	return ch
}

func (c *localSite) SubscribeToStateChange() chan sites.StateChange {
	ch := make(chan sites.StateChange)
	c.stateChangeChs = append(c.stateChangeChs, ch)
	return ch
}

func (c *localSite) GetID() string {
	return c.id
}

func (c *localSite) GetState() sites.SystemState {
	return sites.SystemState{
		ID:            c.id,
		Partitions:    c.getPartitions(),
		Zones:         c.getZones(),
		TroubleStatus: c.systemTroubleStatus,
	}
}

func (c *localSite) getPartitions() []sites.Partition {
	parts := make([]sites.Partition, 0, len(c.partitions))
	for _, p := range c.partitions {
		parts = append(parts, *p)
	}
	return parts
}

func (c *localSite) getZones() []sites.Zone {
	zones := make([]sites.Zone, 0, len(c.zones))
	for _, z := range c.zones {
		zones = append(zones, *z)
	}
	return zones
}

func (c *localSite) Exec(cmd sites.UserCommand) error {

	var msg tpi.ClientMessage
	err := cmd.Validate()
	if err != nil {
		return err
	}

	switch cmd.Code {
	case sites.CmdArmAway:
		msg = tpi.ClientMessage{Code: tpi.ClientCodePartitionArmControlAway, Data: []byte(cmd.PartitionID)}
	case sites.CmdArmStay:
		msg = tpi.ClientMessage{Code: tpi.ClientCodePartitionArmControlStayArm, Data: []byte(cmd.PartitionID)}
	case sites.CmdArmWithZeroEntryDelay:
		msg = tpi.ClientMessage{Code: tpi.ClientCodePartitionArmControlZeroEntryDelay, Data: []byte(cmd.PartitionID)}
	case sites.CmdArmWithPIN:
		data := append([]byte(cmd.PartitionID), []byte(cmd.PIN)...)
		msg = tpi.ClientMessage{Code: tpi.ClientCodePartitionArmControlWithCode, Data: data}
	case sites.CmdDisarm:
		data := append([]byte(cmd.PartitionID), []byte(cmd.PIN)...)
		msg = tpi.ClientMessage{Code: tpi.ClientCodePartitionDisarmControl, Data: data}
	case sites.CmdPanic:
		msg = tpi.ClientMessage{Code: tpi.ClientCodeTriggerPanicAlarm, Data: []byte(cmd.PanicTarget)}
	default:
		logger.Panicf("Unhandled user command %#v", cmd)
	}

	c.enqueueMessage(msg)
	return nil
}

func (c *localSite) enqueueMessage(msg tpi.ClientMessage) {
	c.conn.enqueueMessage(msg)
}

func (c *localSite) startTimersLoop() {
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

func (c *localSite) poll() {
	if c.loggedIn {
		c.enqueueMessage(tpi.ClientMessage{Code: tpi.ClientCodePoll})
	}
}

func (c *localSite) requestStateRefresh() {
	if c.loggedIn {
		c.enqueueMessage(tpi.ClientMessage{Code: tpi.ClientCodeStatusReport})
	}
}

func (c *localSite) processMessage(i interface{}) error {

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
		c.processPartitionState(msg, sites.PartitionStateReady)
	case tpi.ServerCodePartitionNotReady:
		c.processPartitionState(msg, sites.PartitionStateNotReady)
	case tpi.ServerCodePartitionArmed:
		c.processPartitionState(msg, sites.PartitionStateArmed)
	case tpi.ServerCodePartitionInAlarm:
		c.processPartitionState(msg, sites.PartitionStateInAlarm)
	case tpi.ServerCodePartitionDisarmed:
		c.processPartitionState(msg, sites.PartitionStateDisarmed)
	case tpi.ServerCodePartitionBusy:
		c.processPartitionState(msg, sites.PartitionStateBusy)

	case tpi.ServerCodeZoneAlarm:
		c.processZoneState(msg, sites.ZoneStateAlarm)
	case tpi.ServerCodeZoneAlarmRestore:
		c.processZoneState(msg, sites.ZoneStateAlarmRestore)
	case tpi.ServerCodeZoneTemper:
		c.processZoneState(msg, sites.ZoneStateTemper)
	case tpi.ServerCodeZoneTemperRestore:
		c.processZoneState(msg, sites.ZoneStateTemperRestore)
	case tpi.ServerCodeZoneFault:
		c.processZoneState(msg, sites.ZoneStateFault)
	case tpi.ServerCodeZoneFaultRestore:
		c.processZoneState(msg, sites.ZoneStateFaultRestore)
	case tpi.ServerCodeZoneOpen:
		c.processZoneState(msg, sites.ZoneStateOpen)
	case tpi.ServerCodeZoneRestore:
		c.processZoneState(msg, sites.ZoneStateRestore)

	case tpi.ServerCodeTroubleLEDOff, tpi.ServerCodeTroubleLEDOn:
		c.processTroubleLED(msg)

	case tpi.ServerCodeExitDelayInProgress, tpi.ServerCodeEntryDelayInProgress,
		tpi.ServerCodeKeypadLockOut, tpi.ServerCodePartitionArmingFailed,
		tpi.ServerCodePGMOutputInProgress, tpi.ServerCodeChimeEnabled, tpi.ServerCodeChimeDisabled,
		tpi.ServerCodeSystemArmingInProgress, tpi.ServerCodePartialClosing,
		tpi.ServerCodeSpecialClosing, tpi.ServerCodeSpecialOpening:
		c.processPartitionEvent(sites.LevelInfo, msg)

	case tpi.ServerCodeInvalidAccessCode:
		c.processPartitionEvent(sites.LevelWarn, msg)

	case tpi.ServerCodeUserClosing, tpi.ServerCodeUserOpening:
		c.processUserEvent(msg)

	case tpi.ServerCodeVerboseTroubleStatus:
		return c.updateVerboseTroubleStatus(msg)

	case tpi.ServerCodePanelBatteryTrouble, tpi.ServerCodePanelACTrouble,
		tpi.ServerCodeSystemBellTrouble, tpi.ServerCodeFTCTrouble,
		tpi.ServerCodeBufferNearFull, tpi.ServerCodeGeneralSystemTamper:
		c.publishEvent(newServerEvent(sites.LevelTrouble, msg.Code))

	case tpi.ServerCodeDuressAlarm, tpi.ServerCodeFireAlarm, tpi.ServerCodeAuxillaryAlarm,
		tpi.ServerCodeSmokeOrAuxAlarm, tpi.ServerCodeFireTroubleAlarm, tpi.ServerCodePanicAlarm:
		c.publishEvent(newServerEvent(sites.LevelAlarm, msg.Code))

	default:
		c.publishEvent(newServerEvent(sites.LevelInfo, msg.Code))
	}

	return nil
}

func (c *localSite) publishEvent(e *sites.Event) {
	go func() { // async so that blocked consumers do not block caller
		for _, ch := range c.eventChs {
			ch <- *e
		}
	}()
}

func (c *localSite) publishStateChange(chgType sites.StateChangeType, data interface{}) {
	chg := sites.StateChange{Type: chgType, Data: data}
	go func() { // async so that blocked consumers do not block caller
		for _, ch := range c.stateChangeChs {
			ch <- chg
		}
	}()
}

func (c *localSite) processPartitionEvent(level sites.EventLevel, msg tpi.ServerMessage) {
	partID := string(msg.Data)
	c.publishEvent(newServerEvent(level, msg.Code).SetPartitionID(partID))
}

func (c *localSite) processUserEvent(msg tpi.ServerMessage) {
	partID := string(msg.Data[0])
	userID := string(msg.Data[1:])
	c.publishEvent(newServerEvent(sites.LevelInfo, msg.Code).SetPartitionID(partID).SetUserID(userID))
}

func (c *localSite) processLoginResult(msg tpi.ServerMessage) {
	loginRes := tpi.LoginRes(msg.Data)
	if loginRes == tpi.LoginResSuccess { // login success
		c.loggedIn = true
		c.requestStateRefresh()
	} else if loginRes == tpi.LoginResFailure {
		logger.Panicf("Login attempt failed: password rejected!")
		c.loggedIn = false
	} else {
		loginMsg := tpi.ClientMessage{
			Code: tpi.ClientCodeNetworkLogin,
			Data: []byte(c.password),
		}
		c.enqueueMessage(loginMsg)
	}
}

func (c *localSite) processPartitionState(msg tpi.ServerMessage, newState sites.PartitionState) {

	partID := string(msg.Data)

	p := c.getPartition(partID)

	if p.State != newState {
		p.State = newState
		level := sites.LevelInfo
		if newState == sites.PartitionStateInAlarm {
			level = sites.LevelAlarm
		}

		c.publishStateChange(sites.StateChangePartition, p)
		c.publishEvent(newServerEvent(level, msg.Code).SetPartitionID(partID))
	}
}

func (c *localSite) getPartition(partID string) *sites.Partition {
	p, ok := c.partitions[partID]
	if !ok {
		p = sites.NewPartition(partID)
		c.partitions[partID] = p
	}
	return p
}

func (c *localSite) processTroubleLED(msg tpi.ServerMessage) {
	partID := string(msg.Data)
	state := msg.Code == tpi.ServerCodeTroubleLEDOn
	p := c.getPartition(partID)
	if p.TroubleStateLED != state {
		p.TroubleStateLED = state
		level := sites.LevelInfo
		if state {
			level = sites.LevelTrouble
		}
		c.publishStateChange(sites.StateChangePartition, p)

		c.publishEvent(newServerEvent(level, msg.Code).SetPartitionID(partID))
	}
}

func (c *localSite) processKeypadLEDState(msg tpi.ServerMessage) error {
	bitset, err := decodeHexByte(msg.Data)
	if err != nil {
		return err
	}

	state := sites.KeypadLEDState(bitset)

	p := c.getPartition("1")
	if p.KeypadLEDState != state {
		p.KeypadLEDState = state
		c.publishStateChange(sites.StateChangePartition, p)
		c.publishEvent(newServerEvent(sites.LevelInfo, msg.Code).SetPartitionID("1").SetData("state", state))
	}
	return nil
}

func (c *localSite) processKeypadLEDFlashState(msg tpi.ServerMessage) error {
	bitset, err := decodeHexByte(msg.Data)
	if err != nil {
		return err
	}

	state := sites.KeypadLEDFlashState(bitset)

	p := c.getPartition("1")
	if p.KeypadLEDFlashState != state {
		p.KeypadLEDFlashState = state
		c.publishStateChange(sites.StateChangePartition, p)
		c.publishEvent(newServerEvent(sites.LevelInfo, msg.Code).SetPartitionID("1").SetData("state", state))
	}
	return nil
}

func (c *localSite) updateVerboseTroubleStatus(msg tpi.ServerMessage) error {
	bitset, err := decodeHexByte(msg.Data)
	if err != nil {
		return err
	}
	status := sites.SystemTroubleStatus(bitset)
	if c.systemTroubleStatus != status {
		c.systemTroubleStatus = status
		level := sites.LevelInfo
		if status != 0 {
			level = sites.LevelTrouble
		}
		c.publishStateChange(sites.StateChangeSystemTroubleStatus, status)

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

func (c *localSite) processZoneState(msg tpi.ServerMessage, newState sites.ZoneState) {
	var partID string
	var zoneID string
	if newState == sites.ZoneStateFault || newState == sites.ZoneStateFaultRestore || newState == sites.ZoneStateOpen || newState == sites.ZoneStateRestore {
		zoneID = string(msg.Data)
	} else {
		partID = string(msg.Data[:1])
		zoneID = string(msg.Data[1:])
	}

	z := c.getZone(zoneID)
	if z.State != newState {
		z.State = newState
		level := sites.LevelInfo
		if newState == sites.ZoneStateAlarm {
			level = sites.LevelAlarm
		} else if newState == sites.ZoneStateFault || newState == sites.ZoneStateTemper {
			level = sites.LevelTrouble
		}

		c.publishStateChange(sites.StateChangeZone, z)
		c.publishEvent(newServerEvent(level, msg.Code).SetPartitionID(partID).SetZoneID(zoneID))
	}
}

func (c *localSite) getZone(zoneID string) *sites.Zone {
	z, ok := c.zones[zoneID]
	if !ok {
		z = sites.NewZone(zoneID)
		c.zones[zoneID] = z
	}
	return z
}

func (c *localSite) processSystemError(msg tpi.ServerMessage) error {
	errCode, err := tpi.DecodeIntCode(msg.Data)
	errDesc := tpi.GetErrorCodeDescription(errCode)

	c.publishEvent(newServerEvent(sites.LevelError, msg.Code).SetData("error", errDesc))

	if err != nil {
		return err
	}
	return nil
}

func newServerEvent(level sites.EventLevel, code tpi.ServerCode) *sites.Event {
	e := sites.NewEvent(level, code.Name())
	desc := tpi.GetServerCodeDescription(code)
	e.SetDescription(desc)
	return e
}
