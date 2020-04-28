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

			switch msgType := receivedMessage.(type) {
			case pollMessage:
				b.answerPoll(msgType)
			default:
				log.WithFields(log.Fields{
					"bot":     b.conf.NickName,
					"message": msgType,
				}).Warn("read unknow message")
			}
		}
	}
}
