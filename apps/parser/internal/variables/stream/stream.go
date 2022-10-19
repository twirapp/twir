package stream

import "github.com/satont/go-helix/v2"

type HelixStream struct {
	helix.Stream

	Messages int `json:"parsedMessages"`
}
