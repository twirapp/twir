package voteban

import (
	"context"
	"strings"
	"sync"
	"time"

	votebanentity "github.com/twirapp/twir/libs/entities/voteban"
)

type (
	sessionResult struct {
		channelId    string
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
		s.writeResult()
		result, ok := <-s.result
		return result, ok
	case result, ok := <-s.result:
		return result, ok
	}
}

// writeResult computes sessionResult based on current state of session, writes it to the result channel
// and closes it. The result will be written only once, as the session has only one result.
func (s *session) writeResult() {
	var message string

	if s.isTargetUserModerator {
		message = s.voteban.SurviveMessageModerators
	} else {
		message = s.voteban.SurviveMessage
	}

	s.votesMu.RLock()
	defer s.votesMu.RUnlock()

	// We should ban chatter if there are enough votes and most of them are positive.
	isBan := s.yesVotes > s.noVotes && len(s.votes) >= s.voteban.NeededVotes
	if isBan {
		if s.isTargetUserModerator {
			message = s.voteban.BanMessageModerators
		} else {
			message = s.voteban.BanMessage
		}
	}

	message = strings.ReplaceAll(message, "{targetUser}", s.targetUserLogin)

	s.resultOnce.Do(
		func() {
			s.result <- sessionResult{
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
		},
	)
}

func (s *session) tryRegisterVote(chatterUserId, message string) bool {
	s.votesMu.Lock()
	defer s.votesMu.Unlock()

	// Session has already finished or in process of finishing it.
	if len(s.votes) >= s.voteban.NeededVotes {
		return false
	}

	// Chatter has already voted during this session.
	if _, ok := s.votes[chatterUserId]; ok {
		return false
	}

	lastVotesCount := len(s.votes)
	for word := range strings.FieldsSeq(message) {
		if isPositive, exists := s.words[word]; exists {
			s.votes[chatterUserId] = isPositive

			if isPositive {
				s.yesVotes++
			} else {
				s.noVotes++
			}
			break
		}
	}

	if len(s.votes) >= s.voteban.NeededVotes {
		go s.writeResult()
	}

	isVoteRegistered := len(s.votes) > lastVotesCount
	return isVoteRegistered
}
