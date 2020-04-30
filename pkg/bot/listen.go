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
			msgType := struct {
				Action string      `json:"action"`
				Data   interface{} `json:"-"`
			}{}

			if err := json.Unmarshal(data, &msgType); err != nil {
				log.WithFields(log.Fields{
					"bot":   b.conf.NickName,
					"error": err,
				}).Error("json unmarshal")
				break
			}

			switch msgType.Action {
			case "channelChatStreamMessage":
				b.readChat(data)
			case "channelPollStream":
				b.answerPoll(data)
			default:
				log.WithFields(log.Fields{
					"bot": b.conf.NickName,
					"msg": string(data),
				}).Warn("unknow message")
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
