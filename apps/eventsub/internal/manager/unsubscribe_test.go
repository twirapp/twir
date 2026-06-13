package manager

import (
	"errors"
	"testing"

	"github.com/nicklaw5/helix/v2"
)

func TestShouldUnsubscribeChannelSubscription(t *testing.T) {
	tests := []struct {
		name        string
		sub         helix.EventSubSubscription
		channelID   string
		shouldMatch bool
	}{
		{
			name:        "matches broadcaster user id",
			sub:         helix.EventSubSubscription{Condition: helix.EventSubCondition{BroadcasterUserID: "123"}},
			channelID:   "123",
			shouldMatch: true,
		},
		{
			name:        "matches moderator user id",
			sub:         helix.EventSubSubscription{Condition: helix.EventSubCondition{ModeratorUserID: "123"}},
			channelID:   "123",
			shouldMatch: true,
		},
		{
			name:        "matches user id",
			sub:         helix.EventSubSubscription{Condition: helix.EventSubCondition{UserID: "123"}},
			channelID:   "123",
			shouldMatch: true,
		},
		{
			name:        "matches to broadcaster user id",
			sub:         helix.EventSubSubscription{Condition: helix.EventSubCondition{ToBroadcasterUserID: "123"}},
			channelID:   "123",
			shouldMatch: true,
		},
		{
			name:        "matches from broadcaster user id",
			sub:         helix.EventSubSubscription{Condition: helix.EventSubCondition{FromBroadcasterUserID: "123"}},
			channelID:   "123",
			shouldMatch: true,
		},
		{
			name:        "does not match different channel",
			sub:         helix.EventSubSubscription{Condition: helix.EventSubCondition{BroadcasterUserID: "999"}},
			channelID:   "123",
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldUnsubscribeChannelSubscription(tt.channelID, tt.sub)
			if got != tt.shouldMatch {
				t.Fatalf("expected %v, got %v", tt.shouldMatch, got)
			}
		})
	}
}

func TestIsSubscriptionNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		res      *helix.RemoveEventSubSubscriptionParamsResponse
		expected bool
	}{
		{
			name:     "matches not found error text",
			err:      errors.New("not found"),
			expected: true,
		},
		{
			name:     "matches 404 response",
			res:      &helix.RemoveEventSubSubscriptionParamsResponse{ResponseCommon: helix.ResponseCommon{StatusCode: 404}},
			expected: true,
		},
		{
			name:     "matches not found response message",
			res:      &helix.RemoveEventSubSubscriptionParamsResponse{ResponseCommon: helix.ResponseCommon{ErrorMessage: "not found"}},
			expected: true,
		},
		{
			name:     "ignores other errors",
			err:      errors.New("boom"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSubscriptionNotFound(tt.err, tt.res)
			if got != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
