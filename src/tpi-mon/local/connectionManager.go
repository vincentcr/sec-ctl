package main

import (
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

type connState byte

const (
	connStateDisconnected connState = iota
	connStateConnecting
	connStateConnected
)

type attemptConnect func() (interface{}, error)

type connectionManager struct {
	attemptConnect attemptConnect
	connState      connState
	connStateLock  *sync.Cond
	name           string
	conn           interface{}
}

func newConnectionManager(name string, attemptConnect attemptConnect) *connectionManager {
	return &connectionManager{
		name:           name,
		attemptConnect: attemptConnect,
		connState:      connStateDisconnected,
		connStateLock:  sync.NewCond(&sync.Mutex{}),
	}
}

func (mgr *connectionManager) connect() {
	for n := 0; mgr.connState != connStateConnected; n++ {

		conn, err := mgr.attemptConnect()

		if err != nil {
			mgr.log("connection attempt failed: %v", err)
			mgr.backoff(n)
		} else {
			mgr.log("connected")
			mgr.conn = conn
			mgr.connStateLock.L.Lock()
			mgr.connState = connStateConnected
			mgr.connStateLock.L.Unlock()
			mgr.connStateLock.Broadcast()
		}
	}
}

func (mgr *connectionManager) backoff(n int) {
	// 0 -> [250ms, 500ms]
	// 1 -> [500ms, 1000ms]
	// 2 -> [1000ms, 2000ms]
	// ...
	// 8 ... n -> [64s, 128s]
	baseDelayMillis := 250.0
	backoffFactor := math.Pow(2, math.Min(8.0, float64(n)))
	backoff := time.Duration((backoffFactor+rand.Float64()*backoffFactor)*baseDelayMillis) * time.Millisecond
	mgr.log("backoff %v => %v", n, backoff)
	time.Sleep(backoff)
}

func (mgr *connectionManager) startReconnectLoop() {
	go func() {
		for {
			mgr.connStateLock.L.Lock()
			defer mgr.connStateLock.L.Unlock()

			for mgr.connState == connStateConnected {
				mgr.connStateLock.Wait()
			}
			mgr.connState = connStateConnecting

			mgr.connect()
		}
	}()
}

func (mgr *connectionManager) signalConnErrAndWaitReconnected(err error) {
	mgr.signalConnErr(err)
	mgr.waitReconnected()
}

func (mgr *connectionManager) signalConnErr(err error) {
	mgr.log("Connection failure: %v", err)
	mgr.connStateLock.L.Lock()
	defer mgr.connStateLock.L.Unlock()
	if mgr.connState != connStateConnecting {
		mgr.connState = connStateDisconnected
	}
	mgr.connStateLock.Broadcast()
}

func (mgr *connectionManager) log(format string, args ...interface{}) {
	log.Printf(mgr.name+": "+format+"\n", args...)
}

func (mgr *connectionManager) waitReconnected() {
	mgr.connStateLock.L.Lock()
	defer mgr.connStateLock.L.Unlock()
	for mgr.connState != connStateConnected {
		mgr.connStateLock.Wait()
	}
}
