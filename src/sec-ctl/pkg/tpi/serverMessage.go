package tpi

import (
	"fmt"
	"io"
)

// ServerMessage represents a message sent by the server
type ServerMessage struct {
	Code ServerCode
	Data []byte
}

func (m ServerMessage) Write(w io.Writer) error {
	return writeMessage(message{Code: int(m.Code), Data: m.Data}, w)
}

func (m ServerMessage) String() string {
	return fmt.Sprintf("ServerMessage{code: %s(%d), data: '%s'}", m.Code.Name(), m.Code, m.Data)
}

// ReadAvailableServerMessages returns all server messages that can be read from the reader
func ReadAvailableServerMessages(r io.Reader) ([]ServerMessage, error) {

	msgs, err := readAvailableMessages(r)
	if err != nil {
		return nil, err
	}

	serverMsgs := make([]ServerMessage, len(msgs))
	for i, m := range msgs {
		serverMsgs[i] = ServerMessage{Code: ServerCode(m.Code), Data: m.Data}
	}

	return serverMsgs, nil

}
