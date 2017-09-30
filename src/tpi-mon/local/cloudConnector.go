package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
	"tpi-mon/pkg/site"
	"tpi-mon/pkg/ws"
)

type connState byte

const (
	connStateDisconnected connState = iota
	connStateConnecting
	connStateConnected
)

// cloudConnector connects the local tpi client with a remote cloud
type cloudConnector struct {
	url        string
	token      string
	connState  connState
	conn       *ws.Conn
	siteClient site.Client
	sendQueue  []interface{}

	sendCh        chan interface{}
	doneCh        chan struct{}
	loopsWg       sync.WaitGroup
	connStateLock sync.Cond
}

func startCloudConnector(url string, token string, siteClient site.Client) {

	c := &cloudConnector{
		siteClient: siteClient,
		url:        url,
		token:      token,
		sendCh:     make(chan interface{}),
	}

	c.subscribeToTpiEvents()

	go func() {
		c.connect()
	}()
}

func (c *cloudConnector) subscribeToTpiEvents() {

	eventCh := c.siteClient.SubscribeToEvents()
	stateChgCh := c.siteClient.SubscribeToStateChange()

	go func() {
		for {
			select {
			case evt := <-eventCh:
				c.enqueueMessage(evt)
			case chg := <-stateChgCh:
				c.enqueueMessage(chg)
			}
		}
	}()
}

func (c *cloudConnector) connect() {
	for n := 0.0; c.conn == nil; n++ {
		if !c.attemptConnect() {
			c.backoff(n)
		}
	}
	c.startReadLoop()
	c.startWriteLoop()
}

func (c *cloudConnector) attemptConnect() bool {
	conn, err := ws.Dial(c.url, c.token)
	if err != nil {
		log.Printf("Failed to connect to %v: %v\n", c.url, err)
		return false
	}

	c.conn = conn
	return true
}

func (c *cloudConnector) backoff(n float64) {
	// 0 -> [250ms, 500ms]
	// 1 -> [500ms, 1000ms]
	// 2 -> [1000ms, 2000ms]
	// ...
	// 8 ... n -> [64s, 128s]
	backoffFactor := math.Pow(2, math.Min(8.0, n))
	baseDelay := 250 * time.Millisecond
	backoff := time.Duration(backoffFactor+rand.Float64()*backoffFactor) * baseDelay
	log.Printf("reconnect backoff %v => %v\n", n, backoff)
	time.Sleep(backoff)
}

func (c *cloudConnector) signalConnErr(err error) {
	log.Printf("Connection error with %v: %v\n", c.url, err)
	c.connStateLock.L.Lock()
	defer c.connStateLock.L.Unlock()
	if c.connState != connStateConnecting {
		c.connState = connStateDisconnected
	}
	c.connStateLock.Signal()
}

func (c *cloudConnector) reconnectLoop() {

	for {
		c.connStateLock.L.Lock()
		for c.connState == connStateConnected {
			c.connStateLock.Wait()
		}
		c.connState = connStateConnecting
		c.connStateLock.L.Unlock()

		c.stop()
		c.connect()

		c.connStateLock.L.Lock()
		c.connState = connStateConnected
		c.connStateLock.L.Unlock()
	}

}

func (c *cloudConnector) stop() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Println("error closing connection:", err)
		}
		c.conn = nil
	}
	c.doneCh <- struct{}{}
	c.loopsWg.Wait()
}

func (c *cloudConnector) startReadLoop() {

	go func() {

		c.loopsWg.Add(1)
		defer c.loopsWg.Done()

		for {
			i, err := c.conn.Read()
			if err != nil {
				c.signalConnErr(err)
				break
			}

			switch o := i.(type) {
			case site.UserCommand:
				c.processUserCommand(o)
			case ws.ControlMessage:
				c.processControlMessage(o)
			default:
				panic(fmt.Errorf("Unexpected message: %#v", i))
			}

		}
	}()
}

func (c *cloudConnector) processUserCommand(cmd site.UserCommand) {
	if err := c.siteClient.Exec(cmd); err != nil {

		e := site.Event{
			Level:       site.LevelError,
			Code:        "UserCommandError",
			Description: err.Error(),
		}

		c.enqueueMessage(e)
	}
}

func (c *cloudConnector) processControlMessage(msg ws.ControlMessage) {
	switch msg.Code {
	case ws.CtrlGetState:
		st := c.siteClient.GetState()
		c.enqueueMessage(st)
	default:
		panic(fmt.Errorf("Unexpected controlMessage %v", msg))
	}
}

func (c *cloudConnector) enqueueMessage(msg interface{}) {
	c.sendCh <- msg
}

func (c *cloudConnector) startWriteLoop() {
	go func() {
		c.loopsWg.Add(1)
		defer c.loopsWg.Done()

		c.drainSendQueue()

		for {
			select {
			case <-c.doneCh:
				break
			case o := <-c.sendCh:
				if o != nil {
					c.sendQueue = append(c.sendQueue, o)
				}
				c.drainSendQueue()
			}
		}

	}()
}

// drainSendQueue will send all queued messages
func (c *cloudConnector) drainSendQueue() {
	ok := true
	n := 0
	for ok && n < len(c.sendQueue) {
		o := c.sendQueue[n]
		ok = c.sendMessage(o)
		if ok {
			n++
		}
	}

	if n > 0 {
		c.sendQueue = c.sendQueue[n:]
	}
}

func (c *cloudConnector) sendMessage(msg interface{}) bool {
	if c.conn == nil {
		return false
	}

	err := c.conn.Write(msg)
	if err != nil {
		c.signalConnErr(err)
	}
	return err != nil
}
