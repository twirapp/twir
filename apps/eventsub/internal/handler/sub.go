package handler

import (
	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
)

func (c *Handler) handleChannelSubscribe(
	h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventChannelSubscribe,
) {

}

func (c *Handler) handleChannelSubscriptionGift(
	h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventChannelSubscriptionGift,
) {

}

func (c *Handler) handleChannelSubscriptionMessage(
	h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventChannelSubscriptionMessage,
) {

}
