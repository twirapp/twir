package mcp_oauth

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"regexp"
	"slices"
	"time"

	"github.com/google/uuid"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
	repository "github.com/twirapp/twir/libs/repositories/mcp_oauth"
)

const credentialBytes = 32

var verifierPattern = regexp.MustCompile(`^[A-Za-z0-9\-._~]{43,128}$`)

type CreateAuthorizationCodeInput struct {
	Authorize                  AuthorizeInput
	ChannelID, ApprovingUserID uuid.UUID
}
type IssuedAuthorizationCode struct {
	Code      string
	ExpiresAt time.Time
}
type ExchangeAuthorizationCodeInput struct{ ClientID, Code, RedirectURI, CodeVerifier, Resource string }
type RefreshInput struct{ ClientID, RefreshToken, Scope, Resource string }
type RevokeInput struct{ ClientID, Token string }
type TokenSet struct {
	AccessToken, RefreshToken, TokenType, Resource string
	Scopes                                         []entity.Scope
	AccessExpiresAt, RefreshExpiresAt              time.Time
}
type AuthorizedGrant struct {
	Client                     entity.Client
	ChannelID, ApprovingUserID uuid.UUID
	Channel                    channelentity.Channel
	Scopes                     []entity.Scope
	Resource                   string
}

func (s *Service) CreateAuthorizationCode(ctx context.Context, input CreateAuthorizationCodeInput) (IssuedAuthorizationCode, error) {
	request, err := s.ValidateAuthorizeInput(ctx, input.Authorize)
	if err != nil {
		return IssuedAuthorizationCode{}, err
	}
	if _, err := s.authorize(ctx, input.Authorize.ClientID, input.ApprovingUserID, input.ChannelID, entity.CredentialHash{}); err != nil {
		return IssuedAuthorizationCode{}, err
	}
	code, hash, err := s.credential()
	if err != nil {
		return IssuedAuthorizationCode{}, err
	}
	expires := s.clock.Now().Add(5 * time.Minute)
	_, err = s.repository.CreateAuthorizationCode(ctx, repository.CreateAuthorizationCodeInput{CodeHash: hash, ClientID: request.Client.ClientID, ChannelID: input.ChannelID, UserID: input.ApprovingUserID, RedirectURI: request.RedirectURI, PKCEChallenge: request.CodeChallenge, Scopes: request.Scopes, Resource: request.Resource, ExpiresAt: expires})
	if err != nil {
		return IssuedAuthorizationCode{}, err
	}
	return IssuedAuthorizationCode{Code: code, ExpiresAt: expires}, nil
}
func (s *Service) ExchangeAuthorizationCode(ctx context.Context, input ExchangeAuthorizationCodeInput) (TokenSet, error) {
	code, err := s.repository.ConsumeAuthorizationCode(ctx, credentialHash(input.Code))
	if err != nil {
		return TokenSet{}, oauthError(ErrorInvalidGrant, "invalid authorization code")
	}
	if code.ClientID != input.ClientID || code.RedirectURI != input.RedirectURI || code.Resource != s.resource || input.Resource != code.Resource || !code.ExpiresAt.After(s.clock.Now()) || !validVerifier(input.CodeVerifier) || subtle.ConstantTimeCompare([]byte(s256Challenge(input.CodeVerifier)), []byte(code.PKCEChallenge)) != 1 {
		return TokenSet{}, oauthError(ErrorInvalidGrant, "invalid authorization code")
	}
	return s.createToken(ctx, code.ClientID, code.ChannelID, code.UserID, code.Scopes, code.Resource, false, entity.CredentialHash{})
}
func (s *Service) Refresh(ctx context.Context, input RefreshInput) (TokenSet, error) {
	hash := credentialHash(input.RefreshToken)
	current, err := s.repository.GetTokenByRefreshTokenHash(ctx, hash)
	if err != nil || current.ClientID != input.ClientID || input.Resource != s.resource || current.Resource != input.Resource || !current.RefreshExpiresAt.After(s.clock.Now()) {
		return TokenSet{}, oauthError(ErrorInvalidGrant, "invalid refresh token")
	}
	requested := current.Scopes
	if input.Scope != "" {
		requested, err = parseScopes(input.Scope)
		if err != nil || !scopeSubset(requested, current.Scopes) {
			return TokenSet{}, oauthError(ErrorInvalidScope, "requested scope is not permitted")
		}
	}
	if _, err := s.authorize(ctx, current.ClientID, current.UserID, current.ChannelID, hash); err != nil {
		return TokenSet{}, err
	}
	return s.createToken(ctx, current.ClientID, current.ChannelID, current.UserID, requested, current.Resource, true, hash)
}
func (s *Service) Revoke(ctx context.Context, input RevokeInput) error {
	return s.repository.RevokeToken(ctx, input.ClientID, credentialHash(input.Token))
}
func (s *Service) VerifyAccessToken(ctx context.Context, raw string) (AuthorizedGrant, error) {
	hash := credentialHash(raw)
	token, err := s.repository.GetTokenByAccessTokenHash(ctx, hash)
	if err != nil || token.RevokedAt != nil || !token.AccessExpiresAt.After(s.clock.Now()) || token.Resource != s.resource {
		return AuthorizedGrant{}, oauthError(ErrorInvalidToken, "invalid access token")
	}
	channel, err := s.authorize(ctx, token.ClientID, token.UserID, token.ChannelID, hash)
	if err != nil {
		return AuthorizedGrant{}, err
	}
	client, err := s.GetClient(ctx, token.ClientID)
	if err != nil {
		return AuthorizedGrant{}, oauthError(ErrorInvalidToken, "invalid access token")
	}
	return AuthorizedGrant{Client: client, ChannelID: token.ChannelID, ApprovingUserID: token.UserID, Channel: channel, Scopes: slices.Clone(token.Scopes), Resource: s.resource}, nil
}
func (s *Service) createToken(ctx context.Context, clientID string, channelID, userID uuid.UUID, scopes []entity.Scope, resource string, rotate bool, refreshHash entity.CredentialHash) (TokenSet, error) {
	access, accessHash, err := s.credential()
	if err != nil {
		return TokenSet{}, err
	}
	refresh, refreshTokenHash, err := s.credential()
	if err != nil {
		return TokenSet{}, err
	}
	now := s.clock.Now()
	input := repository.CreateTokenInput{ClientID: clientID, ChannelID: channelID, UserID: userID, AccessTokenHash: accessHash, RefreshTokenHash: refreshTokenHash, Scopes: scopes, Resource: resource, AccessExpiresAt: now.Add(time.Hour), RefreshExpiresAt: now.Add(30 * 24 * time.Hour)}
	if rotate {
		_, err = s.repository.RotateRefreshToken(ctx, repository.RotateRefreshTokenInput{PresentedRefreshTokenHash: refreshHash, NextAccessTokenHash: accessHash, NextRefreshTokenHash: refreshTokenHash, Scopes: input.Scopes, AccessExpiresAt: input.AccessExpiresAt, RefreshExpiresAt: input.RefreshExpiresAt})
	} else {
		_, err = s.repository.CreateToken(ctx, input)
	}
	if err != nil {
		if errors.Is(err, repository.ErrRefreshTokenReuse) {
			return TokenSet{}, oauthError(ErrorInvalidGrant, "invalid refresh token")
		}
		return TokenSet{}, fmt.Errorf("issue MCP OAuth token: %w", err)
	}
	return TokenSet{AccessToken: access, RefreshToken: refresh, TokenType: "Bearer", Resource: resource, Scopes: slices.Clone(scopes), AccessExpiresAt: input.AccessExpiresAt, RefreshExpiresAt: input.RefreshExpiresAt}, nil
}
func (s *Service) authorize(ctx context.Context, clientID string, userID, channelID uuid.UUID, hash entity.CredentialHash) (channelentity.Channel, error) {
	user, err := s.users.GetByID(ctx, userID.String())
	if err != nil || user.ID == "" || user.IsBanned {
		return channelentity.Nil, s.reject(ctx, clientID, hash)
	}
	channel, err := s.channels.GetChannelByID(ctx, channelID)
	if err != nil || channel.IsNil() {
		return channelentity.Nil, s.reject(ctx, clientID, hash)
	}
	allowed, err := s.access.CanAccess(ctx, DashboardSubject{ID: user.ID, IsBotAdmin: user.IsBotAdmin}, channelID, manageBotSettings)
	if err != nil {
		return channelentity.Nil, err
	}
	if !allowed {
		return channelentity.Nil, s.reject(ctx, clientID, hash)
	}
	return channel, nil
}
func (s *Service) reject(ctx context.Context, clientID string, hash entity.CredentialHash) error {
	if hash == (entity.CredentialHash{}) {
		return oauthError(ErrorAccessDenied, "permission is required")
	}
	if hash != (entity.CredentialHash{}) {
		if err := s.repository.RevokeToken(ctx, clientID, hash); err != nil {
			return err
		}
	}
	return oauthError(ErrorInvalidToken, "authorization is no longer valid")
}
func (s *Service) credential() (string, entity.CredentialHash, error) {
	raw, err := s.randomValue()
	if err != nil {
		return "", entity.CredentialHash{}, err
	}
	return raw, credentialHash(raw), nil
}
func (s *Service) randomValue() (string, error) {
	bytes := make([]byte, credentialBytes)
	if _, err := io.ReadFull(s.random, bytes); err != nil {
		return "", fmt.Errorf("generate MCP OAuth credential: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
func credentialHash(raw string) entity.CredentialHash {
	return entity.CredentialHash(sha256.Sum256([]byte(raw)))
}
func validVerifier(value string) bool { return verifierPattern.MatchString(value) }
func s256Challenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
