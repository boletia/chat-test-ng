package bot

const (
	// DefaultNumBots default bots
	DefaultNumBots = 1
	// DefaultSendMessages false
	DefaultSendMessages = false
	// DefaultNumMessages Number of messages to send
	DefaultNumMessages = 0
	// DefaultMinDelay minumum of latency between messages
	DefaultMinDelay = 10
	// DefaultMaxDelay minumum of latency between messages
	DefaultMaxDelay = 30
	// DefaultWithGossiper default to false
	DefaultWithGossiper = false
	// DefaultSubdomain subdomain to join
	DefaultSubdomain = "rob-test-event"
	// DefautlEndPoint where we have to connect
	DefautlEndPoint = "wss://7qbaj6pufe.execute-api.us-east-1.amazonaws.com/beta"
)

// Conf Depic new bot configuration
type Conf struct {
	NickName     string
	SendMessages bool
	WithGossiper bool
	SudDomain    string
	NumMessages  int
	MinDelay     int
	MaxDelay     int
	URL          string
}

// Socket interface to send/receive messages
type Socket interface {
	Write(msg []byte) error
	Read(*[]byte) error
	CountCalls(*int, *int)
}

type bot struct {
	socket Socket
	conf   Conf
	quit   chan bool
}

// New Creates new bot instance
func New(cnf Conf, quick chan bool) bot {
	return bot{
		socket: nil,
		conf:   cnf,
		quit:   quick,
	}
}
