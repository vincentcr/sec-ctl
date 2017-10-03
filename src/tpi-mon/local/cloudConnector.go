package main

import (
	"sync"
	"time"
	"tpi-mon/pkg/site"
	"tpi-mon/pkg/ws"

	"golang.org/x/time/rate"
)

const writeBurstLimit = 128
const writeRateLimit = 128 * time.Millisecond

// cloudConnector connects the local tpi client with a remote cloud
type cloudConnector struct {
	connState  connState
	connMgr    *connectionManager
	siteClient site.Client

	sendQueue     *workQueue
	recvQueue     *workQueue
	writeLimiter  *rate.Limiter
	connStateLock sync.Cond
}

func startCloudConnector(url string, token string, siteClient site.Client) {

	c := &cloudConnector{
		siteClient:   siteClient,
		writeLimiter: rate.NewLimiter(rate.Limit(1024), 256),
	}

	c.connMgr = newConnectionManager("cloud", func() (interface{}, error) {
		return ws.Dial(url, token)
	})

	c.sendQueue = newWorkQueue(c.sendMessage)
	c.recvQueue = newWorkQueue(c.recvMessage)

	c.subscribeToTpiEvents()

	go func() {
		c.connMgr.connect()
		c.connMgr.startReconnectLoop()
		c.startReadLoop()
		c.sendQueue.start()
		c.recvQueue.start()
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

func (c *cloudConnector) startReadLoop() {

	go func() {
		for {

			conn := c.connMgr.conn.(*ws.Conn)

			o, err := conn.Read()
			if err != nil {
				c.connMgr.signalConnErrAndWaitReconnected(err)
				logger.Println("cloudConnector: read loop: reconnected, resuming")
			} else {
				c.recvQueue.enqueue(o)
			}
		}
	}()
}

func (c *cloudConnector) recvMessage(i interface{}) error {
	switch o := i.(type) {
	case site.UserCommand:
		c.recvUserCommand(o)
	case ws.ControlMessage:
		c.recvControlMessage(o)
	default:
		logger.Panicf("Unexpected message: %#v", i)
	}
	return nil
}

func (c *cloudConnector) recvUserCommand(cmd site.UserCommand) {
	if err := c.siteClient.Exec(cmd); err != nil {

		e := site.Event{
			Level:       site.LevelError,
			Code:        "UserCommandError",
			Description: err.Error(),
		}

		c.enqueueMessage(e)
	}
}

func (c *cloudConnector) recvControlMessage(msg ws.ControlMessage) {
	switch msg.Code {
	case ws.CtrlGetState:
		st := c.siteClient.GetState()
		c.enqueueMessage(st)
	default:
		logger.Panicf("Unexpected controlMessage %v", msg)
	}
}

func (c *cloudConnector) enqueueMessage(msg interface{}) {
	c.sendQueue.enqueue(msg)
}

func (c *cloudConnector) sendMessage(msg interface{}) error {

	r := c.writeLimiter.Reserve()
	if !r.OK() {
		logger.Panicf("impossible! not allowed to request a burst of 1")
	}
	time.Sleep(r.Delay())

	conn := c.connMgr.conn.(*ws.Conn)

	err := conn.Write(msg)
	if err != nil { // todo: try to separate io errors from others
		c.connMgr.signalConnErrAndWaitReconnected(err)
	}

	return nil
}
