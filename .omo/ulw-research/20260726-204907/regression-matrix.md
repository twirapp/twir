# Cross-Platform Regression Matrix

Legend: PASS = existing focused test passes; FAIL = existing committed test fails; GAP = no focused coverage; UNSUPPORTED = explicit/tested exclusion.

| Boundary | Twitch | Kick | VK Video Live | Required regression |
|---|---|---|---|---|
| Chat ingress normalization | partial PASS | PASS | PASS in `0c9ed8bdc` | Golden fixture asserting internal UUIDs and provider IDs separately |
| Chat persistence ID | PASS by UUID-shaped provider IDs | unverified provider format | runtime FAIL for numeric ID | Persist internal UUID plus provider message ID |
| Command routing | PASS | partial PASS | runtime reached parser | One same-command fixture per platform |
| Mentioned-user stats variables | PASS | GAP / hardcoded Twitch lookup | GAP / hardcoded Twitch lookup | Valid internal DBChannelID with non-Twitch provider ID |
| Multi-response ordering/race | GAP | GAP | GAP | Two or more responses under `-race`, deterministic configured order |
| Stat-only command permissions | GAP | GAP | GAP | Empty roles plus each unmet threshold must deny |
| 7TV profile cache identity | GAP | execution FAIL in profile command path | provider unsupported/unclear | Distinct platform+user keys in one parse |
| Event type mapping | broad GAP | broad GAP | broad GAP | Table test from queue/message type to workflow event type |
| Twitch-ban event gate | Twitch-specific behavior | execution FAIL | conditionally affected | Non-Twitch flow must not depend on Twitch ban unless operation requires Twitch |
| Chat send | PASS | reply supported | adapter test PASS | Adapter conformance with explicit sent/dropped/unsupported/error outcome |
| Rate-limit behavior | GAP; SkipRateLimits dead | intentional 429 drop | provider behavior unknown | Test bypass flag and retry/drop outcome |
| Message splitting | Twitch implementation-specific | FAIL for byte/UTF-8 tests | unverified | Unicode boundary and provider limit fixtures |
| Greetings | Twitch-only implementation | source FAIL/misroute | source FAIL/misroute | Source platform dispatch; processed state only after success |
| Moderation command parsing | panic GAP | unsupported | unsupported | Bare `/ban`, `/timeout`, whitespace, missing target |
| Delete/ban outbound | PASS implementation, limited tests | delete/ban outbound missing | UNSUPPORTED | Capability must match adapter behavior |
| Stream discovery | PASS | PASS | runtime FAIL 404 | Online, offline, unauthorized, not-found provider fixtures |
| Timers | PASS | PASS | UNSUPPORTED and tested | Preserve explicit unsupported result |
| Feature platform selector | Twitch shown | Kick shown | GAP/hidden | Derive options from typed supported-feature matrix |
| Profile URL/current profile | PASS | PASS | source FAIL → Twitch fallback | VK URL/profile mapping test |
| Chat overlay branding | PASS | PASS | GAP/falls through | Platform icon/badge exhaustive test |
| Capability exposure | partial | underdeclared chat.read | underdeclared chat.read/streams.read | Registry/transport/adapter/UI conformance table |
| Watched stat tick | PASS steady state | n/a | n/a | Concurrent/start-first duplicate tick protection |

## Minimum CI Gates

- `go test -race` for parser response generation, event workflow mapping, and bot dispatch packages.
- Provider golden fixtures for all ingress normalizers, including numeric and non-UUID provider IDs.
- Adapter conformance test shared by Twitch/Kick/VK implementations.
- GraphQL/platform-selector exhaustiveness test against the canonical platform capability table.
- No implicit default-to-Twitch branch without an explicit legacy-input test.
