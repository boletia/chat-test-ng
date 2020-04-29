package bot

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

func (b bot) listen(msg chan []byte) {
	for {
		var data []byte
		if err := b.socket.Read(&data); err != nil {
			log.WithFields(log.Fields{
				"bot":   b.conf.NickName,
				"error": err,
			}).Error("socket read error")
			return
		}

		msg <- data
	}
}

func (b bot) readMessage(msg chan []byte) {
	for {
		select {
		case data := <-msg:
			var receivedMessage interface{}
			if err := json.Unmarshal(data, &receivedMessage); err != nil {
				log.WithFields(log.Fields{
					"bot":   b.conf.NickName,
					"error": err,
				}).Error("json unmarshal")
				break
			}

			if data, isValid := receivedMessage.(map[string]interface{}); isValid {
				for key, value := range data {
					if key == "action" {
						switch value {
						case "channelChatStreamMessage":
							b.readChat(data)
						default:
							log.WithFields(log.Fields{
								"bot":    b.conf.NickName,
								"action": value,
							}).Warn("read unknow message")
						}
					}
				}
			}
			/*
				switch msgType := receivedMessage.(type) {
				case pollMessage:
					b.answerPoll(msgType)
				case chatStreamMessage:
					b.readChat(msgType)
				default:
					log.WithFields(log.Fields{
						"bot":  b.conf.NickName,
						"type": fmt.Sprintf("%T", msgType),
						"data": fmt.Sprintf("%#v", msgType),
					}).Warn("read unknow message")
				}
			*/
		}
	}
}
