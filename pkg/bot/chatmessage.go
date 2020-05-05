package bot

import (
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

func (b bot) readChat(chatMsgIface interface{}) {

	if chatMsg, ok := chatMsgIface.([]chatStreamMessageData); ok {
		for _, chat := range chatMsg {
			log.WithFields(log.Fields{
				"bot":  b.conf.NickName,
				"from": chat.Author,
				"msg":  chat.Message,
			}).Info("chat read")
		}
	} else {
		log.WithFields(log.Fields{
			"bot":  b.conf.NickName,
			"data": chatMsgIface,
		}).Error("chat read error")
	}
}
