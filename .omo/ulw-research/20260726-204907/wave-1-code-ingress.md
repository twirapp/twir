# Wave 1: Chat Ingress

## Findings
- `generic.ChatMessage` is Twitch-shaped and richer than Kick/VK normalization.
- Confirmed Twitch mapping defects: delete target login/name swap; ban broadcaster login uses display name; `community_sub_gift` no-op; unknown fragments become text.
- VK drops non-text parts and timestamps; Kick lacks reply metadata.
- EventSub reinit path reloads Twitch/Kick but omits VK.
- Capability declarations do not match implemented ingress: capability tables are advisory.

## EXPAND
- Complete `generic.ChatMessage` field-to-consumer map.
- VK reinit lifecycle and capability mismatch.
- Persistence/API truncation of provider metadata.
