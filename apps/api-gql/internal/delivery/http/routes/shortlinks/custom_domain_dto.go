package shortlinks

import (
	"time"

	shortlinkscustomdomain "github.com/twirapp/twir/libs/entities/short_links_custom_domain"
)

type customDomainOutputDto struct {
	ID                 string    `json:"id"`
	Domain             string    `json:"domain"`
	Verified           bool      `json:"verified"`
	VerificationToken  string    `json:"verification_token"`
	VerificationTarget string    `json:"verification_target"`
	CreatedAt          time.Time `json:"created_at"`
}

func mapCustomDomainOutput(entity shortlinkscustomdomain.Entity, baseUrl string) customDomainOutputDto {
	return customDomainOutputDto{
		ID:                 entity.ID,
		Domain:             entity.Domain,
		Verified:           entity.Verified,
		VerificationToken:  entity.VerificationToken,
		VerificationTarget: entity.GetVerificationTarget(baseUrl),
		CreatedAt:          entity.CreatedAt,
	}
}
