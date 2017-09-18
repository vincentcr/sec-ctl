package proto

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var crlf = []byte("\r\n")

// TPIMessage represents a TPI server or client message
type TPIMessage struct {
	Code int
	Data []byte
}

func (msg TPIMessage) Write(w io.Writer) error {
	msgBytes := msg.encode()
	_, err := w.Write(msgBytes)
	return err
}

func (msg TPIMessage) encode() []byte {
	encoded := EncodeIntCode(msg.Code)

	if msg.Data != nil {
		encoded = append(encoded, msg.Data...)
	}

	checksum := msgChecksum(encoded)

	return append(append(encoded, []byte(checksum)...), crlf...)
}

// ReadAvailableMessages reads all available messages from the supplied reader
func ReadAvailableMessages(reader io.Reader) ([]TPIMessage, error) {

	packetBytes, err := readUntilMarker(reader, crlf)
	if err != nil {
		return nil, err
	}
	packets := bytes.Split(packetBytes, crlf)
	msgs := make([]TPIMessage, 0, len(packets))
	for _, packet := range packets {
		if len(packet) > 0 {
			msg, err := msgDecode(packet)
			if err != nil {
				return nil, err
			}
			// fmt.Printf("Received %v (raw: %v)\n", msg, packet)
			msgs = append(msgs, msg)
		}
	}
	return msgs, nil
}

func readUntilMarker(reader io.Reader, marker []byte) ([]byte, error) {

	data := make([]byte, 0, 4096)
	buf := make([]byte, 2048)
	done := false

	for !done {
		nRead, err := reader.Read(buf)
		if err != nil {
			return nil, err
		}
		if nRead == 0 {
			return nil, fmt.Errorf("Unexpected end of input")
		}
		data = append(data, buf[:nRead]...)

		// done when <marker bytes> are the last bytes of the transmission
		potentialMarker := data[len(data)-len(marker):]
		done = bytes.Compare(marker, potentialMarker) == 0
	}

	return data, nil
}

func msgDecode(msgBytes []byte) (TPIMessage, error) {
	if len(msgBytes) < 5 {
		return TPIMessage{}, fmt.Errorf("Got %d bytes, need at least 5", len(msgBytes))
	}

	dataStart := 3
	dataEnd := len(msgBytes) - 2
	codeBytes := msgBytes[:dataStart]
	data := msgBytes[dataStart:dataEnd]
	expectedChecksum := string(msgBytes[dataEnd:])
	actualChecksum := msgChecksum(msgBytes[:dataEnd])
	if strings.ToLower(expectedChecksum) != strings.ToLower(actualChecksum) {
		return TPIMessage{}, fmt.Errorf("failed to decode message %v: data %v, expected checksum %v, actual %v",
			msgBytes, data, expectedChecksum, actualChecksum)
	}

	code, err := DecodeIntCode(codeBytes)
	if err != nil {
		return TPIMessage{}, err
	}

	msg := TPIMessage{
		Code: code,
		Data: data,
	}

	return msg, nil
}

func msgChecksum(bytes []byte) string {
	var sum int
	for _, b := range bytes {
		sum += int(b)
	}

	sum = sum & 0xff
	checksum := fmt.Sprintf("%02X", sum)
	return checksum
}

// EncodeIntCode encodes an integer as a tpi code
func EncodeIntCode(code int) []byte {
	return []byte(fmt.Sprintf("%03d", code))
}

// DecodeIntCode parses a byte array as an integer
func DecodeIntCode(codeBytes []byte) (int, error) {
	codeInt, err := strconv.Atoi(string(codeBytes))
	if err != nil {
		return -1, fmt.Errorf("invalid code bytes %s", codeBytes)
	}
	return codeInt, nil
}
