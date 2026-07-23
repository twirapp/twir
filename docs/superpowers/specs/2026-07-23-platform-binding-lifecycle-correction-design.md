# Platform Binding Lifecycle Correction Design

## Context

The final review of dashboard platform binding management found two regressions:

1. A logged-out OAuth callback for an already linked platform account creates a
   new channel before checking the existing binding, then fails with
   `errPlatformConflict`.
2. Disabling or disconnecting a binding does not reliably remove its EventSub
   subscriptions. Disconnect deletes the binding before the asynchronous
   listener can load it from the database.

## Goals

- Make logged-out reauthentication reuse the existing channel for the
  authenticated platform account.
- Synchronize EventSub subscriptions after successful enable, disable, and
  disconnect mutations.
- Preserve sufficient information to unsubscribe after a binding is deleted.
- Retain the existing best-effort, post-commit NATS delivery model.

## Non-goals

- Add a transactional outbox, retry worker, schema migration, or database
  migration.
- Change targeted-dashboard OAuth authorization or cross-channel conflict
  behavior.
- Remove deferred legacy platform columns from Task 13.
- Add VK EventSub provider behavior beyond the existing generic no-op path.

## OAuth Reauthentication

`completePlatformAuth` already resolves the provider identity before opening
its transaction. Within that transaction, the no-live-session and
no-target-channel branch will first query for a channel bound to
`(input.Platform, platformUser.ID)`.

- When a channel exists, it becomes the result channel. The existing
  `linkPlatformToChannel` idempotency check verifies that the binding belongs
  to the same provider account and preserves it.
- When no channel exists, the current new-channel creation path remains
  unchanged.
- Lookup errors other than `channelsrepo.ErrNotFound` abort the transaction.
- Live-session and targeted-dashboard callbacks retain their current channel
  selection and conflict protection.

The callback still refreshes the provider token, restores the session, and
sets the selected dashboard after the transaction. Existing-channel
reauthentication must not publish default-channel resources.

## EventSub Lifecycle

### Message contract

`EventsubUnsubscribeRequest` will gain an optional provider-neutral binding
snapshot containing:

- binding ID;
- provider user ID;
- provider channel ID.

The existing `ChannelID` and `Platform` fields remain for logging and for
legacy channel-level unsubscribe callers. The snapshot deliberately contains
only the identity values required for cleanup, not a repository model.

### Mutation behavior

`channel_platforms.Service` receives a narrow EventSub publisher adapter and
logger through dependency injection.

- `SetEnabled` reads and patches the binding inside `transactionRunner.Do`.
  After a successful commit, enabling publishes the existing subscribe-all
  message. Disabling publishes an unsubscribe message with the updated binding
  snapshot.
- `Disconnect` captures the binding inside its existing transaction, deletes
  it, and publishes an unsubscribe message with the captured snapshot only
  after the transaction succeeds.
- A failed database operation or transaction rollback publishes nothing.
- A failed post-commit publish is logged with channel and platform context and
  does not turn a committed GraphQL mutation into an error. This matches the
  current auth/dashboard best-effort NATS semantics.

### Listener behavior

When an unsubscribe request has a snapshot, the EventSub listener does not
load the channel from the database.

- Twitch uses the saved provider channel ID to remove its subscriptions.
- Kick reconstructs the small binding value required by its subscription
  manager, including the binding ID used for Redis subscription keys and the
  provider user ID used to resolve the external account.
- Other platforms remain no-ops until they implement EventSub support.

When a request has no snapshot, the listener retains its current channel
lookup behavior. This preserves existing account-wide/admin/dashboard callers
that request cleanup for all bindings on a channel.

## Regression Coverage

- Table-driven logged-out reauthentication for Twitch, Kick, and VK proves an
  existing binding is reused, no channel/binding is created, and token/session
  completion continues.
- `SetEnabled` tests prove transaction completion precedes publication, enable
  publishes subscribe, disable publishes a complete snapshot, and rollback
  publishes nothing.
- `Disconnect` tests prove the snapshot is captured before deletion, publish
  happens only after commit, and rollback publishes nothing.
- Listener tests prove snapshot-based Twitch and Kick cleanup works without a
  channel lookup; legacy snapshot-less requests still use the existing lookup
  path.

## Verification

- Run focused Go tests for auth, channel-platforms, EventSub listener, and
  affected bus-core packages.
- Run `go test -count=1 ./apps/api-gql/...` and the relevant EventSub package
  tests.
- Run `git diff --check` and perform a fresh targeted final review of the two
  corrected regressions.
