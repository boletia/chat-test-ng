package bot

import (
	"encoding/json"
	"math/rand"

	log "github.com/sirupsen/logrus"
)

const emitVote = "channelPollUserOnVote"

type pollAnswer struct {
	ID     string `json:"id"`
	Option string `json:"option_label"`
	Total  int    `json:"total"`
}

type pollMessageData struct {
	PollID         string       `json:"id"`
	Name           string       `json:"name"`
	Active         bool         `json:"active"`
	Answers        []pollAnswer `json:"answers"`
	EventSubdomain string       `json:"event_subdomain"`
}

type pollMessage struct {
	Action string          `json:"action"`
	Data   pollMessageData `json:"data"`
}

type pollVoteData struct {
	PollID string `json:"poll_id"`
	Answer string `json:"answer_id"`
}

type pollVote struct {
	Action string       `json:"action"`
	Data   pollVoteData `json:"data"`
}

func (b bot) answerPoll(poll pollMessage) {
	choice := rand.Intn(len(poll.Data.Answers))

	vote := pollVote{
		Action: emitVote,
		Data: pollVoteData{
			PollID: poll.Data.PollID,
			Answer: poll.Data.Answers[choice].ID,
		},
	}

	if voteData, err := json.Marshal(vote); err != nil {
		log.WithFields(log.Fields{
			"bot":   b.conf.NickName,
			"poll":  poll.Data.Name,
			"error": err,
		}).Error("unable to answer poll")
	} else {
		if err := b.socket.Write(voteData); err != nil {
			log.WithFields(log.Fields{
				"bot":   b.conf.NickName,
				"error": err,
			}).Error("socket write error")
		}
	}
}
