package stream

import "github.com/nicklaw5/helix"

type HelixStream struct {
	helix.Stream

	Messages int `json:"parsedMessages"`
}
