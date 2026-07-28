# Wave 2 Execution Proofs

## Confirmed Reachable

- 7TV profile cache is shared across users/platform lookup methods. A temporary overlay test fetched one Twitch profile, then observed the same profile for a second Twitch user and a Kick user. The `7tv profile` command can seed the cache through a Twitch mention lookup before expanding platform-aware variables.
- `ChannelUnbanRequestCreate` is received and alerted under the correct type, then dispatched to workflows as `EventTypeRedemptionCreated`.
- A non-Twitch event is suppressed when the separate Twitch binding has `IsTwitchBanned=true`; no platform condition scopes that predicate.

## Confirmed But Currently Unreachable

- `GetGbUserStats` has one unkeyed cache slot and returns the first user's value for a different input. Current callers use only one target per fresh parser cacher, so no production caller currently triggers cross-user contamination.

## Execution

- Parser overlay: `go test -vet=off -race -shuffle=on -count=1 -overlay=/tmp/twir-cache-contamination-probe/overlay.json ./internal/cacher` failed with the expected three contamination assertions.
- Events focused test: `go test -race -shuffle=on -count=1 -run '^TestGetEventChannelBindingsSelectsEventAndTwitchBindingsByPlatform$' -v ./internal/workflows` passed.
- Temporary AST harness proved both event predicates and passed `go vet`; all temporary files were removed.

## Cleanliness

- No tracked source was edited by the probes.
- `git diff --check` remained clean.
- Only `.omo/ulw-research/` is untracked.
