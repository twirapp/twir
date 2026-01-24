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
	Domain          *string

	isNil bool
}

func (s ShortenedUrl) IsNil() bool {
	return s.isNil
}

var Nil = ShortenedUrl{
	isNil: true,
}
