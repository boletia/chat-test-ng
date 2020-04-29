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

type chatMessageData struct {
	NickName       string `json:"nickname"`
	Message        string `json:"message"`
	EventSubdomain string `json:"event_subdomain"`
}

type chatMessage struct {
	Action string          `json:"action"`
	Data   chatMessageData `json:"data"`
}

func (b bot) chat() {
	for msgCount := 1; msgCount <= b.conf.NumMessages; msgCount++ {

		latency := rand.Intn(b.conf.MaxDelay-b.conf.MinDelay) + b.conf.MinDelay
		time.Sleep(time.Duration(latency) * time.Second)

		msg := chatMessage{
			Action: chatActionV1,
			Data: chatMessageData{
				NickName:       b.conf.NickName,
				Message:        fmt.Sprintf("Message %d of %d, latency %d", msgCount, b.conf.NumMessages, latency),
				EventSubdomain: b.conf.SudDomain,
			},
		}

		if msgByte, err := json.Marshal(msg); err == nil {
			if err := b.socket.Write(msgByte); err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"bot":   b.conf.NickName,
				}).Error("unable to send chat message")
			}
		} else {
			log.WithFields(log.Fields{
				"error": err,
				"bot":   b.conf.NickName,
			}).Error("msg chat marshaling error")
		}

	}

	log.WithFields(log.Fields{
		"bot": b.conf.NickName,
	}).Info("has send all its messages")
}
