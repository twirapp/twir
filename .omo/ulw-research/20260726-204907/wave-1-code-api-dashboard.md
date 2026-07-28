# Wave 1: API and Dashboard Backend

## Findings
- Provider registry and capability registry are separate sources of truth.
- VK is feature-gated and has a separate bot-setup provider; Kick also has special auth/token branches.
- Auth/provider paths can degrade to empty URL/config instead of typed failure.
- Commands/timers/events expose platform targeting as string arrays despite a GraphQL Platform enum.
- Public DTOs remain Twitch/Kick-profile-centric; no first-class VK profile field.
- Backend tests cover registry and binding basics but not resolver/provider fallback seams.

## EXPAND
- Trace string platform inputs to validation/storage.
- Verify OAuth empty URL/config behavior.
- Audit profile resolution and ownership fallbacks end-to-end.
