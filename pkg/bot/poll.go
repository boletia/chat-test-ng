package bot

import (
	"encoding/json"
	"math/rand"

	log "github.com/sirupsen/logrus"
)

const emitVote = "v1/poll/vote"

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

func (b bot) answerPoll(pollMsg []byte) {
	poll := pollMessage{}

	if err := json.Unmarshal(pollMsg, &poll); err != nil {
		log.WithFields(log.Fields{
			"bot":   b.conf.NickName,
			"error": err,
		}).Error("json unmarshaling poll message")
		return
	}

	log.WithFields(log.Fields{
		"bot":     b.conf.NickName,
		"poll":    poll.Data.Name,
		"options": poll.Data.Answers,
	}).Info("poll message")

	optionAnswer := rand.Intn(len(poll.Data.Answers))

	vote := pollVote{
		Action: emitVote,
		Data: pollVoteData{
			PollID: poll.Data.PollID,
			Answer: poll.Data.Answers[optionAnswer].ID,
		},
	}

	if voteMsg, err := json.Marshal(vote); err != nil {
		log.WithFields(log.Fields{
			"bot":   b.conf.NickName,
			"poll":  poll.Data.Name,
			"error": err,
		}).Error("json marshaling pollvote message")
	} else {
		if b.socket.Write(voteMsg) != nil {
			log.WithFields(log.Fields{
				"bot":   b.conf.NickName,
				"poll":  poll.Data.Name,
				"error": err,
			}).Error("unable to emmit poll vote message")
		} else {
			log.WithFields(log.Fields{
				"bot":    b.conf.NickName,
				"poll":   poll.Data.Name,
				"answer": poll.Data.Answers[optionAnswer].Option,
			}).Info("poll answer send")

		}
	}
}
