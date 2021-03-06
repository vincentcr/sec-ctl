package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"sync"
	"time"

	"sec-ctl/pkg/sites"
)

const eventExpireDelay = time.Second * 60
const eventCleanupInterval = time.Second * 15

var errNoChange = errors.New("no change")

// state represents the state of the mocked system.
// It is deserialized from json at startup,
// and serialized to json every time it changes.
type state struct {
	stateFname string
	writeLock  *sync.Mutex
	password   string

	Users map[string]string // PIN -> ID
	sites.SystemState
	// Partitions    []*sites.Partition
	// Zones         []*sites.Zone
	// TroubleStatus tpi.SystemTroubleStatus
	// Alarms        []*sites.Alarm
}

// creates a new state object from stateFilename
func newState(password string, stateFilename string) (*state, error) {

	if !path.IsAbs(stateFilename) {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		fmt.Println("stateFilename is not abs:", stateFilename, ". prepending cwd", cwd)
		stateFilename = path.Clean(path.Join(cwd, stateFilename))
	}

	state := &state{
		stateFname: stateFilename,
		password:   password,
		writeLock:  &sync.Mutex{},
	}

	if err := state.load(); err != nil {
		return nil, err
	}

	state.startAlarmCleanupTimer()

	state.dumpToLogs()

	return state, nil
}

// dumpToLogs writes the jsonified state to the logs
func (state *state) dumpToLogs() {
	json, err := state.toJSON()
	if err != nil {
		panic(err)
	}
	logger.Printf("--\nstate:\n%v\n--\n", string(json))
}

func (state *state) toJSON() ([]byte, error) {
	return json.MarshalIndent(state, "", "  ")
}

// load unmarshals the state at `stateFname` into the current state
func (state *state) load() error {
	f, err := os.Open(state.stateFname)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &state); err != nil {
		return err
	}
	return nil
}

// save dumps the jsonified state at `stateFname`
func (state *state) save() error {
	data, err := state.toJSON()
	if err != nil {
		return err
	}

	tmpName := fmt.Sprintf("tmp-tpi-mock-%v%v.json", rand.Int31(), time.Now().Unix())
	tmpPath := path.Join(os.TempDir(), tmpName)
	err = ioutil.WriteFile(tmpPath, data, os.ModePerm)
	if err != nil {
		os.Remove(tmpPath)
		return err
	}

	return os.Rename(tmpPath, state.stateFname)
}

// startAlarmCleanupTimer calls cleanupAlarms every `eventCleanupInterval``
func (state *state) startAlarmCleanupTimer() {

	run := func() {
		if err := state.updateState(state.cleanupAlarms); err != nil {
			logger.Panicln(err)
		}
	}

	go func() {
		timer := time.Tick(eventCleanupInterval)

		for {
			select {
			case <-timer:
				run()
			}
		}
	}()

	run()
}

// cleanupAlarms removes alarms from the alarm if they have been restored
// for longer than `eventExpireDelay` ago.
func (state *state) cleanupAlarms() error {

	expiredIndices := make([]int, 0, len(state.Alarms))
	for i, a := range state.Alarms {
		if !a.Restored.IsZero() && a.Restored.Add(eventExpireDelay).Before(time.Now()) {
			logger.Printf("Alarm %v has expired, removing\n", a)
			expiredIndices = append(expiredIndices, i)
		}
	}

	//go through indices in reverse so we can remove elements while looping
	for i := len(expiredIndices) - 1; i >= 0; i-- {
		idx := expiredIndices[i]
		state.Alarms = append(state.Alarms[:idx], state.Alarms[idx+1:]...)
	}

	if len(expiredIndices) == 0 {
		return errNoChange
	}

	return nil
}

// updateState updates the state of the system in a thread-safe manner:
// first, a lock is acquired, then an `updater` function is called.
// if it returns true, the state is saved
func (state *state) updateState(updater func() error) error {
	state.writeLock.Lock()
	defer state.writeLock.Unlock()

	err := updater()

	if err == errNoChange {
		return nil
	} else if err != nil {
		return err
	}

	if err := state.save(); err != nil {
		return err
	}

	state.dumpToLogs()

	return nil
}

// findPartition finds a partition by id
func (state *state) findPartition(partID string) (sites.Partition, error) {
	for _, part := range state.Partitions {
		if part.ID == partID {
			return part, nil
		}
	}
	return sites.Partition{}, fmt.Errorf("partition %v not found", partID)
}

// findZone finds a zone by id
func (state *state) findZone(zoneID string) (sites.Zone, error) {
	for _, zone := range state.Zones {
		if zone.ID == zoneID {
			return zone, nil
		}
	}
	return sites.Zone{}, fmt.Errorf("zone %v not found", zoneID)
}

// processPartitionAlarm puts the target zone and partition in alarm state
// and returns the relevant messages that must be sent to the client
func (state *state) processAlarm(a sites.Alarm) error {

	return state.updateState(func() error {

		if a.AlarmType == sites.AlarmTypePartition {
			part, err := state.findPartition(a.PartitionID)
			if err != nil {
				return err
			}

			part.State = sites.PartitionStateInAlarm
			part.TroubleStateLED = true

			zone, err := state.findZone(a.ZoneID)
			if err != nil {
				return err
			}

			zone.State = sites.ZoneStateAlarm
		}

		state.Alarms = append(state.Alarms, a)

		return nil
	})
}

func (state *state) processAlarmRestore(a sites.Alarm) error {

	return state.updateState(func() error {

		a.Restored = time.Now()

		if a.AlarmType == sites.AlarmTypePartition {
			part, err := state.findPartition(a.PartitionID)
			if err != nil {
				return err
			}
			zone, err := state.findZone(a.ZoneID)
			if err != nil {
				return err
			}
			part.State = sites.PartitionStateReady
			part.TroubleStateLED = part.KeypadLEDState != 0 && part.KeypadLEDFlashState != 0
			zone.State = sites.ZoneStateRestore
		}

		return nil
	})

}

// findUnrestoredAlarm finds an unrestored alarm by type and partition
func (state *state) findUnrestoredAlarm(a sites.Alarm) (sites.Alarm, error) {
	for _, a2 := range state.Alarms {
		if a.AlarmType == a2.AlarmType && a.PartitionID == a2.PartitionID && a2.Restored.IsZero() {
			return a2, nil
		}
	}
	return sites.Alarm{}, fmt.Errorf("alarm (%v,%v) not found", a.AlarmType, a.PartitionID)
}

func (state *state) armPartition(part sites.Partition) error {
	return state.updateState(func() error {
		part.State = sites.PartitionStateArmed
		return nil
	})
}

func (state *state) disarmPartition(part sites.Partition) error {
	return state.updateState(func() error {
		part.State = sites.PartitionStateReady
		return nil
	})
}
