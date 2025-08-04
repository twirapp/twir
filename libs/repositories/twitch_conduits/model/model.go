package model

import (
	"time"
)

type Conduit struct {
	ID         string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ShardCount int8
}

var Nil = Conduit{}
