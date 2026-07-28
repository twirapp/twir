# Verification Economics

| Claim | Risk | Error cost | Verification cost/time | Chosen path | Decision | Outcome | Residual risk |
|---|---|---|---|---|---|---|---|
| C1 internal vs external message ID | high | Chat persistence failure for every non-UUID provider ID | low | Runtime log + save mapper + Postgres/ClickHouse schemas | verify | Confirmed in both persistence contracts | Correct internal/provider ID split still needs design |
| C2 capability drift | normal | UI/API hide or misstate supported features | medium | Registry-to-provider matrix + tests | verify | Confirmed; ingress itself is not gated | Capability vocabulary may intentionally be coarse |
| C3 VK active stream discovery | high | All VK bot replies fail | medium | First-party docs + runtime HTTP result | verify | Runtime 404; final first-party contract review pending | Correct catalog pagination/offline behavior pending |
| C9 unban workflow miswire | critical | Wrong automation runs; intended workflow skipped | low | Exact queue trace + temporary AST harness | verify | Confirmed reachable | No dedicated repository test |
| C10 Twitch ban suppresses Kick/VK events | high | Non-Twitch automation silently skipped | low | Exact predicate harness + binding trace | verify | Confirmed conditional | Operation-specific Twitch dependencies still need classification |
| C14 7TV cache contamination | normal | Wrong profile returned | medium | Real client call through in-process transport | verify | Confirmed reachable in profile command | Concurrency behavior still needs dedicated race test |
| C17 malformed moderation panic | critical | User-configurable response crashes bot | low | Focused subprocess/harness | verify | Wave-3 execution pending | Exact recovery boundary unknown |
| C18 Kick Unicode split | high | Invalid or oversized outbound payloads | low | Existing focused tests | verify | Two committed tests already fail | Provider validation/runtime behavior not called live |
