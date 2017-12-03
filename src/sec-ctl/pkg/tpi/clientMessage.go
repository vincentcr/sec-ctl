package tpi

import (
	"fmt"
	"io"
)

// ClientMessage represents a message sent by the client
type ClientMessage struct {
	Code ClientCode
	Data []byte
}

func (m ClientMessage) String() string {
	return fmt.Sprintf("ClientMessage{code: %v(%d), data: '%s'}", m.Code, m.Code, m.Data)
}

func (m ClientMessage) Write(w io.Writer) error {
	return writeMessage(message{Code: int(m.Code), Data: m.Data}, w)
}

// ReadAvailableClientMessages returns all client messages that can be read from the reader
func ReadAvailableClientMessages(r io.Reader) ([]ClientMessage, error) {

	msgs, err := readAvailableMessages(r)
	if err != nil {
		return nil, err
	}

	clientMsgs := make([]ClientMessage, len(msgs))
	for i, m := range msgs {
		clientMsgs[i] = ClientMessage{Code: ClientCode(m.Code), Data: m.Data}
	}

	return clientMsgs, nil
}
