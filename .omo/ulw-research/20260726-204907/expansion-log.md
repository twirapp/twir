# Expansion Log

## Wave 1

- Status: complete
- Axes: platform core; ingress; identity/stats; parser/commands; variables; persistence; events/timers; bots outbound; API/dashboard; frontend/overlays; provider docs; OSS implementations; history/migrations; skeptic/test matrix.
- Leads opened: ID boundary map; cache contamination; event miswire/Twitch-ban gate; provider API mismatch; platform selector/profile fallback; capability consistency; direct stats writers; OSS conformance patterns.
- Leads closed as duplicates: repeated Twitch-default findings across parser/API/frontend; repeated VK outbound 404 finding.

## Wave 2

- Status: complete; all nine reports recovered and cross-checked.
- Targets: ID boundary; parser cache proof; event flow proof; capability matrix; API/UI fallback; provider outbound; stats writers; OSS comparison; skeptic cross-critique.
- Converged leads: provider/internal ID persistence boundary; capability underdeclaration; VK UI/profile fallback; hidden Kick 429 drops; blind stats updates.
- Counter-evidence: VK timer exclusion is intentional and test-covered; migration-level channel-platform backfills are strongly guarded; legacy VK classic API drift is separate from VK Video Live.
- Rejected lead: the raw Twitch mapper's provider IDs do not reach consumers; `processChannelChatMessage` overwrites `ChannelID` and `UserID` with internal UUIDs before publishing.

## Wave 3

- Status: preparing final adversarial expansion from wave-2 leads.
- Targets: event-ID semantics; parser race/permission/platform lookup execution; greeting/moderation/Kick split execution; VK Video Live first-party discovery contract; remediation/test architecture.
