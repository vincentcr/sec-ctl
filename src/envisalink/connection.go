package envisalink

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

func Demo() {
	servAddr := "192.168.86.12:4025"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("connected to %s\n", servAddr)

	defer conn.Close()

	ch := make(chan ServerMessage)

	startReadLoop(conn, ch)

	for {
		select {
		case msg := <-ch:
			processServerMessage(conn, msg)

		}
	}

	// loginMsg := ClientMessage{
	// 	Code: NetworkLogin,
	// 	Data: []byte("Q4m1gh"),
	// }

	// loginBytes := msgEncode(loginMsg)

	// fmt.Printf("msg %v encoded to %v\n", loginMsg, loginBytes)

	// _, err = conn.Write(loginBytes)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("wrote msg to socket")

	// for {
	// 	data, err := readUntilBytes(conn, crlf)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Printf("read data: %v\n", data)

	// 	msgPackets := bytes.Split(data, crlf)
	// 	for i, packet := range msgPackets {
	// 		fmt.Printf("decoding packet: %v\n", packet)
	// 		if len(packet) > 0 {
	// 			msg, err := msgDecode(packet)
	// 			if err != nil {
	// 				panic(err)
	// 			}

	// 			fmt.Printf("Decoded msg %d: %v\n", i, msg)

	// 		}
	// 	}

	// }
}

func processServerMessage(conn io.Writer, msg ServerMessage) {
	fmt.Printf("received message {code:%d, data:'%s'}\n", msg.Code, msg.Data)
	switch msg.Code {
	case LoginRes:
		if string(msg.Data) != "1" {
			loginMsg := ClientMessage{
				Code: NetworkLogin,
				Data: []byte("Q4m1gh"),
			}
			if err := writeMessage(conn, loginMsg); err != nil {
				panic(err)
			}
		}
	}
}

func writeMessage(writer io.Writer, msg ClientMessage) error {
	fmt.Printf("sending message {code:%d, data:'%s'}\n", msg.Code, msg.Data)
	msgBytes := msgEncode(msg)

	fmt.Printf("msg %v encoded to %v\n", msg, msgBytes)

	_, err := writer.Write(msgBytes)
	return err
}

func startReadLoop(reader io.Reader, ch chan ServerMessage) {
	go func() {
		for {
			msgs, err := readMessages(reader)
			for _, msg := range msgs {
				ch <- msg
			}
			if err != nil {
				panic(err)
			}
		}
	}()
}

func readMessages(reader io.Reader) ([]ServerMessage, error) {
	packetBytes, err := readUntilMarker(reader, crlf)
	if err != nil {
		return nil, err
	}
	packets := bytes.Split(packetBytes, crlf)
	msgs := make([]ServerMessage, len(packets))
	for i, packet := range packets {
		if len(packet) > 0 {
			msg, err := msgDecode(packet)
			if err != nil {
				return nil, err
			}
			msgs[i] = msg
		}
	}
	return msgs, nil
}

func readUntilMarker(reader io.Reader, marker []byte) ([]byte, error) {

	data := make([]byte, 0, 4096)
	buf := make([]byte, 2048)
	done := false

	for !done {
		fmt.Println("reading next buffer")
		nRead, err := reader.Read(buf)
		if err != nil {
			return nil, err
		}
		if nRead == 0 {
			return nil, fmt.Errorf("Unexpected end of input")
		}
		data = append(data, buf[:nRead]...)

		// done when <marker bytes> are the last bytes of the transmission
		done = bytes.Compare(marker, data[len(data)-len(marker):]) == 0
	}

	return data, nil
}
