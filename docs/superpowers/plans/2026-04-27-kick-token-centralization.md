# Kick Token Centralization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Centralize Kick bot token refresh in `apps/tokens` and make Kick chat length enforcement/splitting match Kick API byte limit behavior.

**Architecture:** Reuse existing tokens bus queue so bots service asks tokens service for Kick bot access tokens instead of reading and refreshing encrypted credentials directly. Move Kick-specific token refresh/persist logic into tokens service, then make Kick chat validation and splitting byte-based with UTF-8-safe boundaries.

**Tech Stack:** Go, fx, NATS bus via `libs/bus-core`, pgx repositories, local `gokick` client, Go testing package.

---

### Task 1: Add bus contract for bot platform

**Files:**
- Modify: `libs/bus-core/tokens/tokens.go`
- Test: none

- [ ] **Step 1: Add `Platform` to bot token request**

```go
type GetBotTokenRequest struct {
	BotId    string                    `json:"botId"`
	Platform platformentity.Platform   `json:"platform"`
}
```

- [ ] **Step 2: Run focused build to verify compile surface**

Run: `go test ./libs/bus-core/...`
Expected: PASS

### Task 2: Make `gokick` validate by bytes

**Files:**
- Modify: `gokick/chat.go`
- Test: `gokick/chat_test.go`

- [ ] **Step 1: Write failing test for multibyte overflow**

```go
func TestValidateChatMessageContent_BytesLimit(t *testing.T) {
	require.NoError(t, gokick.ValidateChatMessageContent(strings.Repeat("a", 500)))
	require.NoError(t, gokick.ValidateChatMessageContent(strings.Repeat("ы", 250)))
	require.ErrorIs(t, gokick.ValidateChatMessageContent(strings.Repeat("ы", 251)), gokick.ErrChatMessageContentTooLong)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./gokick -run TestValidateChatMessageContent_BytesLimit`
Expected: FAIL because current validator counts runes, not bytes

- [ ] **Step 3: Implement byte-based validation**

```go
func ValidateChatMessageContent(content string) error {
	if len(content) > ChatMessageContentMaxRunes {
		return ErrChatMessageContentTooLong
	}
	return nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./gokick -run TestValidateChatMessageContent_BytesLimit`
Expected: PASS

### Task 3: Add Kick splitter tests in bots service

**Files:**
- Modify: `apps/bots/internal/kick/chat_client_test.go`

- [ ] **Step 1: Write failing tests for byte-based split behavior**

```go
func TestSplitMessage_UsesByteLimit(t *testing.T) {
	parts := splitMessage(strings.Repeat("ы", 251))
	require.Len(t, parts, 2)
	require.Len(t, []byte(parts[0]), 500)
	require.Len(t, []byte(parts[1]), 2)
}

func TestSplitMessage_PreservesUTF8Boundaries(t *testing.T) {
	parts := splitMessage(strings.Repeat("界", 167))
	for _, part := range parts {
		require.True(t, utf8.ValidString(part))
		require.LessOrEqual(t, len([]byte(part)), 500)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./apps/bots/internal/kick -run TestSplitMessage`
Expected: FAIL because current splitter chunks by runes

- [ ] **Step 3: Implement byte-based splitting**

```go
func splitMessage(text string) []string {
	normalizedText := strings.ReplaceAll(text, "\n", " ")
	if normalizedText == "" {
		return nil
	}

	if err := gokick.ValidateChatMessageContent(normalizedText); err == nil {
		return []string{normalizedText}
	}

	parts := []string{}
	for len(normalizedText) > 0 {
		end := 0
		for i := range normalizedText {
			if i == 0 {
				continue
			}
			if len(normalizedText[:i]) > gokick.ChatMessageContentMaxRunes {
				break
			}
			end = i
		}

		if end == 0 || len(normalizedText) <= gokick.ChatMessageContentMaxRunes {
			end = len(normalizedText)
		}

		parts = append(parts, normalizedText[:end])
		normalizedText = normalizedText[end:]
	}

	return parts
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./apps/bots/internal/kick -run TestSplitMessage`
Expected: PASS

### Task 4: Move Kick bot token refresh into tokens service

**Files:**
- Modify: `apps/tokens/internal/bus_listener/bus_listener.go`
- Test: `apps/tokens/internal/bus_listener/bus_listener_test.go`
- Read: `libs/repositories/kick_bots/repository.go`

- [ ] **Step 1: Add failing tests for platform default and Kick bot path**

```go
func TestRequestBotToken_DefaultsToTwitchWhenPlatformEmpty(t *testing.T) {}
func TestRequestBotToken_KickRefreshesDefaultBot(t *testing.T) {}
```

Test expectations:
- Empty `Platform` uses existing Twitch bot token path.
- `PlatformKick` loads default Kick bot, decrypts refresh token, refreshes through Kick client, persists updated encrypted token, returns decrypted access token.

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./apps/tokens/internal/bus_listener -run TestRequestBotToken`
Expected: FAIL because Kick branch does not exist yet

- [ ] **Step 3: Inject Kick bot repository and implement Kick branch**

Implementation outline:

```go
type Opts struct {
	// existing fields
	KickBotsRepository kick_bots.Repository
}

func (c *tokensImpl) RequestBotToken(ctx context.Context, data tokens.GetBotTokenRequest) (tokens.TokenResponse, error) {
	platform := data.Platform
	if platform == "" {
		platform = platformentity.PlatformTwitch
	}

	if platform == platformentity.PlatformKick {
		return c.requestKickBotToken(ctx)
	}

	// existing Twitch flow
}
```

Add helper:

```go
func (c *tokensImpl) requestKickBotToken(ctx context.Context) (tokens.TokenResponse, error) {
	// get default bot
	// decrypt refresh token
	// refresh if expired
	// persist encrypted tokens
	// decrypt access token if needed
	// return TokenResponse
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./apps/tokens/internal/bus_listener -run TestRequestBotToken`
Expected: PASS

### Task 5: Make bots Kick client consume tokens service

**Files:**
- Modify: `apps/bots/internal/kick/chat_client.go`
- Modify: `apps/bots/app/app.go`
- Test: `apps/bots/internal/kick/chat_client_test.go`

- [ ] **Step 1: Write failing test proving bus token usage**

```go
func TestSendMessage_RequestsKickTokenFromBus(t *testing.T) {
	// fake bus returns token response
	// fake HTTP transport captures Authorization header
	// expect no repo decrypt/refresh path required
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./apps/bots/internal/kick -run TestSendMessage_RequestsKickTokenFromBus`
Expected: FAIL because client currently depends on `kick_bots.Repository` and local decrypt logic

- [ ] **Step 3: Replace local token logic with bus request**

Implementation outline:

```go
type ChatClient struct {
	twirBus    *buscore.Bus
	config     cfg.Config
	httpClient *http.Client
	logger     *slog.Logger
}

tokenResp, err := c.twirBus.Tokens.RequestBotToken.Request(ctx, tokens.GetBotTokenRequest{
	Platform: platformentity.PlatformKick,
})
```

Remove:
- `crypto` import
- `kick_bots.Repository` dependency
- `refreshBotToken`
- unauthorized retry recursion

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./apps/bots/internal/kick -run 'TestSendMessage_RequestsKickTokenFromBus|TestSplitMessage'`
Expected: PASS

### Task 6: Full verification

**Files:**
- No code changes

- [ ] **Step 1: Run package tests**

Run: `go test ./gokick ./apps/tokens/internal/bus_listener ./apps/bots/internal/kick`
Expected: PASS

- [ ] **Step 2: Run app builds/tests covering compile integration**

Run: `bun cli build tokens && bun cli build bots`
Expected: PASS

- [ ] **Step 3: Review git diff for scope**

Run: `git diff -- apps/tokens/internal/bus_listener/bus_listener.go apps/bots/internal/kick/chat_client.go gokick/chat.go libs/bus-core/tokens/tokens.go`
Expected: only token-centralization and byte-limit changes
