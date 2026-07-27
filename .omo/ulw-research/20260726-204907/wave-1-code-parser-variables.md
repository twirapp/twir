# Wave 1: Parser Variables

## Findings
- Confirmed shared-cache bugs: 7TV Twitch/Kick profile values share one slot; `GetGbUserStats` shares one slot across user IDs.
- `top.emotes.users` always asks Twitch for provider IDs, breaking Kick/VK display names.
- `NewParseContextChannel` intentionally carries TwitchUserID even for another selected platform.
- Many variables correctly reject unsupported platforms, but several default empty/unknown to Twitch or return empty values.
- Current-song cache is integration-driven and includes VK; risk is precedence rather than platform gating.

## EXPAND
- Verify shared cache contamination by execution.
- Audit all TwitchUserID consumers under Kick/VK parse contexts.
- Build a variable/platform support matrix.
