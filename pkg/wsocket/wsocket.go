package wsocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type wsocket struct {
	*websocket.Conn
	mux sync.Mutex
}

// New Creates new websocket
func New(url string) (*wsocket, error) {
	conn, err := connect(url)
	if err != nil {
		return nil, err
	}

	return &wsocket{
		conn,
		sync.Mutex{},
	}, nil
}

func connect(url string) (*websocket.Conn, error) {
	var conn *websocket.Conn
	var err error

	if conn, _, err = websocket.DefaultDialer.Dial(url, nil); err != nil {
		return nil, err
	}
	return conn, nil
}

func (ws *wsocket) Write(msg []byte) error {
	var err error

	ws.mux.Lock()
	err = ws.WriteMessage(websocket.TextMessage, msg)
	ws.mux.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (ws *wsocket) Read(data *[]byte) error {
	var err error
	var msgType int

	for {
		msgType, *data, err = ws.Conn.ReadMessage()

		if msgType != websocket.TextMessage {
			continue
		}

		if err != nil {
			return err
		}

		if len(*data) > 0 {
			return nil
		}
	}
}

func (ws *wsocket) Disconnect() {
	ws.Conn.Close()
}
