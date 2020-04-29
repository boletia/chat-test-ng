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

func (b bot) readChat(msg map[string]interface{}) {

	for key, value := range msg {
		switch key {
		case "data":
			if chatData, isChatData := value.(map[string]interface{}); isChatData {
				log.WithFields(log.Fields{
					"bot":  b.conf.NickName,
					"from": chatData["author"],
					"msg":  chatData["message"],
				}).Info("chat read")
			}

			return
		}
	}
}
