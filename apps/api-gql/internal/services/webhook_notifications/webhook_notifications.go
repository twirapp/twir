package webhook_notifications

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	gojson "github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	"github.com/twirapp/twir/libs/audit"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/entities/webhook_notifications"
	channelsmoduleswebhooks "github.com/twirapp/twir/libs/repositories/channels_modules_webhooks"
	"go.uber.org/fx"
)

var ErrInvalidPayload = errors.New("invalid webhook payload")

const webhookNotificationsPageSize = 500

type Opts struct {
	fx.In

	Repository    channelsmoduleswebhooks.Repository
	AuditRecorder audit.Recorder
	Bus           *buscore.Bus
	Logger        *slog.Logger
	ShortenedUrls *shortenedurls.Service
	Config        cfg.Config
}

type Service struct {
	repo          channelsmoduleswebhooks.Repository
	auditRecorder audit.Recorder
	bus           *buscore.Bus
	logger        *slog.Logger
	shortenedUrls *shortenedurls.Service
	config        cfg.Config
}

func New(opts Opts) *Service {
	return &Service{
		repo:          opts.Repository,
		auditRecorder: opts.AuditRecorder,
		bus:           opts.Bus,
		logger:        opts.Logger,
		shortenedUrls: opts.ShortenedUrls,
		config:        opts.Config,
	}
}

func (s *Service) GetByChannelID(
	ctx context.Context,
	channelID string,
) (webhook_notifications.Settings, error) {
	settings, err := s.repo.GetByChannelID(ctx, channelID)
	if err != nil {
		return webhook_notifications.Nil, fmt.Errorf("failed to get webhook module settings: %w", err)
	}

	return settings, nil
}

type CreateInput struct {
	ChannelID string
	ActorID   string

	Enabled                   bool
	GithubIssuesEnabled       bool
	GithubPullRequestsEnabled bool
	GithubCommitsEnabled      bool
}

func (s *Service) Create(ctx context.Context, input CreateInput) (
	webhook_notifications.Settings,
	error,
) {
	settings, err := s.repo.Create(
		ctx,
		channelsmoduleswebhooks.CreateInput{
			ChannelID:                 input.ChannelID,
			Enabled:                   input.Enabled,
			GithubIssuesEnabled:       input.GithubIssuesEnabled,
			GithubPullRequestsEnabled: input.GithubPullRequestsEnabled,
			GithubCommitsEnabled:      input.GithubCommitsEnabled,
		},
	)
	if err != nil {
		return webhook_notifications.Nil, fmt.Errorf("failed to create webhook module settings: %w", err)
	}

	_ = s.auditRecorder.RecordCreateOperation(
		ctx,
		audit.CreateOperation{
			Metadata: audit.OperationMetadata{
				System:    "channels_modules_webhooks",
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(settings.ID),
			},
			NewValue: settings,
		},
	)

	return settings, nil
}

type UpdateInput struct {
	ChannelID string
	ActorID   string

	Enabled                   *bool
	GithubIssuesEnabled       *bool
	GithubPullRequestsEnabled *bool
	GithubCommitsEnabled      *bool
}

func (s *Service) Update(
	ctx context.Context,
	id uuid.UUID,
	input UpdateInput,
) (webhook_notifications.Settings, error) {
	existing, err := s.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		if errors.Is(err, channelsmoduleswebhooks.ErrSettingsNotFound) {
			return webhook_notifications.Nil, errors.New("webhook module settings not found")
		}
		return webhook_notifications.Nil, err
	}

	if existing.ChannelID != input.ChannelID {
		return webhook_notifications.Nil, errors.New("webhook module settings do not belong to this channel")
	}

	settings, err := s.repo.Update(
		ctx,
		id,
		channelsmoduleswebhooks.UpdateInput{
			Enabled:                   input.Enabled,
			GithubIssuesEnabled:       input.GithubIssuesEnabled,
			GithubPullRequestsEnabled: input.GithubPullRequestsEnabled,
			GithubCommitsEnabled:      input.GithubCommitsEnabled,
		},
	)
	if err != nil {
		return webhook_notifications.Nil, fmt.Errorf("failed to update webhook module settings: %w", err)
	}

	_ = s.auditRecorder.RecordUpdateOperation(
		ctx,
		audit.UpdateOperation{
			Metadata: audit.OperationMetadata{
				System:    "channels_modules_webhooks",
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(settings.ID),
			},
			NewValue: settings,
			OldValue: existing,
		},
	)

	return settings, nil
}

type DeleteInput struct {
	ID        uuid.UUID
	ChannelID string
	ActorID   string
}

func (s *Service) Delete(ctx context.Context, input DeleteInput) error {
	oldSettings, err := s.repo.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return fmt.Errorf("get webhook module settings: %w", err)
	}

	if oldSettings.ChannelID != input.ChannelID {
		return errors.New("webhook module settings do not belong to this channel")
	}

	if err := s.repo.Delete(ctx, input.ID); err != nil {
		return fmt.Errorf("delete webhook module settings: %w", err)
	}

	_ = s.auditRecorder.RecordDeleteOperation(
		ctx,
		audit.DeleteOperation{
			Metadata: audit.OperationMetadata{
				System:    "channels_modules_webhooks",
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(oldSettings.ID),
			},
			OldValue: &oldSettings,
		},
	)

	return nil
}

func (s *Service) HandleGithubWebhook(
	ctx context.Context,
	event string,
	payload []byte,
) error {
	event = strings.ToLower(event)
	switch event {
	case "ping":
		return nil
	case "issues":
		return s.handleGithubIssues(ctx, payload)
	case "issue_comment":
		return s.handleGithubIssueComment(ctx, payload)
	case "pull_request":
		return s.handleGithubPullRequest(ctx, payload)
	case "push":
		return s.handleGithubPush(ctx, payload)
	default:
		return nil
	}
}

func (s *Service) handleGithubIssues(
	ctx context.Context,
	payload []byte,
) error {
	var data githubIssuesPayload
	if err := gojson.Unmarshal(payload, &data); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidPayload, err)
	}

	action := strings.ToLower(data.Action)
	if action == "" {
		return nil
	}

	verb := ""
	switch action {
	case "opened":
		verb = "opened"
	case "edited":
		verb = "updated"
	case "closed":
		verb = "closed"
		if data.Issue.StateReason != nil && *data.Issue.StateReason == "completed" {
			verb = "resolved"
		}
	default:
		return nil
	}

	user := githubActor(data.Sender, data.Issue.User)

	repo := data.Repository.FullName
	if repo == "" {
		repo = "repository"
	}

	issueURL := s.shortenUrl(ctx, data.Issue.HTMLURL)
	message := fmt.Sprintf(
		"[GitHub] Issue #%d %s by %s in %s: %s",
		data.Issue.Number,
		verb,
		user,
		repo,
		strings.TrimSpace(data.Issue.Title),
	)
	if issueURL != "" {
		message = fmt.Sprintf("%s - %s", message, issueURL)
	}

	return s.sendMessageToEnabledChannels(
		ctx,
		channelsmoduleswebhooks.GetEnabledChannelsInput{
			GithubIssuesEnabled: lo.ToPtr(true),
		},
		[]string{message},
	)
}

func (s *Service) handleGithubIssueComment(
	ctx context.Context,
	payload []byte,
) error {
	var data githubIssueCommentPayload
	if err := gojson.Unmarshal(payload, &data); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidPayload, err)
	}

	if strings.ToLower(data.Action) != "created" {
		return nil
	}

	if data.Issue.PullRequest != nil {
		return nil
	}

	repo := data.Repository.FullName
	if repo == "" {
		repo = "repository"
	}

	actor := githubActor(data.Sender, data.Comment.User)
	comment := firstLine(data.Comment.Body)
	commentURL := s.shortenUrl(ctx, data.Comment.HTMLURL)

	message := fmt.Sprintf(
		"[GitHub] Issue #%d comment by %s in %s",
		data.Issue.Number,
		actor,
		repo,
	)

	if comment != "" {
		message = fmt.Sprintf("%s: %s", message, comment)
	}

	if commentURL != "" {
		message = fmt.Sprintf("%s - %s", message, commentURL)
	}

	return s.sendMessageToEnabledChannels(
		ctx,
		channelsmoduleswebhooks.GetEnabledChannelsInput{
			GithubIssuesEnabled: lo.ToPtr(true),
		},
		[]string{message},
	)
}

func (s *Service) handleGithubPullRequest(
	ctx context.Context,
	payload []byte,
) error {
	var data githubPullRequestPayload
	if err := gojson.Unmarshal(payload, &data); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidPayload, err)
	}

	action := strings.ToLower(data.Action)
	if action == "" {
		return nil
	}

	verb := ""
	switch action {
	case "opened":
		verb = "opened"
	case "closed":
		if data.PullRequest.Merged {
			verb = "merged"
		} else {
			verb = "closed"
		}
	default:
		return nil
	}

	user := githubActor(data.Sender, data.PullRequest.User)

	repo := data.Repository.FullName
	if repo == "" {
		repo = "repository"
	}

	pullRequestURL := s.shortenUrl(ctx, data.PullRequest.HTMLURL)
	message := fmt.Sprintf(
		"[GitHub] PR #%d %s by %s in %s: %s",
		data.PullRequest.Number,
		verb,
		user,
		repo,
		strings.TrimSpace(data.PullRequest.Title),
	)
	if pullRequestURL != "" {
		message = fmt.Sprintf("%s - %s", message, pullRequestURL)
	}

	return s.sendMessageToEnabledChannels(
		ctx,
		channelsmoduleswebhooks.GetEnabledChannelsInput{
			GithubPullRequestsEnabled: lo.ToPtr(true),
		},
		[]string{message},
	)
}

func (s *Service) handleGithubPush(
	ctx context.Context,
	payload []byte,
) error {
	var data githubPushPayload
	if err := gojson.Unmarshal(payload, &data); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidPayload, err)
	}

	if len(data.Commits) == 0 {
		return nil
	}

	repo := data.Repository.FullName
	if repo == "" {
		repo = "repository"
	}

	branch := trimGithubRef(data.Ref)
	if branch == "" {
		branch = "unknown"
	}

	maxMessages := 5
	messages := make([]string, 0, maxMessages)
	for i, commit := range data.Commits {
		if i >= maxMessages {
			break
		}

		author := commit.Author.Username
		if author == "" {
			author = commit.Author.Name
		}
		if author == "" {
			author = "unknown"
		}

		commitURL := commit.URL
		if repo != "" && repo != "repository" {
			commitURL = fmt.Sprintf("https://github.com/%s/commit/%s", repo, commit.ID)
		}

		message := fmt.Sprintf(
			"[GitHub] %s@%s %s by %s: %s - %s",
			repo,
			branch,
			shortSHA(commit.ID),
			author,
			firstLine(commit.Message),
			commitURL,
		)

		messages = append(messages, message)
	}

	return s.sendMessageToEnabledChannels(
		ctx,
		channelsmoduleswebhooks.GetEnabledChannelsInput{
			GithubCommitsEnabled: lo.ToPtr(true),
		},
		messages,
	)
}

func (s *Service) sendMessage(ctx context.Context, channelID string, message string) error {
	if strings.TrimSpace(message) == "" {
		return nil
	}

	if err := s.bus.Bots.SendMessage.Publish(
		ctx,
		bots.SendMessageRequest{
			ChannelId:      channelID,
			Message:        message,
			SkipRateLimits: true,
		},
	); err != nil {
		s.logger.Error("failed to send webhook notification message", slog.String("channel_id", channelID), slog.String("message", message), slog.Any("error", err))
		return err
	}

	return nil
}

func (s *Service) sendMessageToChannels(
	ctx context.Context,
	channels []string,
	message string,
) error {
	if len(channels) == 0 {
		return nil
	}

	failed := 0
	for _, channelID := range channels {
		if err := s.sendMessage(ctx, channelID, message); err != nil {
			failed++
		}
	}

	if failed == len(channels) {
		return errors.New("failed to send message to all channels")
	}

	return nil
}

func (s *Service) sendMessageToEnabledChannels(
	ctx context.Context,
	input channelsmoduleswebhooks.GetEnabledChannelsInput,
	messages []string,
) error {
	if len(messages) == 0 {
		return nil
	}

	filtered := make([]string, 0, len(messages))
	for _, message := range messages {
		if strings.TrimSpace(message) == "" {
			continue
		}
		filtered = append(filtered, message)
	}

	if len(filtered) == 0 {
		return nil
	}

	return s.withEnabledChannelsPaged(
		ctx,
		input,
		func(channels []string) error {
			for _, message := range filtered {
				if err := s.sendMessageToChannels(ctx, channels, message); err != nil {
					return err
				}
			}
			return nil
		},
	)
}

func (s *Service) withEnabledChannelsPaged(
	ctx context.Context,
	input channelsmoduleswebhooks.GetEnabledChannelsInput,
	handle func([]string) error,
) error {
	perPage := input.PerPage
	if perPage <= 0 {
		perPage = webhookNotificationsPageSize
	}

	page := input.Page
	if page < 0 {
		page = 0
	}

	for {
		pageInput := input
		pageInput.Page = page
		pageInput.PerPage = perPage

		channels, err := s.repo.GetEnabledChannels(ctx, pageInput)
		if err != nil {
			return err
		}
		if len(channels) == 0 {
			return nil
		}

		if err := handle(channels); err != nil {
			return err
		}

		if len(channels) < perPage {
			return nil
		}

		page++
	}
}

type githubRepo struct {
	FullName string `json:"full_name"`
}

type githubUser struct {
	Login string `json:"login"`
}

type githubIssue struct {
	Number      int        `json:"number"`
	Title       string     `json:"title"`
	HTMLURL     string     `json:"html_url"`
	User        githubUser `json:"user"`
	StateReason *string    `json:"state_reason"`
	PullRequest *struct{}  `json:"pull_request"`
}

type githubIssuesPayload struct {
	Action     string      `json:"action"`
	Issue      githubIssue `json:"issue"`
	Repository githubRepo  `json:"repository"`
	Sender     githubUser  `json:"sender"`
}

type githubPullRequest struct {
	Number  int        `json:"number"`
	Title   string     `json:"title"`
	HTMLURL string     `json:"html_url"`
	User    githubUser `json:"user"`
	Merged  bool       `json:"merged"`
}

type githubPullRequestPayload struct {
	Action      string            `json:"action"`
	PullRequest githubPullRequest `json:"pull_request"`
	Repository  githubRepo        `json:"repository"`
	Sender      githubUser        `json:"sender"`
}

type githubIssueComment struct {
	HTMLURL string     `json:"html_url"`
	Body    string     `json:"body"`
	User    githubUser `json:"user"`
}

type githubIssueCommentPayload struct {
	Action     string             `json:"action"`
	Issue      githubIssue        `json:"issue"`
	Comment    githubIssueComment `json:"comment"`
	Repository githubRepo         `json:"repository"`
	Sender     githubUser         `json:"sender"`
}

type githubCommitAuthor struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

type githubCommit struct {
	ID      string             `json:"id"`
	Message string             `json:"message"`
	URL     string             `json:"url"`
	Author  githubCommitAuthor `json:"author"`
}

type githubPushPayload struct {
	Ref        string         `json:"ref"`
	Repository githubRepo     `json:"repository"`
	Commits    []githubCommit `json:"commits"`
}

func githubActor(sender githubUser, fallback githubUser) string {
	if sender.Login != "" {
		return sender.Login
	}
	if fallback.Login != "" {
		return fallback.Login
	}
	return "unknown"
}

func (s *Service) shortenUrl(ctx context.Context, rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" || s.shortenedUrls == nil {
		return rawURL
	}

	link, err := s.shortenedUrls.GetByUrl(ctx, nil, rawURL)
	if err != nil {
		s.logger.Warn("failed to get short link", slog.String("url", rawURL), slog.Any("error", err))
		return rawURL
	}

	if link.IsNil() {
		link, err = s.shortenedUrls.Create(
			ctx,
			shortenedurls.CreateInput{
				URL: rawURL,
			},
		)
		if err != nil {
			s.logger.Warn("failed to create short link", slog.String("url", rawURL), slog.Any("error", err))
			return rawURL
		}
	}

	shortURL, err := s.buildShortUrl(ctx, link.ShortID)
	if err != nil {
		s.logger.Warn(
			"failed to build short link",
			slog.String("short_id", link.ShortID),
			slog.Any("error", err),
		)
		return rawURL
	}

	return shortURL
}

func (s *Service) buildShortUrl(ctx context.Context, shortID string) (string, error) {
	baseUrl, _ := gincontext.GetBaseUrlFromContext(ctx, s.config.SiteBaseUrl)

	parsedBaseUrl, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	parsedBaseUrl.Path = "/s/" + shortID

	return parsedBaseUrl.String(), nil
}

func trimGithubRef(ref string) string {
	const prefix = "refs/heads/"
	if strings.HasPrefix(ref, prefix) {
		return strings.TrimPrefix(ref, prefix)
	}
	return ref
}

func shortSHA(sha string) string {
	if len(sha) <= 7 {
		return sha
	}
	return sha[:7]
}

func firstLine(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	if idx := strings.Index(text, "\n"); idx > -1 {
		return text[:idx]
	}
	return text
}
