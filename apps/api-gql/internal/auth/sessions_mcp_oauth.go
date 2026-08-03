package auth

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"
)

const mcpOAuthAttemptsKey = "mcpOAuthAttempts"

var (
	ErrMCPOAuthAttemptNotFound = errors.New("mcp oauth attempt not found")
	ErrMCPOAuthAttemptExpired  = errors.New("mcp oauth attempt expired")
)

type MCPOAuthAttempt struct {
	ClientID        string
	RedirectURI     string
	ClientState     string
	CodeChallenge   string
	RequestedScopes []string
	Resource        string
	CSRFToken       string
	ExpiresAt       time.Time
}

func (s *Auth) SetMCPOAuthAttempt(ctx context.Context, attemptID string, attempt MCPOAuthAttempt) error {
	attempts := s.mcpOAuthAttempts(ctx)
	attempts[attemptID] = cloneMCPOAuthAttempt(attempt)
	if err := s.saveMCPOAuthAttempts(ctx, attempts); err != nil {
		return fmt.Errorf("cannot commit MCP OAuth attempt: %w", err)
	}

	return nil
}

func (s *Auth) GetMCPOAuthAttempt(ctx context.Context, attemptID string) (MCPOAuthAttempt, error) {
	attempts := s.mcpOAuthAttempts(ctx)
	attempt, ok := attempts[attemptID]
	if !ok {
		return MCPOAuthAttempt{}, fmt.Errorf("%w: %s", ErrMCPOAuthAttemptNotFound, attemptID)
	}

	if !attempt.ExpiresAt.After(s.now()) {
		delete(attempts, attemptID)
		if err := s.saveMCPOAuthAttempts(ctx, attempts); err != nil {
			return MCPOAuthAttempt{}, fmt.Errorf("cannot remove expired MCP OAuth attempt: %w", err)
		}

		return MCPOAuthAttempt{}, fmt.Errorf("%w: %s", ErrMCPOAuthAttemptExpired, attemptID)
	}

	return cloneMCPOAuthAttempt(attempt), nil
}

func (s *Auth) DeleteMCPOAuthAttempt(ctx context.Context, attemptID string) error {
	attempts := s.mcpOAuthAttempts(ctx)
	if _, ok := attempts[attemptID]; !ok {
		return fmt.Errorf("%w: %s", ErrMCPOAuthAttemptNotFound, attemptID)
	}

	delete(attempts, attemptID)
	if err := s.saveMCPOAuthAttempts(ctx, attempts); err != nil {
		return fmt.Errorf("cannot commit MCP OAuth attempt deletion: %w", err)
	}

	return nil
}

func (s *Auth) mcpOAuthAttempts(ctx context.Context) map[string]MCPOAuthAttempt {
	storedAttempts, _ := s.sessionManager.Get(ctx, mcpOAuthAttemptsKey).(map[string]MCPOAuthAttempt)
	attempts := make(map[string]MCPOAuthAttempt, len(storedAttempts)+1)
	for attemptID, attempt := range storedAttempts {
		attempts[attemptID] = attempt
	}

	return attempts
}

func (s *Auth) saveMCPOAuthAttempts(ctx context.Context, attempts map[string]MCPOAuthAttempt) error {
	if len(attempts) == 0 {
		s.sessionManager.Remove(ctx, mcpOAuthAttemptsKey)
	} else {
		s.sessionManager.Put(ctx, mcpOAuthAttemptsKey, attempts)
	}

	if _, _, err := s.sessionManager.Commit(ctx); err != nil {
		return fmt.Errorf("commit session: %w", err)
	}

	return nil
}

func cloneMCPOAuthAttempt(attempt MCPOAuthAttempt) MCPOAuthAttempt {
	attempt.RequestedScopes = slices.Clone(attempt.RequestedScopes)
	return attempt
}
