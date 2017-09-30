package main

import (
	"fmt"
	"sync"
	"time"
	"tpi-mon/pkg/site"
	"tpi-mon/pkg/tpi"
)

const defaultArmDelay = time.Second * 15

type controller struct {
	state       *state
	sessions    []*clientSession
	sessionLock *sync.Mutex
}

func newController(s *state) *controller {
	return &controller{
		state:       s,
		sessions:    make([]*clientSession, 0),
		sessionLock: &sync.Mutex{},
	}
}

// sessionStarted notifies the state of a new session
func (ctrl *controller) sessionStarted(session *clientSession) {
	ctrl.sessionLock.Lock()
	defer ctrl.sessionLock.Unlock()
	ctrl.sessions = append(ctrl.sessions, session)
}

// sessionStarted notifies the state that a session has ended
func (ctrl *controller) sessionEnded(session *clientSession) {
	ctrl.sessionLock.Lock()
	defer ctrl.sessionLock.Unlock()

	for i, s := range ctrl.sessions {
		if s == session {
			ctrl.sessions = append(ctrl.sessions[:i], ctrl.sessions[i+1:]...)
			return
		}
	}
	panic(fmt.Errorf("session not found in session list: %v", session))
}

// broadcastMessagesToClients sends the specified message to all connected clients
func (ctrl *controller) broadcastMessagesToClients(msgs ...tpi.ServerMessage) {
	for _, msg := range msgs {
		for _, s := range ctrl.sessions {
			s.writeCh <- msg
		}
	}
}

// processLoginRequest verifies the password,
// and return success/failure accordingly
func (ctrl *controller) processLoginRequest(msg tpi.ClientMessage) (bool, []tpi.ServerMessage, error) {
	var res tpi.LoginRes

	password := string(msg.Data)
	success := password == ctrl.state.password
	if success {
		res = tpi.LoginResSuccess
	} else {
		res = tpi.LoginResFailure
	}

	reply := tpi.ServerMessage{Code: tpi.ServerCodeLoginRes, Data: []byte(res)}

	return success, []tpi.ServerMessage{reply}, nil
}

func (ctrl *controller) processStatusReport(msg tpi.ClientMessage) ([]tpi.ServerMessage, error) {

	s := ctrl.state

	replies := make([]tpi.ServerMessage, 0)

	for _, z := range s.Zones {
		status := zoneStateToServerCode(z.State)
		m := tpi.ServerMessage{Code: status, Data: []byte(z.ID)}
		replies = append(replies, m)
	}

	for _, p := range s.Partitions {
		status := partitionStateToServerCode(p.State)
		m := tpi.ServerMessage{Code: status, Data: []byte(p.ID)}
		replies = append(replies, m)
	}

	for _, p := range s.Partitions {
		var code tpi.ServerCode
		if p.TroubleStateLED || p.KeypadLEDState != 0 || p.KeypadLEDFlashState != 0 {
			code = tpi.ServerCodeTroubleLEDOn
		} else if p.State != site.PartitionStateBusy {
			code = tpi.ServerCodeTroubleLEDOff
		}

		if code != 0 {
			m := tpi.ServerMessage{Code: code, Data: []byte(p.ID)}
			replies = append(replies, m)
		}
	}

	if s.TroubleStatus != 0 {
		data := tpi.EncodeIntCode(int(s.TroubleStatus))
		m := tpi.ServerMessage{Code: tpi.ServerCodeVerboseTroubleStatus, Data: data}
		replies = append(replies, m)
	}

	return replies, nil
}

func zoneStateToServerCode(state site.ZoneState) tpi.ServerCode {
	switch state {
	case site.ZoneStateAlarm:
		return tpi.ServerCodeZoneAlarm
	case site.ZoneStateAlarmRestore:
		return tpi.ServerCodeZoneAlarmRestore
	case site.ZoneStateTemper:
		return tpi.ServerCodeZoneTemper
	case site.ZoneStateTemperRestore:
		return tpi.ServerCodeZoneTemperRestore
	case site.ZoneStateFault:
		return tpi.ServerCodeZoneFault
	case site.ZoneStateFaultRestore:
		return tpi.ServerCodeZoneFaultRestore
	case site.ZoneStateOpen:
		return tpi.ServerCodeZoneOpen
	case site.ZoneStateRestore:
		return tpi.ServerCodeZoneRestore
	}

	panic(fmt.Errorf("Unmapped zone state: %v", state))
}

func partitionStateToServerCode(state site.PartitionState) tpi.ServerCode {
	switch state {
	case site.PartitionStateReady:
		return tpi.ServerCodePartitionReady
	case site.PartitionStateNotReady:
		return tpi.ServerCodePartitionNotReady
	case site.PartitionStateArmed:
		return tpi.ServerCodePartitionArmed
	case site.PartitionStateInAlarm:
		return tpi.ServerCodePartitionInAlarm
	case site.PartitionStateDisarmed:
		return tpi.ServerCodePartitionDisarmed
	case site.PartitionStateBusy:
		return tpi.ServerCodePartitionBusy
	}

	panic(fmt.Errorf("Unmapped partition state: %v", state))
}

// triggerAlarm triggers an alarm of specfied type on specified partition and zone:
//  - if applicable change the state of the target zone and partition
//  - broadcast messages to connected clients
//  - update the state to record this alarm
//
// The system will stay in alarm state until a corresponding call
// to restoreAlarm is made
func (ctrl *controller) triggerAlarm(t site.AlarmType, partID string, zoneID string) error {

	s := ctrl.state

	a := site.Alarm{
		AlarmType:   t,
		PartitionID: partID,
		ZoneID:      zoneID,
		Triggered:   time.Now(),
	}

	//prevent dupes
	if a2, _ := s.findUnrestoredAlarm(a); a2.AlarmType != "" {
		return fmt.Errorf("alarm already triggered")
	}

	if a.AlarmType == site.AlarmTypePartition {
		if a.PartitionID == "" || a.ZoneID == "" {
			return fmt.Errorf("partitionID and zoneID are required for alarm of type %v", a.AlarmType)
		}
	} else {
		if a.PartitionID != "" || a.ZoneID != "" {
			return fmt.Errorf("partitionID and zoneID are not allowed for alarm of type %v", a.AlarmType)
		}
	}

	if err := s.processAlarm(a); err != nil {
		return err
	}

	msgs, err := ctrl.processAlarm(a)
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		ctrl.broadcastMessagesToClients(msg)
	}
	return nil
}

func (ctrl *controller) processAlarm(a site.Alarm) ([]tpi.ServerMessage, error) {
	//for partition alarms we need to send partition, zone alarm, and trigger the trouble led
	msgs := make([]tpi.ServerMessage, 0, 3)

	switch a.AlarmType {
	case site.AlarmTypePartition:
		msgs = append(msgs,
			tpi.ServerMessage{
				Code: tpi.ServerCodePartitionInAlarm,
				Data: []byte(a.PartitionID),
			},
			tpi.ServerMessage{
				Code: tpi.ServerCodeZoneAlarm,
				Data: []byte(a.PartitionID + a.ZoneID),
			},
			tpi.ServerMessage{
				Code: tpi.ServerCodeTroubleLEDOn,
				Data: []byte(a.PartitionID),
			},
		)
	case site.AlarmTypeDuress:
		msgs = append(msgs, tpi.ServerMessage{Code: tpi.ServerCodeDuressAlarm})
	case site.AlarmTypeFire:
		msgs = append(msgs, tpi.ServerMessage{Code: tpi.ServerCodeFireAlarm})
	case site.AlarmTypeAux:
		msgs = append(msgs, tpi.ServerMessage{Code: tpi.ServerCodeAuxillaryAlarm})
	case site.AlarmTypeSmokeOrAux:
		msgs = append(msgs, tpi.ServerMessage{Code: tpi.ServerCodeSmokeOrAuxAlarm})
	case site.AlarmTypePanic:
		msgs = append(msgs, tpi.ServerMessage{Code: tpi.ServerCodePanicAlarm})
	default:
		return nil, fmt.Errorf("Invalid alarm type %v", a.AlarmType)
	}
	return msgs, nil
}

// restoreAlarm marks the specified alarm as restored:
//  - if applicable, reset the state of the target partition and zone
//  - broadcast relevant messages to clients
//  - the target alarm is marked as restored
//
// Note that the alarm is not immediately removed from the state.
// Instead it will be removed later by a cleanup goroutine.
// This is so that the alarm lingers and can be seen by clients that
// were not connected at the precise moment it occurred.
func (ctrl *controller) restoreAlarm(t site.AlarmType, partID string) error {

	s := ctrl.state

	a, err := s.findUnrestoredAlarm(site.Alarm{AlarmType: t, PartitionID: partID})
	if err != nil {
		return err
	}

	if err = s.processAlarmRestore(a); err != nil {
		return err
	}

	msgs, err := ctrl.processAlarmRestore(a)

	for _, msg := range msgs {
		ctrl.broadcastMessagesToClients(msg)
	}
	return nil
}

func (ctrl *controller) processAlarmRestore(a site.Alarm) ([]tpi.ServerMessage, error) {
	//for partition alarms we need to send partition, zone alarm, and trigger the trouble led
	msgs := make([]tpi.ServerMessage, 0, 3)

	switch a.AlarmType {
	case site.AlarmTypePartition:

		part, _ := ctrl.state.findPartition(a.PartitionID)
		var troubleLedCode tpi.ServerCode
		if part.TroubleStateLED {
			troubleLedCode = tpi.ServerCodeTroubleLEDOn
		} else {
			troubleLedCode = tpi.ServerCodeTroubleLEDOff
		}

		msgs = append(msgs,
			tpi.ServerMessage{
				Code: tpi.ServerCodePartitionReady,
				Data: []byte(a.PartitionID),
			},
			tpi.ServerMessage{
				Code: tpi.ServerCodeZoneRestore,
				Data: []byte(a.PartitionID + a.ZoneID),
			},
			tpi.ServerMessage{
				Code: troubleLedCode,
				Data: []byte(a.PartitionID),
			},
		)
	case site.AlarmTypeDuress: //no restore msg
	case site.AlarmTypeFire:
		msgs = append(msgs, tpi.ServerMessage{Code: tpi.ServerCodeFireAlarmRestore})
	case site.AlarmTypeAux:
		msgs = append(msgs, tpi.ServerMessage{Code: tpi.ServerCodeAuxillaryAlarmRestore})
	case site.AlarmTypeSmokeOrAux:
		msgs = append(msgs, tpi.ServerMessage{Code: tpi.ServerCodeSmokeOrAuxAlarmRestore})
	case site.AlarmTypePanic:
		msgs = append(msgs, tpi.ServerMessage{Code: tpi.ServerCodePanicAlarmRestore})

	default:
		return nil, fmt.Errorf("Invalid alarm type %v", a.AlarmType)
	}
	return msgs, nil
}

// processPartitionAlarmRestore resets the alarm state from the target zone and partition,
// and return the relevant messages to send to the clients.
func (ctrl *controller) processPartitionAlarmRestore(a *site.Alarm) ([]tpi.ServerMessage, error) {
	s := ctrl.state

	part, err := s.findPartition(a.PartitionID)
	if err != nil {
		return nil, err
	}

	part.State = site.PartitionStateReady
	part.TroubleStateLED = part.KeypadLEDState != 0 && part.KeypadLEDFlashState != 0

	zone, err := s.findZone(a.ZoneID)
	if err != nil {
		return nil, err
	}
	zone.State = site.ZoneStateRestore

	var troubleLedCode tpi.ServerCode
	if part.TroubleStateLED {
		troubleLedCode = tpi.ServerCodeTroubleLEDOn
	} else {
		troubleLedCode = tpi.ServerCodeTroubleLEDOff
	}

	msgs := []tpi.ServerMessage{
		tpi.ServerMessage{
			Code: tpi.ServerCodePartitionReady,
			Data: []byte(a.PartitionID),
		},
		tpi.ServerMessage{
			Code: tpi.ServerCodeZoneRestore,
			Data: []byte(a.PartitionID + a.ZoneID),
		},
		tpi.ServerMessage{
			Code: troubleLedCode,
			Data: []byte(a.PartitionID),
		},
	}

	return msgs, nil
}

// processPartitionLessAlarmRestore returns the relevent message that is sent
// to the clients when the specifed alarm occurs
func (ctrl *controller) processPartitionLessAlarmRestore(a *site.Alarm) ([]tpi.ServerMessage, error) {
	//simple case: just map the alarm type to the server code
	var code tpi.ServerCode
	switch a.AlarmType {
	case site.AlarmTypeDuress:
		code = 0
	case site.AlarmTypeFire:
		code = tpi.ServerCodeFireAlarmRestore
	case site.AlarmTypeAux:
		code = tpi.ServerCodeAuxillaryAlarmRestore
	case site.AlarmTypeSmokeOrAux:
		code = tpi.ServerCodeSmokeOrAuxAlarmRestore
	case site.AlarmTypePanic:
		code = tpi.ServerCodePanicAlarmRestore
	default:
		return nil, fmt.Errorf("Invalid alarm type %v", a.AlarmType)
	}

	if code == 0 {
		return []tpi.ServerMessage{}, nil
	}

	return []tpi.ServerMessage{tpi.ServerMessage{Code: code}}, nil
}

func (ctrl *controller) processArmControlAway(msg tpi.ClientMessage) ([]tpi.ServerMessage, error) {
	partID := string(msg.Data)
	return ctrl.requestArm(tpi.ArmModeAway, partID, "", defaultArmDelay)
}

func (ctrl *controller) processArmControlStay(msg tpi.ClientMessage) ([]tpi.ServerMessage, error) {
	partID := string(msg.Data)
	return ctrl.requestArm(tpi.ArmModeStay, partID, "", defaultArmDelay)
}

func (ctrl *controller) processArmControlWithCode(msg tpi.ClientMessage) ([]tpi.ServerMessage, error) {
	partID := string(msg.Data[0])
	pin := string(msg.Data[1:])

	userID, ok := ctrl.state.Users[pin]
	if !ok {
		ctrl.broadcastMessagesToClients(tpi.ServerMessage{Code: tpi.ServerCodeInvalidAccessCode})
		return []tpi.ServerMessage{}, nil
	}

	return ctrl.requestArm(tpi.ArmModeAway, partID, userID, defaultArmDelay)
}

func (ctrl *controller) processArmControlZeroEntryDelay(msg tpi.ClientMessage) ([]tpi.ServerMessage, error) {
	partID := string(msg.Data)

	return ctrl.requestArm(tpi.ArmModeAway, partID, "", 0)
}

func (ctrl *controller) requestArm(mode tpi.ArmMode, partID string, userID string, delay time.Duration) ([]tpi.ServerMessage, error) {
	s := ctrl.state

	part, err := s.findPartition(partID)
	if err != nil {
		return nil, err
	}

	if part.State != site.PartitionStateReady {
		errCode := "24"
		reply := tpi.ServerMessage{Code: tpi.ServerCodeSysErr, Data: []byte(errCode)}
		return []tpi.ServerMessage{reply}, nil
	}

	ctrl.arm(mode, part, userID, delay)

	return []tpi.ServerMessage{}, nil
}

func (ctrl *controller) arm(mode tpi.ArmMode, part site.Partition, userID string, delay time.Duration) {

	ctrl.beginArm(part, userID, delay)

	go func() {
		select {
		case <-time.After(delay):
			ctrl.completeArm(part, userID, mode)
		}
	}()
}

func (ctrl *controller) beginArm(part site.Partition, userID string, delay time.Duration) {
	msgs := make([]tpi.ServerMessage, 0)
	if delay > 0 {
		msgs = append(msgs, tpi.ServerMessage{Code: tpi.ServerCodeExitDelayInProgress})

		if userID != "" {
			msgs = append(msgs, tpi.ServerMessage{Code: tpi.ServerCodeSystemArmingInProgress})
		}
	}

	ctrl.broadcastMessagesToClients(msgs...)
}

func (ctrl *controller) completeArm(part site.Partition, userID string, mode tpi.ArmMode) {
	msgs := make([]tpi.ServerMessage, 0)

	if userID != "" {
		data := []byte(part.ID + userID)
		msgs = append(msgs, tpi.ServerMessage{Code: tpi.ServerCodeUserClosing, Data: data})
	}

	data := append(tpi.EncodeIntCode(int(mode)), []byte(part.ID)...)
	msgs = append(msgs, tpi.ServerMessage{Code: tpi.ServerCodePartitionArmed, Data: data})

	ctrl.state.armPartition(part)

	ctrl.broadcastMessagesToClients(msgs...)
}

func (ctrl *controller) processDisarm(msg tpi.ClientMessage) ([]tpi.ServerMessage, error) {

	logger.Println("processing disarm", msg)
	s := ctrl.state

	partID := string(msg.Data[0])
	part, err := s.findPartition(partID)
	if err != nil {
		return nil, err
	}

	if part.State != site.PartitionStateArmed {
		errCode := "23"
		reply := tpi.ServerMessage{Code: tpi.ServerCodeSysErr, Data: []byte(errCode)}
		return []tpi.ServerMessage{reply}, nil
	}

	pin := string(msg.Data[1:])
	userID, ok := s.Users[pin]
	if !ok {
		ctrl.broadcastMessagesToClients(tpi.ServerMessage{Code: tpi.ServerCodeInvalidAccessCode})
		return []tpi.ServerMessage{}, nil
	}

	if err = s.disarmPartition(part); err != nil {
		return nil, err
	}

	ctrl.broadcastMessagesToClients(
		tpi.ServerMessage{Code: tpi.ServerCodeUserOpening, Data: []byte(part.ID + userID)},
		tpi.ServerMessage{Code: tpi.ServerCodePartitionDisarmed, Data: []byte(part.ID)},
		tpi.ServerMessage{Code: tpi.ServerCodePartitionReady, Data: []byte(part.ID)},
	)

	return []tpi.ServerMessage{}, nil
}
