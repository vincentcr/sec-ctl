package main

import (
	"fmt"
	"net"
	"tpi-mon/pkg/tpi"
)

type localClientConnector struct {
	sendQueue *workQueue
	recvQueue *workQueue
	connMgr   *connectionManager
}

// NewLocalClient creates a new local client, from the supplied local server info
func newLocalConnection(hostname string, port uint16, recvFunc workQueueFunc) *localClientConnector {
	c := &localClientConnector{}

	c.connMgr = newConnectionManager("local site", func() (interface{}, error) {
		servAddr := fmt.Sprintf("%s:%d", hostname, port)
		tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
		if err != nil {
			return nil, err
		}

		return net.DialTCP("tcp", nil, tcpAddr)
	})

	c.recvQueue = newWorkQueue(recvFunc)
	c.sendQueue = newWorkQueue(c.sendMessage)

	go func() {
		c.connMgr.connect()
		c.connMgr.startReconnectLoop()
		c.startReadLoop()
		c.sendQueue.start()
		c.recvQueue.start()
	}()

	return c
}

func (c *localClientConnector) startReadLoop() {
	go func() {
		for {
			conn := c.connMgr.conn.(*net.TCPConn)
			msgs, err := tpi.ReadAvailableServerMessages(conn)
			if err != nil {
				c.connMgr.signalConnErrAndWaitReconnected(err)
			} else {
				for _, msg := range msgs {
					c.recvQueue.enqueue(msg)
				}
			}
		}
	}()
}

func (c *localClientConnector) enqueueMessage(msg tpi.ClientMessage) {
	c.sendQueue.enqueue(msg)
}

func (c *localClientConnector) sendMessage(i interface{}) error {
	msg := i.(tpi.ClientMessage)
	conn := c.connMgr.conn.(*net.TCPConn)
	err := msg.Write(conn)
	if err != nil {
		c.connMgr.signalConnErrAndWaitReconnected(err)
	}
	return nil
}
