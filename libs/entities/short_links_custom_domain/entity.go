package shortlinkscustomdomain

import (
	"errors"
	"regexp"
	"time"
)

var domainRegex = regexp.MustCompile(
	`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`,
)

type Entity struct {
	ID                string
	UserID            string
	Domain            string
	Verified          bool
	VerificationToken string
	CreatedAt         time.Time
	UpdatedAt         time.Time

	isNil bool
}

func (c Entity) IsNil() bool {
	return c.isNil
}

var Nil = Entity{
	isNil: true,
}

func (c Entity) GetVerificationTarget() string {
	return c.VerificationToken + ".shortener.twir.app"
}

func (c Entity) Validate() error {
	if c.Domain == "" {
		return errors.New("domain cannot be empty")
	}
	if c.UserID == "" {
		return errors.New("user ID cannot be empty")
	}
	if !domainRegex.MatchString(c.Domain) {
		return errors.New("domain has invalid format")
	}

	return nil
}
