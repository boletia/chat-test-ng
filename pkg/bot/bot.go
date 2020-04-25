package bot

// Conf Depic new bot configuration
type Conf struct {
	NickName     string
	SendMessages bool
	WithGossiper bool
	SudDomain    string
	NumMessages  int
	MinDelay     int64
	MaxDelay     int64
	URL          string
}

// Socket interface to send/receive messages
type Socket interface {
	Write(msg []byte) error
	Read(*[]byte) error
}

type bot struct {
	socket Socket
	conf   Conf
	quit   chan bool
}

// New Creates new bot instance
func New(cnf Conf, sock Socket, quick chan bool) bot {
	return bot{
		socket: nil,
		conf:   cnf,
		quit:   quick,
	}
}
