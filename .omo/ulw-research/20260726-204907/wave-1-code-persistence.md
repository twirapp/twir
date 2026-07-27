# Wave 1: Persistence IDs

## Findings
- `generic.ChatMessage` exposes multiple untyped string IDs; stores infer semantics.
- Confirmed runtime: numeric VK message ID entered ClickHouse UUID `chat_messages.id`.
- `ChannelBindingID` is expected to be an internal UUID and parsed in bots; provider channel ID is separate.
- API-GQL has provider-ID-to-UUID parse choke points in moderation paths.
- Redis keys are string-safe but often omit platform namespace, creating collision/replay ambiguity.
- Worker return was incomplete; exact full consumer inventory remains open.

## EXPAND
- Complete every `generic.ChatMessage` ID consumer and UUID parse map.
- Verify Redis key namespace collision risk.
- Separate internal event UUID from provider message ID at ingress.
