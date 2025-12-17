package voteban

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/twirapp/twir/libs/bus-core/twitch"
	votebanentity "github.com/twirapp/twir/libs/entities/voteban"
)

type (
	sessionResult struct {
		channelId    string
		haveDecision bool
		isModerator  bool
		isBan        bool
		targetUserId string
		message      string
		yesVotes     int
		noVotes      int
		banDuration  int
	}

	session struct {
		words                 map[string]bool
		votes                 map[string]bool
		yesVotes              int
		noVotes               int
		votesMu               sync.Mutex
		result                chan sessionResult
		voteban               votebanentity.Voteban
		targetUserId          string
		targetUserLogin       string
		isTargetUserModerator bool
	}
)

func newSession(
	voteban votebanentity.Voteban,
	targetUserId string,
	targetUserLogin string,
	isTargetUserModerator bool,
) *session {
	sess := session{
		words:  make(map[string]bool),
		votes:  make(map[string]bool),
		result: make(chan sessionResult),

		voteban:               voteban,
		targetUserId:          targetUserId,
		targetUserLogin:       targetUserLogin,
		isTargetUserModerator: isTargetUserModerator,
	}

	for _, word := range voteban.ChatVotesWordsPositive {
		sess.words[word] = true
	}

	for _, word := range voteban.ChatVotesWordsNegative {
		sess.words[word] = false
	}

	sess.votes[targetUserId] = true
	return &sess
}

func (s *session) waitResult() (sessionResult, bool) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(s.voteban.VoteDuration)*time.Second,
	)
	defer cancel()

	select {
	case <-ctx.Done():
		s.writeResult()
		result := <-s.result
		return result, true
	case result := <-s.result:
		return result, true
	}
}

func (s *session) writeResult() {
	s.votesMu.Lock()
	defer s.votesMu.Unlock()

	if len(s.votes) < s.voteban.NeededVotes {
		s.result <- sessionResult{
			haveDecision: false,
		}
		return
	}

	var (
		message string
		isBan   = s.yesVotes > s.noVotes
	)

	if s.isTargetUserModerator {
		message = s.voteban.SurviveMessageModerators
	} else {
		message = s.voteban.SurviveMessage
	}

	if isBan {
		if s.isTargetUserModerator {
			message = s.voteban.BanMessageModerators
		} else {
			message = s.voteban.BanMessage
		}
	}

	message = strings.ReplaceAll(message, "{targetUser}", s.targetUserLogin)

	s.result <- sessionResult{
		haveDecision: true,
		isModerator:  s.isTargetUserModerator,
		isBan:        isBan,
		targetUserId: s.targetUserId,
		message:      message,
		yesVotes:     s.yesVotes,
		noVotes:      s.noVotes,
		channelId:    s.voteban.ChannelID,
		banDuration:  s.voteban.TimeoutSeconds,
	}
}

func (s *session) tryRegisterVote(msg twitch.TwitchChatMessage) {
	if msg.Message == nil {
		return
	}

	s.votesMu.Lock()
	defer s.votesMu.Unlock()

	// Session has already finished or in process of finishing it.
	if len(s.votes) >= s.voteban.NeededVotes {
		return
	}

	// Chatter has already voted during this session.
	if _, ok := s.votes[msg.ChatterUserId]; ok {
		return
	}

	for word := range strings.FieldsSeq(msg.Message.Text) {
		if isYay, exists := s.words[word]; exists {
			s.votes[msg.ChatterUserId] = isYay

			if isYay {
				s.yesVotes++
			} else {
				s.noVotes++
			}
		}
	}

	if len(s.votes) >= s.voteban.NeededVotes {
		s.writeResult()
	}
}
