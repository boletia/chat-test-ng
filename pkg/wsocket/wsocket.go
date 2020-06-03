package wsocket

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type wsocket struct {
	*websocket.Conn
	mux        sync.Mutex
	writtenOps int
	readOps    int
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
		0,
		0,
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

	ws.writtenOps++
	return nil
}

func (ws *wsocket) Read(data *[]byte) error {
	var err error
	var msgType int

	for {
		msgType, *data, err = ws.Conn.ReadMessage()
		if err != nil {
			return err
		}

		if msgType != websocket.TextMessage {
			continue
		}

		if len(*data) > 0 {
			ws.readOps++
			return nil
		}
	}
}

func (ws *wsocket) CountCalls(written *int, read *int) {
	*written = ws.writtenOps
	*read = ws.readOps
}

func (ws *wsocket) SendCloseMessage(deadLine time.Time) error {
	err := ws.WriteControl(websocket.CloseMessage, nil, deadLine)

	return err
}

func (ws *wsocket) CloseSocket() bool {
	return (ws.Conn.Close() == nil)
}

func (ws *wsocket) GetSocket() *websocket.Conn {
	return ws.Conn
}
