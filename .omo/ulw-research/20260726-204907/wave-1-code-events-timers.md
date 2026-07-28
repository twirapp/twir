# Wave 1: Events, Timers, Scheduler

## Findings
- Confirmed likely miswire: `ChannelUnbanRequestCreate` is emitted as `RedemptionCreated`.
- Generic event flow applies Twitch-ban state before platform-specific execution, potentially suppressing non-Twitch work.
- Scheduler stream synchronization reuses a stale snapshot during concurrent DB mutations/publishes.
- Follow/stream chat alerts are Twitch-specific and empty source defaults to Twitch.
- Timers intentionally support Twitch/Kick only; VK is skipped and tests encode that limitation.

## EXPAND
- Verify event miswire and Twitch-ban cross-platform suppression with focused tests.
- Execute stream transition scenarios against `processStreams`.
- Audit event activity helpers and chat-alert send targets.
