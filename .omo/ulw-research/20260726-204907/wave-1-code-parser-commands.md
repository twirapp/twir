# Wave 1: Parser Commands

## Findings
- Empty command platform allowlist means all platforms; default command creation persists no restriction.
- Default/help prefix remains hardcoded `!` in argument usage output.
- Response filters remain Twitch-category-specific.
- Potential nil dereference on missing default-command handler and a concurrent append race in response parsing.
- Stat-only restrictions may be bypassed when roles are empty.
- Twitch-only command bodies remain in VIP, voteban, and duel paths.
- Async response sends and dedupe errors are log-only/best-effort.

## EXPAND
- Execute focused race/permission/default-command proofs.
- Matrix all default commands against platform dependencies.
- Verify command prefix and response-order behavior.
