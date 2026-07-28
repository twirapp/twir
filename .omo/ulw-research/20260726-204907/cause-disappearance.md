# Cause Disappearance Ledger

| cause_id | Expected truth | Previous observation | Last seen | Disconfirming observation | Replacement cause | Current status | Violation gone? |
|---|---|---|---|---|---|---|---|
| CA1 | VK ingress accepts DevAPI event wrapper. | Parser expected legacy `type=message`. | 2026-07-26 20:41 | Live event reached parser/bots at 20:43 after wrapper fix. | Downstream ID and outbound API mismatches. | replaced | yes |
| CA2 | Provider IDs reach UUID chat consumer fields. | Raw Twitch mapper temporarily assigns provider IDs to `ChannelID`/`UserID`. | Wave 2 initial map | `processChannelChatMessage` overwrites both with internal UUIDs before publication. | Persistence `ID`, not channel/user IDs, remains mismatched for VK. | disproved | yes |
| CA3 | `GetGbUserStats` leaks between production parser targets. | Method-level unkeyed cache returns first user for a second input. | Wave 2 execution | Every current caller uses one target per fresh cacher. | Twitch-hardcoded `getDbChannel` is the reachable Kick/VK stats defect. | replaced | yes |
| CA4 | Kick 429 disappears silently by accident. | Adapter returns nil on rate limit. | Wave 2 provider audit | Adapter explicitly logs `dropping message`; policy is intentional. | Outcome contract cannot tell callers or metrics that delivery was dropped. | narrowed | no |
| CA5 | VK is accidentally omitted from all timers/events UI. | Shared selector has only Twitch/Kick. | Wave 2 UI audit | Timer implementation and tests explicitly exclude VK; VK advertises no event capability. | Commands, keywords, chat filtering, profile URLs, and chat branding remain fragmented. | narrowed | partial |
