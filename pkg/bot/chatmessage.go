package bot

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

//"map[action:channelChatStreamMessage data:map[author:bot-0 event_subdomain:rob-test-event message:Message 1 of 4, latency 8]
type chatStreamMessageData struct {
	Author         string `json:"author"`
	EventSubdomain string `json:"event_subdomain"`
	Message        string `json:"Message"`
}

type chatStreamMessage struct {
	Action string                `json:"action"`
	Data   chatStreamMessageData `json:"data"`
}

func (b bot) readChat(chatMsg []byte) {
	chat := chatStreamMessage{}

	if err := json.Unmarshal(chatMsg, &chat); err != nil {
		log.WithFields(log.Fields{
			"bot":   b.conf.NickName,
			"error": err,
		}).Error("json unmarshaling chat message")
		return
	}

	log.WithFields(log.Fields{
		"bot":  b.conf.NickName,
		"from": chat.Data.Author,
		"msg":  chat.Data.Message,
	}).Info("chat read")
}
