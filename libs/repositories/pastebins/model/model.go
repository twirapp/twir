package model

import (
	"net"
	"time"
)

type Pastebin struct {
	ID          string
	CreatedAt   time.Time
	Content     string
	ExpireAt    *time.Time
	OwnerUserID *string
	UserIP      *net.IP
	UserAgent   *string
}

var Nil = Pastebin{}
