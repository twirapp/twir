package stream

import "github.com/nicklaw5/helix/v2"

type HelixStream struct {
	helix.Stream

	Messages int `json:"parsedMessages"`
}
