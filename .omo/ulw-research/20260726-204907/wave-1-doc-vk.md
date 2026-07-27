# Wave 1: VK Contracts

## Findings
- Official docs cover WebSocket connection/subscription tokens, chat read/write/settings, channel/catalog discovery, and typed IDs.
- Official send requires `/v1/chat/message/send` with `channel_url` and `stream_id`; catalog exposes `/v1/catalog/active_channels` and `/online_channels`.
- Channel/user IDs are int64; stream IDs are strings.
- Public docs do not expose reply/moderation-action contracts or concrete WebSocket event payload docs.
- Accessible OAuth scope list is incomplete; chat send permission is `chat:message:send`.

## EXPAND
- Verify active catalog pagination and owner matching.
- Determine offline chat-send semantics empirically or mark unresolved.
- Compare implementation endpoints against documented methods.
