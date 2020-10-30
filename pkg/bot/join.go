package bot

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	joinChatAction   = "channelStreamJoinUser"
	joinChatActionV2 = "channelStreamJoinUserV2"
	joinChatActionV3 = "v1/users/channel/subscribe"
)

type fingerprintJoin struct {
	ID       string `json:"id,ommitempty"`
	OS       string `json:"os,ommitempty"`
	Browser  string `json:"browser,ommitempty"`
	Location string `json:"location,ommitempty"`
	IP       string `json:"ip,ommitempty"`
}

type joinChannelData struct {
	EventSubdomain    string           `json:"event_subdomain"`
	IsOrganizer       bool             `json:"is_organizer"`
	NickName          string           `json:"nickname"`
	TemporaryID       *string          `json:"temporary_id,ommitempty"`
	HasCustomStickers *bool            `json:"has_custom_stickers,ommitempty"`
	Fingerprint       *fingerprintJoin `json:"fingerprint,ommitempty"`
}

type joinChannel struct {
	Action string          `json:"action"`
	Data   joinChannelData `json:"data"`
}

func (b bot) JoinChat() bool {
	joinChat := joinChannel{
		Action: joinChatActionV3,
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
		return false
	}

	nano := fmt.Sprintf("%d", time.Now().UnixNano())
	tempID := string(nano[len(nano)-7:])

	fprint := &fingerprintJoin{
		ID:       tempID,
		OS:       "Mac OS 10.15.7",
		Browser:  "Chrome 86.0.4240.111",
		Location: "cant tell you",
		IP:       "cant tell you",
	}

	hasSticker := false

	joinChat.Data.TemporaryID = &tempID
	joinChat.Data.HasCustomStickers = &hasSticker
	joinChat.Data.Fingerprint = fprint

	msgByte, err = json.Marshal(joinChat)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"bot":   b.conf.NickName,
		}).Error("joinchat marshal error")
		return false
	}

	if err := b.apiGateWaySocket.Write(msgByte); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"bot":   b.conf.NickName,
		}).Error("unable to join to chat")
		return false
	}

	return true
}
