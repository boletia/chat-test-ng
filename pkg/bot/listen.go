package bot

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func (b bot) listen(msg chan []byte, counter chan count) {
	sock := b.socket.GetSocket()

	for {
		data := make([]byte, 0)
		msgType, reader, err := sock.NextReader()
		if err != nil {
			return
		}

		if msgType != websocket.TextMessage {
			continue
		}

		buf := make([]byte, 0)
		for {
			_, err := reader.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				return
			}
			data = append(data, buf[:]...)
		}

		if b.conf.Decode {
			msg <- data
		} else {
			b.msgFile.WriteString(fmt.Sprintf("%s:%s\n", b.conf.NickName, string(data)))
		}

		counter <- count{
			read: true,
		}
		data = nil
	}
}

func (b bot) readMessage(msg chan []byte) {
	for {
		select {
		case data := <-msg:
			msgType := struct {
				Action string                  `json:"action"`
				Data   []chatStreamMessageData `json:"data"`
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
				b.readChat(msgType.Data)
			case "channelPollStream":
				b.answerPoll(data)
			default:
				log.WithFields(log.Fields{
					"bot": b.conf.NickName,
					"msg": string(data),
				}).Warn("unknow message")
			}
		}
	}
}
