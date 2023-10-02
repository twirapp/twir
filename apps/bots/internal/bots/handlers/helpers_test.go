package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlers_ParseEmotes(t *testing.T) {
	table := []struct {
		name        string
		messageText string
		raw         string
		expected    []MessageEmote
	}{
		{
			name:        "empty",
			messageText: "",
			raw:         "",
			expected:    nil,
		},
		{
			name:        "one emote",
			messageText: "йцу Kappa Kappa",
			raw:         "25:4-8,10-14",
			expected: []MessageEmote{
				{
					ID:    "25",
					Name:  "Kappa",
					Count: 2,
					Positions: []EmotePosition{
						{
							Start: 4,
							End:   8,
						},
						{
							Start: 10,
							End:   14,
						},
					},
				},
			},
		},
	}

	for _, tt := range table {
		t.Run(
			tt.name,
			func(t *testing.T) {
				h := &Handlers{}
				actual := h.ParseEmotes(tt.messageText, tt.raw)
				assert.Equal(t, tt.expected, actual)
			},
		)
	}
}
