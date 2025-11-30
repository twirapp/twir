package model

import (
	"net/netip"
	"time"
)

type ShortenedUrl struct {
	ShortID         string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	URL             string
	CreatedByUserId *string
	Views           int
	UserAgent       *string
	UserIp          *netip.Addr
}

var Nil = ShortenedUrl{}
