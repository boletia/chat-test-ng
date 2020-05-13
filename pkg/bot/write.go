package bot

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	chatAction   = "channelChatUserOnMessage"
	chatActionV1 = "v1/chat/message"
)

type handlerWritter func(msg []byte) error

type writter interface {
	Write(msg []byte) error
}

type chatMessageData struct {
	NickName       string `json:"nickname"`
	Message        string `json:"message"`
	EventSubdomain string `json:"event_subdomain"`
	Avatar         string `json:"avatar"`
}

type chatMessage struct {
	Action string          `json:"action"`
	Data   chatMessageData `json:"data"`
}

func (b bot) chat() {
	var fWritter handlerWritter

	if b.conf.Sent2Dynamo {
		fWritter = b.dynamo.Write
	} else {
		fWritter = b.socket.Write
	}

	timestamp := time.Now()
	for msgCount := 1; msgCount <= b.conf.NumMessages; msgCount++ {

		latency := rand.Int63n(b.conf.MaxDelay-b.conf.MinDelay) + b.conf.MinDelay
		time.Sleep(time.Duration(latency) * time.Millisecond)

		msg := chatMessage{
			Action: chatActionV1,
			Data: chatMessageData{
				NickName:       b.conf.NickName,
				Message:        fmt.Sprintf("Message %d of %d, latency %d", msgCount, b.conf.NumMessages, latency),
				EventSubdomain: b.conf.SudDomain,
				Avatar:         "bot-1",
			},
		}

		if msgByte, err := json.Marshal(msg); err == nil {
			if !b.conf.Sent2Dynamo {
				if err := fWritter(msgByte); err != nil {
					log.WithFields(log.Fields{
						"error": err,
						"bot":   b.conf.NickName,
					}).Error("unable to send message")
				}
			} else {
				go func() {
					if err := fWritter(msgByte); err != nil {
						log.WithFields(log.Fields{
							"error": err,
							"bot":   b.conf.NickName,
						}).Error("unable to send message")
					}
				}()
			}
		} else {
			log.WithFields(log.Fields{
				"error": err,
				"bot":   b.conf.NickName,
			}).Error("msg chat marshaling error")
		}

	}

	log.WithFields(log.Fields{
		"bot":        b.conf.NickName,
		"epalseTime": time.Since(timestamp),
	}).Info("has send all its messages")
}

func (b bot) writeMessage(w writter, msg []byte) {
	if err := w.Write(msg); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"bot":   b.conf.NickName,
		}).Error("unable to send chat message")
	}
}
