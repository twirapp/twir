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
		votesMu               sync.RWMutex
		result                chan sessionResult
		resultOnce            sync.Once
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
		words:                 make(map[string]bool),
		votes:                 make(map[string]bool),
		result:                make(chan sessionResult, 1),
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
		s.writeResultOnce()
		result := <-s.result
		return result, true
	case result := <-s.result:
		return result, true
	}
}

func (s *session) writeResultOnce() {
	s.resultOnce.Do(
		func() {
			s.writeResult()
		},
	)
}

func (s *session) writeResult() {
	s.votesMu.Lock()
	defer s.votesMu.Unlock()

	if len(s.votes) < s.voteban.NeededVotes {
		s.result <- sessionResult{
			haveDecision: false,
		}
		close(s.result)
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
	close(s.result)
}

func (s *session) tryRegisterVote(msg twitch.TwitchChatMessage) {
	if msg.Message == nil {
		return
	}

	s.votesMu.RLock()
	if !s.canVote(msg.ChatterUserId) {
		s.votesMu.RUnlock()
		return
	}
	s.votesMu.RUnlock()

	for word := range strings.FieldsSeq(msg.Message.Text) {
		if isPositive, exists := s.words[word]; exists {
			s.votesMu.Lock()
			if !s.canVote(msg.ChatterUserId) {
				s.votesMu.Unlock()
				break
			}

			s.votes[msg.ChatterUserId] = isPositive

			if isPositive {
				s.yesVotes++
			} else {
				s.noVotes++
			}
			s.votesMu.Unlock()
		}
	}

	s.votesMu.RLock()
	if len(s.votes) >= s.voteban.NeededVotes {
		s.votesMu.RUnlock()
		s.writeResultOnce()
		return
	}
	s.votesMu.RUnlock()
}

func (s *session) canVote(chatterUserId string) bool {
	// Session has already finished or in process of finishing it.
	if len(s.votes) >= s.voteban.NeededVotes {
		return false
	}

	// Chatter has already voted during this session.
	if _, ok := s.votes[chatterUserId]; ok {
		return false
	}

	return true
}
