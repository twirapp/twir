package manager

import (
	"context"
	"log/slog"
	"time"

	"github.com/kvizyx/twitchy/eventsub"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// todo: signal for shutdown
func (c *Manager) startWebSocket() {
	wsCtx := context.TODO()
	ws := c.eventsub.Websocket()

	ws.OnWelcome(
		func(message eventsub.WebsocketWelcomeMessage) {
			c.logger.Info(
				"websocket welcome received",
				slog.String("session_id", message.Payload.Session.Id),
			)
			c.wsCurrentSessionId = &message.Payload.Session.Id

			if err := c.twitchUpdateConduitShard(context.Background()); err != nil {
				c.logger.Error("failed to update conduit shard", slog.Any("err", err))
			} else {
				c.logger.Info(
					"conduit shard updated",
					slog.String("conduit_id", c.currentConduit.Id),
				)
			}
		},
	)

	ws.OnChannelChatMessage(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelChatMessage))
	ws.OnChannelChatMessageDelete(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelChatMessageDelete))
	ws.OnChannelChatNotification(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelChatNotification))
	ws.OnChannelChatClear(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelChatClear))
	ws.OnChannelFollow(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelFollow))
	ws.OnChannelModeratorAdd(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelModeratorAdd))
	ws.OnChannelModeratorRemove(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelModeratorRemove))
	ws.OnChannelPointsCustomRewardRedemptionAdd(
		wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelPointsRewardRedemptionAdd),
	)
	ws.OnChannelPointsCustomRewardRedemptionUpdate(
		wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelPointsRewardRedemptionUpdate),
	)
	ws.OnChannelPointsCustomRewardRemove(
		wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelPointsRewardRemove),
	)
	ws.OnChannelPointsCustomRewardAdd(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelPointsRewardAdd))
	ws.OnChannelPointsCustomRewardUpdate(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelPointsRewardUpdate))
	ws.OnChannelPollBegin(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelPollBegin))
	ws.OnChannelPollProgress(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelPollProgress))
	ws.OnChannelPollEnd(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelPollEnd))
	ws.OnChannelPredictionBegin(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelPredictionBegin))
	ws.OnChannelPredictionProgress(
		wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelPredictionProgress),
	)
	ws.OnChannelPredictionLock(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelPredictionLock))
	ws.OnChannelPredictionEnd(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelPredictionEnd))
	ws.OnChannelBan(wrapWsEventsubHandlerWithCtx(c.handler.HandleBan))
	ws.OnChannelSubscribe(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelSubscribe))
	ws.OnChannelSubscriptionMessage(
		wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelSubscriptionMessage),
	)
	ws.OnChannelRaid(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelRaid))
	ws.OnChannelUnbanRequestCreate(
		wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelUnbanRequestCreate),
	)
	ws.OnChannelUnbanRequestResolve(
		wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelUnbanRequestResolve),
	)
	ws.OnChannelVipAdd(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelVipAdd))
	ws.OnChannelVipRemove(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelVipRemove))
	ws.OnUserAuthorizationRevoke(
		wrapWsEventsubHandlerWithCtx(c.handler.HandleUserAuthorizationRevoke),
	)
	ws.OnChannelUpdate(wrapWsEventsubHandlerWithCtx(c.handler.HandleChannelUpdate))
	ws.OnStreamOnline(wrapWsEventsubHandlerWithCtx(c.handler.HandleStreamOnline))
	ws.OnStreamOffline(wrapWsEventsubHandlerWithCtx(c.handler.HandleStreamOffline))
	ws.OnUserUpdate(wrapWsEventsubHandlerWithCtx(c.handler.HandleUserUpdate))

	for {
		if err := ws.Connect(wsCtx); err != nil {
			c.logger.Error("websocket connection failed", slog.Any("err", err))
		}

		if err := ws.Disconnect(); err != nil {
			c.logger.Warn("websocket disconnection failed", slog.Any("err", err))
		}

		time.Sleep(1 * time.Second)
	}
}

var tracer = otel.Tracer("eventsub.websocket")

func wrapWsEventsubHandlerWithCtx[Event any](
	handler func(
		context.Context,
		Event,
		eventsub.WebsocketNotificationMetadata,
	),
) eventsub.Handler[Event, eventsub.WebsocketNotificationMetadata] {
	return func(event Event, metadata eventsub.WebsocketNotificationMetadata) {
		ctx := context.Background()
		newCtx, span := tracer.Start(ctx, "eventsub.websocket.handler")
		defer span.End()
		ctx = trace.ContextWithSpan(newCtx, span)

		span.SetAttributes(
			attribute.String(
				"twitcheventsub.message_type",
				metadata.MessageType,
			),
			attribute.String(
				"twitcheventsub.subscription_type",
				string(metadata.SubscriptionType),
			),
			attribute.String(
				"twitcheventsub.subscription_version",
				metadata.SubscriptionVersion,
			),
		)

		handler(ctx, event, metadata)
	}
}
