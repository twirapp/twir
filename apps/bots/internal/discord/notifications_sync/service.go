package notifications_sync

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/twirapp/twir/apps/bots/internal/discord/discord_go"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/notifications"
	"github.com/twirapp/twir/libs/wsrouter"
	"go.uber.org/fx"
)

const notificationsSubscriptionKey = "api.newNotifications"

type Opts struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    cfg.Config
	Logger    *slog.Logger
	Discord   *discord_go.Discord
	Repo      notifications.Repository
	WsRouter  wsrouter.WsRouter
}

type eventKind uint8

const (
	eventBackfill eventKind = iota
	eventCreate
	eventUpdate
	eventDelete
)

type syncEvent struct {
	kind       eventKind
	message    *discord.Message
	messageIDs []discord.MessageID
}

type notificationEvent struct {
	ID           string    `json:"id"`
	UserID       *string   `json:"userId,omitempty"`
	Text         *string   `json:"text,omitempty"`
	EditorJSJSON *string   `json:"editorJsJson,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	Deleted      bool      `json:"deleted"`
}

type Service struct {
	discord      *discord_go.Discord
	repo         notifications.Repository
	wsRouter     wsrouter.WsRouter
	logger       *slog.Logger
	store        *attachmentStore
	channelID    discord.ChannelID
	historyLimit uint
	events       chan syncEvent

	lifecycleMu sync.RWMutex
	runContext  context.Context
	cancel      context.CancelFunc
	done        chan struct{}
}

func New(opts Opts) (*Service, error) {
	service := &Service{
		discord:      opts.Discord,
		repo:         opts.Repo,
		wsRouter:     opts.WsRouter,
		logger:       logger.WithComponent(opts.Logger, "discord-notifications-sync"),
		historyLimit: opts.Config.DiscordNotificationsHistoryLimit,
		events:       make(chan syncEvent, 128),
	}

	if opts.Config.DiscordNotificationsChannelID == "" {
		service.logger.Info("Discord notifications sync is disabled")
		return service, nil
	}
	if opts.Config.DiscordBotToken == "" {
		return nil, fmt.Errorf(
			"DISCORD_BOT_TOKEN is required when DISCORD_NOTIFICATIONS_CHANNEL_ID is set",
		)
	}

	channelSnowflake, err := discord.ParseSnowflake(opts.Config.DiscordNotificationsChannelID)
	if err != nil {
		return nil, fmt.Errorf("parse Discord notifications channel ID: %w", err)
	}
	service.channelID = discord.ChannelID(channelSnowflake)

	service.store, err = newAttachmentStore(opts.Config)
	if err != nil {
		return nil, err
	}
	if service.store == nil {
		service.logger.Warn(
			"S3 attachment mirroring is disabled; Discord CDN URLs may expire",
		)
	}

	opts.Discord.AddHandler(service.handleMessageCreate)
	opts.Discord.AddHandler(service.handleMessageUpdate)
	opts.Discord.AddHandler(service.handleMessageDelete)
	opts.Discord.AddHandler(service.handleMessageDeleteBulk)

	opts.Lifecycle.Append(fx.Hook{
		OnStart: service.start,
		OnStop:  service.stop,
	})

	return service, nil
}

func (s *Service) start(context.Context) error {
	runContext, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	s.lifecycleMu.Lock()
	s.runContext = runContext
	s.cancel = cancel
	s.done = done
	s.lifecycleMu.Unlock()

	go func() {
		defer close(done)
		s.run(runContext)
	}()
	s.events <- syncEvent{kind: eventBackfill}

	s.logger.Info(
		"Discord notifications sync started",
		slog.String("channel_id", s.channelID.String()),
		slog.Uint64("history_limit", uint64(s.historyLimit)),
	)
	return nil
}

func (s *Service) stop(ctx context.Context) error {
	s.lifecycleMu.RLock()
	cancel := s.cancel
	done := s.done
	s.lifecycleMu.RUnlock()
	if cancel == nil || done == nil {
		return nil
	}

	cancel()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Service) enqueue(event syncEvent) {
	s.lifecycleMu.RLock()
	runContext := s.runContext
	s.lifecycleMu.RUnlock()
	if runContext == nil {
		return
	}

	select {
	case s.events <- event:
	case <-runContext.Done():
	}
}

func (s *Service) handleMessageCreate(event *gateway.MessageCreateEvent) {
	if event.ChannelID != s.channelID {
		return
	}
	message := event.Message
	s.enqueue(syncEvent{kind: eventCreate, message: &message})
}

func (s *Service) handleMessageUpdate(event *gateway.MessageUpdateEvent) {
	if event.ChannelID != s.channelID {
		return
	}
	s.enqueue(syncEvent{kind: eventUpdate, messageIDs: []discord.MessageID{event.ID}})
}

func (s *Service) handleMessageDelete(event *gateway.MessageDeleteEvent) {
	if event.ChannelID != s.channelID {
		return
	}
	s.enqueue(syncEvent{kind: eventDelete, messageIDs: []discord.MessageID{event.ID}})
}

func (s *Service) handleMessageDeleteBulk(event *gateway.MessageDeleteBulkEvent) {
	if event.ChannelID != s.channelID {
		return
	}
	s.enqueue(syncEvent{kind: eventDelete, messageIDs: event.IDs})
}

func (s *Service) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-s.events:
			if err := s.processEvent(ctx, event); err != nil {
				s.logger.ErrorContext(
					ctx,
					"failed to process Discord notification event",
					logger.Error(err),
				)
			}
		}
	}
}

func (s *Service) processEvent(ctx context.Context, event syncEvent) error {
	switch event.kind {
	case eventBackfill:
		return s.backfill(ctx)
	case eventCreate:
		return s.syncMessage(ctx, *event.message, false)
	case eventUpdate:
		for _, messageID := range event.messageIDs {
			message, err := s.discord.Message(ctx, s.channelID, messageID)
			if err != nil {
				return fmt.Errorf("fetch updated Discord message %s: %w", messageID, err)
			}
			if err := s.syncMessage(ctx, *message, false); err != nil {
				return err
			}
		}
		return nil
	case eventDelete:
		return s.deleteMessages(ctx, event.messageIDs)
	default:
		return fmt.Errorf("unknown Discord notification event %d", event.kind)
	}
}

func (s *Service) backfill(ctx context.Context) error {
	messages, err := s.discord.Messages(ctx, s.channelID, s.historyLimit)
	if err != nil {
		return fmt.Errorf("fetch Discord notification history: %w", err)
	}

	for index := len(messages) - 1; index >= 0; index-- {
		if err := s.syncMessage(ctx, messages[index], true); err != nil {
			return err
		}
	}
	s.logger.InfoContext(
		ctx,
		"Discord notification history synchronized",
		slog.Int("messages", len(messages)),
	)
	return nil
}

func (s *Service) syncMessage(
	ctx context.Context,
	message discord.Message,
	historical bool,
) error {
	media := make([]renderedMedia, 0, len(message.Attachments))
	attachmentKeys := make([]string, 0, len(message.Attachments))
	for _, source := range mediaSources(message) {
		item, objectKey, err := s.store.persist(ctx, message.ID.String(), source)
		if err != nil {
			s.logger.WarnContext(
				ctx,
				"failed to mirror Discord attachment; using original URL",
				slog.String("message_id", message.ID.String()),
				slog.String("attachment", source.Filename),
				logger.Error(err),
			)
		}
		media = append(media, item)
		if objectKey != "" {
			attachmentKeys = append(attachmentKeys, objectKey)
		}
	}

	editorJSJSON, hasContent, err := buildEditorJS(message, media)
	if err != nil {
		_ = s.store.remove(ctx, attachmentKeys)
		return fmt.Errorf("render Discord message %s: %w", message.ID, err)
	}
	if !hasContent {
		return s.deleteMessages(ctx, []discord.MessageID{message.ID})
	}

	createdAt := message.Timestamp.Time().UTC()
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	var updatedAt *time.Time
	if message.EditedTimestamp.IsValid() {
		value := message.EditedTimestamp.Time().UTC()
		updatedAt = &value
	}

	result, err := s.repo.UpsertDiscord(ctx, notifications.UpsertDiscordInput{
		DiscordMessageID:      message.ID.String(),
		DiscordChannelID:      s.channelID.String(),
		EditorJSJSON:          &editorJSJSON,
		CreatedAt:             createdAt,
		UpdatedAt:             updatedAt,
		DiscordAttachmentKeys: attachmentKeys,
	})
	if err != nil {
		_ = s.store.remove(ctx, attachmentKeys)
		return fmt.Errorf("upsert Discord message %s: %w", message.ID, err)
	}

	if err := s.removeReplacedAttachments(
		ctx,
		result.PreviousAttachmentKeys,
		attachmentKeys,
	); err != nil {
		s.logger.ErrorContext(ctx, "failed to remove replaced attachment", logger.Error(err))
	}

	if historical && !result.Created {
		return nil
	}
	event := notificationEvent{
		ID:           result.Notification.ID,
		UserID:       result.Notification.UserID,
		Text:         result.Notification.Text,
		EditorJSJSON: result.Notification.EditorJSJSON,
		CreatedAt:    result.Notification.CreatedAt,
	}
	if err := s.wsRouter.Publish(notificationsSubscriptionKey, event); err != nil {
		return fmt.Errorf("publish Discord notification %s: %w", message.ID, err)
	}

	return nil
}

func (s *Service) removeReplacedAttachments(
	ctx context.Context,
	previous []string,
	current []string,
) error {
	currentSet := make(map[string]struct{}, len(current))
	for _, objectKey := range current {
		currentSet[objectKey] = struct{}{}
	}

	var obsolete []string
	for _, objectKey := range previous {
		if _, ok := currentSet[objectKey]; !ok {
			obsolete = append(obsolete, objectKey)
		}
	}
	return s.store.remove(ctx, obsolete)
}

func (s *Service) deleteMessages(ctx context.Context, messageIDs []discord.MessageID) error {
	ids := make([]string, len(messageIDs))
	for index, messageID := range messageIDs {
		ids[index] = messageID.String()
	}

	deleted, err := s.repo.DeleteDiscord(ctx, s.channelID.String(), ids)
	if err != nil {
		return fmt.Errorf("delete Discord notifications: %w", err)
	}
	for _, item := range deleted {
		if err := s.store.remove(ctx, item.AttachmentKeys); err != nil {
			s.logger.ErrorContext(ctx, "failed to remove deleted attachment", logger.Error(err))
		}
		if err := s.wsRouter.Publish(notificationsSubscriptionKey, notificationEvent{
			ID:        item.ID,
			CreatedAt: time.Now().UTC(),
			Deleted:   true,
		}); err != nil {
			return fmt.Errorf("publish deleted Discord notification %s: %w", item.ID, err)
		}
	}

	return nil
}
