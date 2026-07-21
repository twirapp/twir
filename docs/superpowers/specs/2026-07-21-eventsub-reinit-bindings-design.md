# EventSub Binding Reinitialization Design

## Scope

Fix `BusListener.reinitChannels` so startup reinitialization reaches every unique channel with a Twitch or Kick binding. Preserve the existing all-platform subscription lifecycle and its eligibility checks.

## Data Flow

After the existing Twitch subscription cleanup, the listener will:

1. Load every Twitch-bound channel with `GetAllByBindingPlatform(ctx, platform.PlatformTwitch)`.
2. Load every Kick-bound channel with `GetAllByBindingPlatform(ctx, platform.PlatformKick)`.
3. Deduplicate the returned channel UUIDs.
4. Invoke the existing all-platform subscription path once for each unique UUID.

An error from either binding lookup stops reinitialization before subscriptions begin, matching the previous all-or-nothing lookup behavior. Subscription errors continue to be logged per channel, as they are today.

## Structure

Extract the post-cleanup work into an unexported helper that accepts the subscription callback. `reinitChannels` passes a callback that invokes `subscribeToAllEventsByChannelID(ctx, id, "")`, so platform-specific lifecycle behavior remains unchanged. The callback seam keeps the regression test local and avoids Twitch, token, database, or Docker dependencies.

## Testing

Add a unit test with more than ten binding results across both platforms and at least one dual-bound channel. A fake repository records both platform lookups; the callback records channel IDs. The assertions require both complete binding queries, zero legacy `GetMany` calls, and exactly one callback invocation per unique channel.

## Non-Goals

- Do not change Task 6 payload handling or Task 7 lifecycle keys.
- Do not alter the Twitch cleanup phase, subscription ordering guarantees, or platform eligibility rules.
- Do not add infrastructure-backed tests.
