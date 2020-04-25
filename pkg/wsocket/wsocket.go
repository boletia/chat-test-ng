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

func (ws wsocket) Write(msg []byte) error {
	ws.mux.Lock()
	// write to the socket
	ws.mux.Unlock()
	return nil
}

func (ws wsocket) Read(*[]byte) error {
	return nil
}
