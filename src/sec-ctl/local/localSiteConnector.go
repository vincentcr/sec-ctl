package main

import (
	"fmt"
	"net"
	"sec-ctl/pkg/tpi"
)

type localSiteConnector struct {
	sendQueue *workQueue
	recvQueue *workQueue
	connMgr   *connectionManager
}

// NewLocalClient creates a new local client, from the supplied local server info
func newLocalSiteConnector(hostname string, port uint16, recvFunc workQueueFunc) *localSiteConnector {
	c := &localSiteConnector{}

	c.connMgr = newConnectionManager("local sites", func() (interface{}, error) {
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

func (c *localSiteConnector) startReadLoop() {
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

func (c *localSiteConnector) enqueueMessage(msg tpi.ClientMessage) {
	c.sendQueue.enqueue(msg)
}

func (c *localSiteConnector) sendMessage(i interface{}) error {
	msg := i.(tpi.ClientMessage)
	conn := c.connMgr.conn.(*net.TCPConn)
	err := msg.Write(conn)
	if err != nil {
		c.connMgr.signalConnErrAndWaitReconnected(err)
	}
	return nil
}
