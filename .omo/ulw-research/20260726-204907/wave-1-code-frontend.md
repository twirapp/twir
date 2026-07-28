# Wave 1: Frontend and Overlay UX

## Findings
- Shared platform selector exposes Twitch/Kick only, blocking VK configuration in commands/events/timers/keywords.
- Auth/dashboard profile model lacks VK.
- Commands/timers retain Twitch categories/announce colors.
- Multiple profile links, embeds, chat walls, giveaways, and rewards point to Twitch URLs.
- Overlay chat can show Kick partially, but badges/emotes remain Twitch-first; TMI stamps messages as Twitch.
- Client-only browser assumptions are safe in current entrypoints but fragile under SSR reuse.

## EXPAND
- Inventory every platform selector consumer.
- Verify `resolveProfile` unknown/VK fallback.
- Build overlay message/emote/badge parity matrix.
