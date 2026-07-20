# Dota Multi-Replica Lifecycle Design

## Status

Approved for implementation on 2026-07-21.

## Problem

`apps/dota` runs six replicas, but its match state is currently restored once
into process-local memory and then written back to Redis without a revision
check. A stale replica can therefore abandon or settle a newer match. Core
NATS lifecycle messages are also fire-and-forget, so prediction create,
resolve, and cancel failures are not retried durably. Finally, match-result
updates are not idempotent across crash or replay.

## Goals

- Accept each GSI state transition once across all Dota replicas.
- Reject delayed or mismatched GSI payloads before they change state.
- Deliver prediction lifecycle actions at least once with durable retry.
- Make MMR and session W/L settlement idempotent per channel and Dota match.
- Ensure only one terminal prediction action can win for a match.
- Preserve existing Dota bus events for chat alerts and API state updates.

## Non-Goals

- Change existing global Core NATS queue semantics.
- Reduce Dota service replicas.
- Add a general-purpose project-wide outbox framework.
- Change Twitch prediction titles beyond the already approved managed marker.

## Authoritative Match State

Redis is the authoritative match-state store. The state machine no longer
caches channel state in process memory.

Each snapshot includes:

- `revision`: monotonically increasing fencing value.
- `lastProviderTimestamp` and last game-time data for stale-payload checks.
- Current match identity and game state.
- Last terminal match identity where needed to suppress replayed terminal data.

Every GSI request follows this loop:

1. Load the current Redis snapshot.
2. Reject a payload older than the snapshot source timestamp. For a matching
   match, use game time to order payloads with equal provider timestamps.
3. Require a `post_game` payload match ID to equal the tracked match ID.
4. Never abandon a tracked match from a map-less or inactive payload unless its
   source timestamp is strictly newer than the tracked state.
5. Build the next snapshot and lifecycle action, if any.
6. Atomically compare the expected serialized snapshot, write the next
   revision, and append the lifecycle action to a Redis Stream using one Lua
   script.
7. Retry from step 1 when the compare-and-swap loses to another replica.

A different in-game match can replace the tracked match only when it is newer
than the tracked source state. The same atomic transition emits a cancel action
for the old match before the create action for the new match.

## Durable Prediction Actions

Lifecycle actions are stored in a Redis Stream in the state-transition Lua
script. Actions are `create`, `resolve`, and `cancel` and carry the channel ID,
match ID, state revision, and terminal result when applicable.

Workers form a `dota-predictions` consumer group:

- New entries are read with `XREADGROUP`.
- An entry is acknowledged with `XACK` only after its action succeeds or is a
  safe idempotent no-op.
- Failed entries stay pending for retry.
- Each replica periodically uses `XAUTOCLAIM` to recover work owned by a dead
  consumer.
- Stream retention is bounded only after all consumer-group requirements are
  satisfied.

Workers validate a create action against current authoritative state before
calling Twitch. If the match has already reached a terminal state, the create
action is acknowledged without creating a late prediction. Terminal actions
operate on the prediction record for their exact channel and match regardless
of the currently active newer match.

The prediction store adds an atomic terminal claim. A worker must own that
claim before reading and ending a Twitch prediction. Successful terminal work
marks the action complete; transient failures release the claim for retry. This
prevents concurrent resolve and cancel workers from both ending the same
prediction.

## Idempotent Settlement

Postgres gains `dota_match_settlements` with primary key `(channel_id,
match_id)`. It records the immutable result and makes settlement idempotent.

`ApplyMatchResultOnce` executes in one database transaction:

1. Insert the match settlement with `ON CONFLICT DO NOTHING`.
2. Only when the insert succeeds, apply the MMR delta and W/L increment to
   `channels_dota_settings`.
3. Return the resulting settings for both first processing and replay.

The resolve worker applies this settlement before resolving the prediction.
Retries can therefore repeat the complete action without double-counting MMR
or W/L. After settlement, it publishes the existing `MatchEnded` chat/API data
with the returned settings. The state snapshot receives the settled MMR and
session values through a revision-fenced update so overlays remain current even
if a new match has already begun.

## Event Semantics

- `MatchStarted` and `MatchAbandoned` bus messages remain available for
  non-critical consumers such as chat alerts.
- Prediction correctness no longer depends on Core NATS delivery or handler
  replies.
- `MatchEnded` is emitted by the durable resolve worker after idempotent
  settlement, rather than before a durable terminal action is accepted.
- Existing GetData and API state-update behavior remains compatible.

## Failure Handling

- Lost Redis CAS: reload and retry without external side effects.
- Worker crash: another replica claims the unacknowledged Stream entry.
- Twitch timeout: leave the stream entry pending and preserve the existing
  managed prediction correlation record for safe recovery.
- Database retry: settlement ledger returns the prior result without applying
  the delta again.
- Mismatched or stale GSI data: ignore it without mutating state or emitting a
  lifecycle action.

## Verification

Tests must cover:

- Two independent state-machine instances racing on one channel.
- Delayed post-game and map-less payloads after a newer match.
- Exactly one emitted terminal action under concurrent resolve/cancel paths.
- Stream action retry, acknowledgement, and stale-consumer claim behavior.
- Create action suppression after a terminal transition.
- Settlement replay without duplicate MMR/W/L updates.
- Existing GSI, prediction, bus Gob, and chat-alert behavior.

Database migration and repository tests must validate the unique settlement
key, atomic first application, and replay return path.
