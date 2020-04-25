package bot

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

const (
	joinChatAction   = "channelStreamJoinUser"
	joinChatActionV2 = "channelStreamJoinUserV2"
)

type joinChannelData struct {
	EventSubdomain string `json:"event_subdomain"`
	IsOrganizer    bool   `json:"is_organizer"`
	NickName       string `json:"nickname"`
}

type joinChannel struct {
	Action string          `json:"action"`
	Data   joinChannelData `json:"data"`
}

func (b bot) JoinChat() bool {
	joinChat := joinChannel{
		Action: joinChatActionV2,
		Data: joinChannelData{
			EventSubdomain: b.conf.SudDomain,
			NickName:       b.conf.NickName,
			IsOrganizer:    false,
		},
	}

	msgByte, err := json.Marshal(joinChat)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"bot":   b.conf.NickName,
		}).Error("joinchat marshal error")
		return false
	}

	if err := b.socket.Write(msgByte); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"bot":   b.conf.NickName,
		}).Error("unable to join to chat")
	}
	return true
}
