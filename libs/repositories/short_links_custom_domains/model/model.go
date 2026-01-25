package model

import "time"

type CustomDomain struct {
	ID                string
	UserID            string
	Domain            string
	Verified          bool
	VerificationToken string
	CreatedAt         time.Time
	UpdatedAt         time.Time

	isNil bool
}

func (c CustomDomain) IsNil() bool {
	return c.isNil
}

var Nil = CustomDomain{
	isNil: true,
}
