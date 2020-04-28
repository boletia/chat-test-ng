package bot

import (
	"sync"
	"time"

	"github.com/boletia/chat-test-ng/pkg/wsocket"
	log "github.com/sirupsen/logrus"
)

func (b *bot) connect() bool {
	connected := false
	startTime := time.Now()

	defer func() {
		if connected {
			elapseTime := time.Since(startTime)
			log.WithFields(log.Fields{
				"bot":     b.conf.NickName,
				"seconds": elapseTime.Seconds(),
			}).Info("connection time")
		}
	}()

	wsocket, err := wsocket.New(b.conf.URL)
	if err != nil {
		log.WithFields(log.Fields{
			"bot": b.conf.NickName,
			"url": b.conf.URL,
		}).Error("unable to connect bot")
		return connected
	}

	connected = true
	b.socket = wsocket

	return connected
}

func (b bot) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	if !b.connect() || !b.JoinChat() {
		return
	}

	msg := make(chan []byte)
	go b.readMessage(msg)
	go b.listen(msg)

	if b.conf.SendMessages == true {
		go b.chat()
	}

	for {
		select {
		case <-b.quit:
			return
		}
	}
}
