package shortlinkscustomdomains

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	shortlinkscustomdomain "github.com/twirapp/twir/libs/entities/short_links_custom_domain"
	"github.com/twirapp/twir/libs/repositories/plans"
	shortlinkscustomdomainsrepo "github.com/twirapp/twir/libs/repositories/short_links_custom_domains"
	"github.com/twirapp/twir/libs/repositories/short_links_custom_domains/model"
	"go.uber.org/fx"
)

var reservedDomains = map[string]struct{}{
	"twir.app":                  {},
	"cf.twir.app":               {},
	"music-recognizer.twir.app": {},
	"services-bots.twir.app":    {},
}

type Opts struct {
	fx.In

	Repository      shortlinkscustomdomainsrepo.Repository
	PlansRepository plans.Repository
}

func New(opts Opts) *Service {
	return &Service{
		repository:      opts.Repository,
		plansRepository: opts.PlansRepository,
	}
}

type Service struct {
	repository      shortlinkscustomdomainsrepo.Repository
	plansRepository plans.Repository
}

type CreateInput struct {
	UserID string
	Domain string
}

func (s *Service) Create(ctx context.Context, input CreateInput) (shortlinkscustomdomain.Entity, error) {
	if input.UserID == "" {
		return shortlinkscustomdomain.Nil, errors.New("user ID cannot be empty")
	}

	normalizedDomain, err := normalizeDomain(input.Domain)
	if err != nil {
		return shortlinkscustomdomain.Nil, err
	}

	if err := validateDomain(input.UserID, normalizedDomain); err != nil {
		return shortlinkscustomdomain.Nil, err
	}

	plan, err := s.plansRepository.GetByChannelID(ctx, input.UserID)
	if err != nil {
		return shortlinkscustomdomain.Nil, fmt.Errorf("failed to get plan: %w", err)
	}
	if plan.IsNil() {
		return shortlinkscustomdomain.Nil, fmt.Errorf("plan not found for user")
	}
	if plan.LinksShortenerCustomDomains <= 0 {
		return shortlinkscustomdomain.Nil, fmt.Errorf("custom domains are not available on your plan")
	}

	// Generate verification token
	token := generateVerificationToken()

	existingCount, err := s.repository.CountByUserID(ctx, input.UserID)
	if err != nil {
		return shortlinkscustomdomain.Nil, fmt.Errorf("failed to count custom domains: %w", err)
	}
	if existingCount >= plan.LinksShortenerCustomDomains {
		return shortlinkscustomdomain.Nil, shortlinkscustomdomainsrepo.ErrUserAlreadyHasDomain
	}

	created, err := s.repository.Create(
		ctx, shortlinkscustomdomainsrepo.CreateInput{
			UserID:            input.UserID,
			Domain:            normalizedDomain,
			VerificationToken: token,
		},
	)
	if err != nil {
		return shortlinkscustomdomain.Nil, err
	}

	return mapToEntity(created), nil
}

func (s *Service) VerifyDomain(ctx context.Context, userID string) error {
	customDomain, err := s.repository.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if customDomain.IsNil() {
		return errors.New("custom domain not found")
	}

	if customDomain.Verified {
		return nil
	}

	expectedTarget := strings.ToLower(mapToEntity(customDomain).GetVerificationTarget())
	cname, err := net.LookupCNAME(customDomain.Domain)
	if err != nil {
		return errors.New("DNS lookup failed: " + err.Error())
	}

	actualTarget := strings.ToLower(strings.TrimSuffix(cname, "."))
	expectedTarget = strings.ToLower(strings.TrimSuffix(expectedTarget, "."))

	if actualTarget != expectedTarget {
		return errors.New("CNAME record does not match expected target")
	}

	return s.repository.VerifyDomain(ctx, customDomain.ID)
}

func (s *Service) GetByUserID(ctx context.Context, userID string) (shortlinkscustomdomain.Entity, error) {
	customDomain, err := s.repository.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, shortlinkscustomdomainsrepo.ErrNotFound) {
			return shortlinkscustomdomain.Nil, nil
		}
		return shortlinkscustomdomain.Nil, err
	}

	return mapToEntity(customDomain), nil
}

func (s *Service) GetByDomain(ctx context.Context, domain string) (shortlinkscustomdomain.Entity, error) {
	customDomain, err := s.repository.GetByDomain(ctx, domain)
	if err != nil {
		if errors.Is(err, shortlinkscustomdomainsrepo.ErrNotFound) {
			return shortlinkscustomdomain.Nil, nil
		}
		return shortlinkscustomdomain.Nil, err
	}

	return mapToEntity(customDomain), nil
}

func (s *Service) IsDomainAllowed(ctx context.Context, domain string) (bool, error) {
	normalizedDomain, err := normalizeDomain(domain)
	if err != nil {
		return false, nil
	}

	if _, ok := reservedDomains[normalizedDomain]; ok {
		return false, nil
	}

	customDomain, err := s.repository.GetByDomain(ctx, normalizedDomain)
	if err != nil {
		if errors.Is(err, shortlinkscustomdomainsrepo.ErrNotFound) {
			return false, nil
		}

		return false, err
	}

	return customDomain.Verified, nil
}

func (s *Service) Delete(ctx context.Context, userID string) error {
	customDomain, err := s.repository.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if customDomain.IsNil() {
		return errors.New("custom domain not found")
	}

	return s.repository.Delete(ctx, customDomain.ID)
}

func generateVerificationToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func mapToEntity(m model.CustomDomain) shortlinkscustomdomain.Entity {
	return shortlinkscustomdomain.Entity{
		ID:                m.ID,
		UserID:            m.UserID,
		Domain:            m.Domain,
		Verified:          m.Verified,
		VerificationToken: m.VerificationToken,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func validateDomain(userID, domain string) error {
	entity := shortlinkscustomdomain.Entity{
		UserID: userID,
		Domain: domain,
	}
	if err := entity.Validate(); err != nil {
		return err
	}

	if _, ok := reservedDomains[domain]; ok {
		return errors.New("domain is reserved")
	}

	return nil
}

func normalizeDomain(domain string) (string, error) {
	trimmed := strings.ToLower(strings.TrimSpace(domain))
	if trimmed == "" {
		return "", errors.New("domain cannot be empty")
	}

	if strings.Contains(trimmed, "://") {
		parsed, err := url.Parse(trimmed)
		if err != nil {
			return "", errors.New("domain has invalid format")
		}
		if parsed.Host == "" {
			return "", errors.New("domain has invalid format")
		}
		trimmed = parsed.Host
	}

	if strings.Contains(trimmed, "/") {
		trimmed = strings.Split(trimmed, "/")[0]
	}

	if strings.Contains(trimmed, ":") {
		if host, _, err := net.SplitHostPort(trimmed); err == nil {
			trimmed = host
		}
	}

	trimmed = strings.TrimSuffix(trimmed, ".")

	if trimmed == "" {
		return "", errors.New("domain cannot be empty")
	}

	return trimmed, nil
}
