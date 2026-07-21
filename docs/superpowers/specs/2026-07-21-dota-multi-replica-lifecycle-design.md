# Dota Multi-Replica Lifecycle Design

## Status

Approved for implementation on 2026-07-21. Revised to PostgreSQL state and
outbox on 2026-07-21 after Redis action-fencing analysis.

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

PostgreSQL is the authoritative match-state and action-outbox store. Redis
continues to hold non-authoritative live snapshots and prediction ownership
records, but it is not the source of transition truth. The state machine no
longer caches channel state in process memory.

Each snapshot includes:

- `revision`: monotonically increasing fencing value.
- `lastProviderTimestamp` and last game-time data for stale-payload checks.
- Current match identity and game state.
- Last terminal match identity where needed to suppress replayed terminal data.

`dota_channel_match_states` has one row per channel with a JSONB snapshot,
revision, and source-order fields. Every GSI request follows one transaction:

1. Insert an idle row if absent, then select the channel row `FOR UPDATE`.
2. Reject a payload older than the snapshot source timestamp. For a matching
   match, use game time to order payloads with equal provider timestamps.
3. Require a `post_game` payload match ID to equal the tracked match ID.
4. Never abandon a tracked match from a map-less or inactive payload unless its
   source timestamp is strictly newer than the tracked state.
5. Build the next revision and zero or more lifecycle action rows.
6. Update the locked state row and insert action rows into the PostgreSQL
   outbox in the same transaction.
7. Commit before performing any Twitch or bus side effect.

A different in-game match can replace the tracked match only when it is newer
than the tracked source state. The same atomic transition emits a cancel action
for the old match before the create action for the new match.

## Durable Prediction Actions

`dota_prediction_outbox` stores `create`, `resolve`, and `cancel` actions with
their channel, match, revision/sequence, terminal result, payload, attempt
count, availability time, lock token, and completion time. A unique action key
prevents duplicate inserts for the same channel, match, and action.

Workers claim rows with `FOR UPDATE SKIP LOCKED`, assign a random lock token,
and perform Twitch work outside the claim transaction. Completion and retry
updates require that token. A stale lease can be claimed by another replica.
Failures clear the lease and move `available_at` forward with bounded backoff.

Workers validate a create action against current authoritative state before
calling Twitch. If the match has already reached a terminal state, the create
action completes as a safe no-op rather than creating a late prediction.
Terminal actions operate on the prediction record for their exact channel and
match regardless of a newer active match. Workers process only the earliest
unfinished action for a match, so a terminal action cannot overtake creation.

The Redis prediction store retains an atomic terminal claim. A worker must own
that claim before reading and ending a Twitch prediction. Successful terminal
work marks the action complete; transient failures release the claim for retry.
This is a second idempotency boundary if a database action lease expires during
an external Twitch call.

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
with the returned settings and transactionally updates the authoritative state
statistics without changing any newer match identity.

## Event Semantics

- `MatchStarted` and `MatchAbandoned` bus messages remain available for
  non-critical consumers such as chat alerts.
- Prediction correctness no longer depends on Core NATS delivery or handler
  replies.
- `MatchEnded` is emitted by the durable resolve worker after idempotent
  settlement, rather than before a durable terminal action is accepted.
- Existing GetData and API state-update behavior remains compatible.

## Failure Handling

- Database transaction conflict: release the row lock and retry the transition
  without external side effects.
- Worker crash: another replica claims the expired outbox lease.
- Twitch timeout: leave the outbox action available for retry and preserve the existing
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
- Outbox action retry, completion, sequence ordering, and stale-lease claim
  behavior.
- Create action suppression after a terminal transition.
- Settlement replay without duplicate MMR/W/L updates.
- Existing GSI, prediction, bus Gob, and chat-alert behavior.

Database migration and repository tests must validate state-row locking,
outbox uniqueness/claim completion, the unique settlement key, atomic first
application, and replay return path.
