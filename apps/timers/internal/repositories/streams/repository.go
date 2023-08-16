package streams

import "errors"

var NotFound = errors.New("stream not found")

type Repository interface {
	// GetByChannelId returns NotFound error if stream not found
	GetByChannelId(id string) (Stream, error)
}

type Stream struct {
	ID        string
	UserLogin string
	UserID    string
}
