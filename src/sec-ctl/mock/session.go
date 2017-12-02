package main

import (
	"net"
	"sec-ctl/pkg/tpi"
)

type clientSession struct {
	controller       *controller
	conn             net.Conn
	readCh           chan tpi.ClientMessage
	writeCh          chan tpi.ServerMessage
	loggedIn         bool
	numLoginRequests int
}

// Handles incoming requests.
func handleClientSession(ctrl *controller, conn net.Conn) {

	session := &clientSession{
		controller: ctrl,
		conn:       conn,
		readCh:     make(chan tpi.ClientMessage),
		writeCh:    make(chan tpi.ServerMessage),
		loggedIn:   false,
	}

	ctrl.sessionStarted(session)

	session.startReadLoop()
	session.startWriteLoop()
	session.startProcessingLoop()

	// request client login
	session.writeCh <- tpi.ServerMessage{Code: tpi.ServerCodeLoginRes, Data: []byte(tpi.LoginResLoginRequest)}
}

func (s *clientSession) startReadLoop() {

	go func() {
		defer s.conn.Close()

		for {
			msgs, err := tpi.ReadAvailableClientMessages(s.conn)
			if err != nil {
				logger.Println("read error:", err)
				break
			}

			for _, m := range msgs {
				logger.Println("read:", m)
				s.readCh <- m
			}
		}

		//notify state of session end
		s.controller.sessionEnded(s)
		logger.Println("client session ended")

	}()

}

func (s *clientSession) startWriteLoop() {
	go func() {
		for {
			select {
			case m := <-s.writeCh:
				logger.Println("write:", m)
				err := m.Write(s.conn)
				if err != nil {
					logger.Println("write error:", err)
				}
			}
		}
	}()
}

func (s *clientSession) startProcessingLoop() {
	go func() {
		for {
			select {
			case msg := <-s.readCh:
				s.processClientMessage(msg)
			}
		}
	}()
}

func (s *clientSession) processClientMessage(msg tpi.ClientMessage) error {
	var replies []tpi.ServerMessage
	var err error

	ctrl := s.controller

	if msg.Code != tpi.ClientCodeNetworkLogin && !s.loggedIn {
		replies = []tpi.ServerMessage{tpi.ServerMessage{Code: tpi.ServerCodeLoginRes, Data: []byte(tpi.LoginResLoginRequest)}}
	} else {
		switch msg.Code {
		case tpi.ClientCodePoll: // noop, will just ack
		case tpi.ClientCodeNetworkLogin:
			var loggedIn bool
			loggedIn, replies, err = ctrl.processLoginRequest(msg)
			if err == nil {
				s.loggedIn = loggedIn
			}

		case tpi.ClientCodeStatusReport:
			replies, err = ctrl.processStatusReport(msg)
		case tpi.ClientCodePartitionArmControlAway:
			replies, err = ctrl.processArmControlAway(msg)
		case tpi.ClientCodePartitionArmControlStayArm:
			replies, err = ctrl.processArmControlStay(msg)
		case tpi.ClientCodePartitionArmControlWithCode:
			replies, err = ctrl.processArmControlWithCode(msg)
		case tpi.ClientCodePartitionArmControlZeroEntryDelay:
			replies, err = ctrl.processArmControlZeroEntryDelay(msg)
		case tpi.ClientCodePartitionDisarmControl:
			replies, err = ctrl.processDisarm(msg)
		default:
			logger.Println("WARN: Unhandled client message:", msg)
		}
	}

	if err != nil {
		return err
	}

	s.throttleLoggedOutClient(replies)

	s.reply(msg, replies...)
	return nil
}

func (s *clientSession) throttleLoggedOutClient(msgs []tpi.ServerMessage) {
	for _, m := range msgs {
		if m.Code == tpi.ServerCodeLoginRes {
			s.numLoginRequests++
			break
		}
	}
	if s.numLoginRequests > 4 {
		logger.Println("too many login attempts, closing session")
		s.conn.Close()
	}

}

func (s *clientSession) reply(msg tpi.ClientMessage, replies ...tpi.ServerMessage) error {
	for _, reply := range replies {
		s.writeCh <- reply
	}
	s.writeCh <- tpi.ServerMessage{
		Code: tpi.ServerCodeAck,
		Data: tpi.EncodeIntCode(int(msg.Code)),
	}

	return nil
}
