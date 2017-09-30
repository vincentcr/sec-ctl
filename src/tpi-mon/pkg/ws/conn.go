package ws

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"tpi-mon/pkg/site"

	"github.com/gorilla/websocket"
)

func init() {
	gob.Register(ControlMessage{})
	gob.Register(site.UserCommand{})
	gob.Register(site.Event{})
	gob.Register(site.Partition{})
	gob.Register(site.Zone{})
	gob.Register(site.StateChange{})
	gob.Register(site.SystemState{})
	gob.Register(site.SystemTroubleStatus(0))
	gob.Register(site.Alarm{})
}

// Conn is a wrapper type of the websocket connection
type Conn struct {
	ws *websocket.Conn
}

// Dial opens a connection to the specified server, using the specified auth token
func Dial(url string, token string) (*Conn, error) {
	var dialer *websocket.Dialer
	authVal := fmt.Sprintf("Bearer %v", token)
	authHead := http.Header{"Authorisation": []string{authVal}}

	conn, _, err := dialer.Dial(url, authHead)
	if err != nil {
		return nil, err
	}

	return &Conn{ws: conn}, nil
}

// UpgradeRequest upgrades an http request connection to a websocket connection
func UpgradeRequest(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	var wsupgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return &Conn{conn}, nil
}

type ControlMessageCode byte
type ControlMessage struct {
	Code ControlMessageCode
}

const (
	CtrlGetState ControlMessageCode = 1
)

func (conn *Conn) Write(data interface{}) error {
	w, err := conn.ws.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(w)
	if err = enc.Encode(&data); err != nil {
		return err
	}

	log.Println("wrote message:", data)

	return w.Close()
}

func (conn *Conn) Read() (interface{}, error) {

	_, r, err := conn.ws.NextReader()
	if err != nil {
		return nil, err
	}

	dec := gob.NewDecoder(r)

	var res interface{}
	if err = dec.Decode(&res); err != nil {
		return nil, err
	}

	log.Println("read message:", res)

	return res, nil
}

// Close closes the connection
func (conn *Conn) Close() error {
	return conn.ws.Close()
}
