package bot

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// DefaultNumBots default bots
	DefaultNumBots = 1
	// DefaultSendMessages false
	DefaultSendMessages = false
	// DefaultNumMessages Number of messages to send
	DefaultNumMessages = 0
	// DefaultMinDelay minumum of latency between messages
	DefaultMinDelay = 1000
	// DefaultMaxDelay minumum of latency between messages
	DefaultMaxDelay = 5000
	// DefaultWithGossiper default to false
	DefaultWithGossiper = false
	// DefaultSubdomain subdomain to join
	DefaultSubdomain = "neermetv-test"
	// DefautlEndPoint where we have to connect
	DefautlEndPoint = "wss://chat.neerme.io:81/"
	// DefaultWsapiEndPoint endpoint to send second join
	DefaultWsapiEndPoint = "wss://wsapi.neerme.io/"
	// DefaultRamping number of bots/sec
	DefaultRamping = 10
	// DefaultSecondsToReport sleep time before reports
	DefaultSecondsToReport = 10
)

// Conf Depic new bot configuration
type Conf struct {
	NickName        string
	SendMessages    bool
	WithGossiper    bool
	SudDomain       string
	NumMessages     int
	MinDelay        int64
	MaxDelay        int64
	URL             string
	URLWsapi        string
	Ramping         int
	OnlyError       bool
	Sent2Dynamo     bool
	SecondsToReport uint64
}

type count struct {
	read bool
}

// Socket interface to send/receive messages
type Socket interface {
	Write(msg []byte) error
	Read(*[]byte) error
	CountCalls(*int, *int)
	SendCloseMessage(time.Time) error
	CloseSocket() bool
}

// Dynamo interface to send messages to dynamo
type Dynamo interface {
	Write(msg []byte) error
}

type bot struct {
	socket           Socket
	apiGateWaySocket Socket
	dynamo           Dynamo
	conf             Conf
	quit             chan bool
	rcvdMessages     *uint64
	sendMessages     *uint64
	summaryCount     *uint64
}

// New Creates new bot instance
func New(cnf Conf, quick chan bool) bot {
	return bot{
		socket:           nil,
		apiGateWaySocket: nil,
		dynamo:           nil,
		conf:             cnf,
		quit:             quick,
		rcvdMessages:     new(uint64),
		sendMessages:     new(uint64),
		summaryCount:     new(uint64),
	}
}

func (b *bot) AddDynamo(dy Dynamo) {
	b.dynamo = dy
}

func (b bot) messageCounter(event chan count) {
	for {
		select {
		case msg := <-event:
			if msg.read {
				*b.rcvdMessages++
			} else {
				*b.sendMessages++
			}
		}
	}
}

func (b bot) printMsgsSumary() {
	defer func() {
		b.quit <- true
	}()

	for {
		select {
		case <-b.quit:
			return
		default:
			time.Sleep(time.Second * time.Duration(b.conf.SecondsToReport))
			log.WithFields(log.Fields{
				"bot":           b.conf.NickName,
				"rcvd_messages": *b.rcvdMessages,
				"sent_messages": fmt.Sprintf("%d/%d", *b.sendMessages, b.conf.NumMessages),
				"sleep_time":    b.conf.SecondsToReport,
				"iteration":     *b.summaryCount,
			}).Info("msg received")
			*b.summaryCount++
		}
	}
}
