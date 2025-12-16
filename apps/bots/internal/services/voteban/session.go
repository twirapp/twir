package voteban

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/twirapp/twir/libs/bus-core/twitch"
	votebanentity "github.com/twirapp/twir/libs/entities/voteban"
)

type votebanSession struct {
	data     votebanentity.Voteban
	yesWords map[string]struct{}
	noWords  map[string]struct{}
	votes    map[string]bool

	targerUserId          string
	targerUserLogin       string
	isTargerUserModerator bool

	mu         sync.Mutex
	stopFunc   context.CancelFunc
	finishChan chan sessionResult
}

func createVotebanSession(data votebanentity.Voteban, initiatorUserId, targerUserId, targerUserLogin string, isTargerUserModerator bool) *votebanSession {
	s := votebanSession{
		data:     data,
		mu:       sync.Mutex{},
		yesWords: map[string]struct{}{},
		noWords:  map[string]struct{}{},
		votes: map[string]bool{
			targerUserId: true,
		},
		targerUserId:          targerUserId,
		targerUserLogin:       targerUserLogin,
		isTargerUserModerator: isTargerUserModerator,
	}

	for _, w := range data.ChatVotesWordsPositive {
		s.yesWords[w] = struct{}{}
	}

	for _, w := range data.ChatVotesWordsNegative {
		s.noWords[w] = struct{}{}
	}

	return &s
}

type sessionResult struct {
	channelId    string
	haveDecision bool
	isModerator  bool
	isBan        bool
	targerUserId string
	message      string
	yesVotes     int
	noVotes      int
	banDuration  int
}

func (c *votebanSession) start() chan sessionResult {
	c.finishChan = make(chan sessionResult)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(c.data.VoteDuration)*time.Second,
	)

	c.stopFunc = cancel

	go func() {
		<-ctx.Done()
		c.finish()
		close(c.finishChan)
	}()

	return c.finishChan
}

func (c *votebanSession) finish() {
	defer c.stopFunc()

	if len(c.votes) < c.data.NeededVotes {
		c.finishChan <- sessionResult{
			haveDecision: false,
		}
		return
	}

	var (
		yesVotes int
		noVotes  int
	)

	for _, v := range c.votes {
		if v {
			yesVotes++
		} else {
			noVotes++
		}
	}

	var (
		message string
		isBan   bool
	)

	if yesVotes <= noVotes {
		if c.isTargerUserModerator {
			message = c.data.SurviveMessageModerators
		} else {
			message = c.data.SurviveMessage
		}
	} else if yesVotes > noVotes {
		isBan = true
		if c.isTargerUserModerator {
			message = c.data.BanMessageModerators
		} else {
			message = c.data.BanMessage
		}
	}

	message = strings.ReplaceAll(message, "{targetUser}", c.targerUserLogin)

	c.finishChan <- sessionResult{
		haveDecision: true,
		isModerator:  c.isTargerUserModerator,
		isBan:        isBan,
		targerUserId: c.targerUserId,
		message:      message,
		yesVotes:     yesVotes,
		noVotes:      noVotes,
		channelId:    c.data.ChannelID,
		banDuration:  c.data.TimeoutSeconds,
	}
}

func (c *votebanSession) tryRegisterVote(msg twitch.TwitchChatMessage) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if msg.Message == nil {
		return
	}

	if len(c.votes) >= c.data.NeededVotes {
		return
	}

	if _, ok := c.votes[msg.ChatterUserId]; ok {
		return
	}

	for part := range strings.FieldsSeq(msg.Message.Text) {
		if _, ok := c.yesWords[part]; ok {
			c.votes[msg.ChatterUserId] = true
			break
		}

		if _, ok := c.noWords[part]; ok {
			c.votes[msg.ChatterUserId] = false
			break
		}
	}

	if len(c.votes) >= c.data.NeededVotes {
		c.finish()
	}
}
