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
			msgType := []struct {
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

			for _, mType := range msgType {
				if mData, ok := mType.Data.([]byte); ok {
					switch mType.Action {
					case "channelChatStreamMessage":
						b.readChat(mData)
					case "channelPollStream":
						b.answerPoll(mData)
					default:
						log.WithFields(log.Fields{
							"bot": b.conf.NickName,
							"msg": string(mData),
						}).Warn("unknow message")
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
