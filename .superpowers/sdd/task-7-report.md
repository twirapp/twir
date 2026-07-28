# Task 7 Report: Generalize Event Transport Subscription Lifecycle

## Implementation Summary

- Added a narrow `EventTransport` contract with platform, capabilities, binding-aware subscribe/unsubscribe methods, and callback base URL setup.
- Added an EventSub transport registry backed by the shared platform registry. The current composition registers the Kick webhook transport; future transports can use the same registry and lifecycle path.
- Changed the webhook lifecycle manager to enumerate registered platforms, query `GetAllByBindingPlatform`, extract the matching normalized binding, and route it through the registry.
- Preserved enabled-only subscription behavior and cleanup of disabled bindings. Batch dispatch joins individual failures only after attempting all bindings, so one failure does not stop later subscriptions.
- Migrated Kick subscription state to `kick:sub:<binding-id>:<event-type>:<transport-kind>` and introduced webhook/WebSocket transport kinds in `libs/bus-core/eventsub`.
- Updated direct Kick subscribe/unsubscribe and the resubscribe job to pass complete bindings, ensuring every state-key path has the binding ID.
- No VK protocol, webhook handler, or Task 8/Task 10 work was added.

## Files

- Added `apps/eventsub/internal/platforms/transport.go`
- Added `apps/eventsub/internal/platforms/transport_test.go`
- Added `apps/eventsub/internal/kick/subscription_manager_test.go`
- Modified `apps/eventsub/internal/webhook/manager.go`
- Modified `apps/eventsub/internal/webhook/manager_test.go`
- Modified `apps/eventsub/internal/kick/subscription_manager.go`
- Modified `apps/eventsub/internal/kick/resubscribe_job.go`
- Modified `apps/eventsub/internal/kick/resubscribe_job_test.go`
- Modified `apps/eventsub/internal/bus-listener/bus-listener.go`
- Modified `apps/eventsub/app/app.go`
- Modified `libs/bus-core/eventsub/eventsub.go`

## RED Evidence

The first lifecycle regression test was run before production changes:

```text
$ go test ./internal/webhook -run '^TestSubscribeAllPlatformsUsesRegisteredTransportPlatform$' -count=1
INFO webhook manager: subscribing to kick events channels_count=0
INFO webhook manager: finished subscribing to kick events channels_count=0
--- FAIL: TestSubscribeAllPlatformsUsesRegisteredTransportPlatform
    transport subscriptions = [], want [1cbb738c-85cf-4c8a-aa25-b46ec612d424]
    binding platform lookups = [kick], want [twitch]
FAIL
```

This proved the old lifecycle manager queried the hard-coded Kick binding list and never invoked a registered Twitch transport.

The contract and state-key tests also failed before implementation:

```text
$ go test ./internal/platforms -count=1
undefined: EventTransport
undefined: SubscribeAll
undefined: NewRegistry
FAIL [build failed]

$ go test ./internal/kick -run '^TestRedisKeyIncludesBindingEventAndTransportKind$' -count=1
too many arguments in call to redisKey
undefined: buscoreeventsub.TransportWebhook
FAIL [build failed]
```

## GREEN Evidence

Focused tests after implementation:

```text
$ go test ./internal/platforms -count=1
ok github.com/twirapp/twir/apps/eventsub/internal/platforms

$ go test ./internal/webhook -run '^(TestBulkKickOperationsUseCompleteBindingList|TestSubscribeAllPlatformsUsesRegisteredTransportPlatform)$' -count=1
ok github.com/twirapp/twir/apps/eventsub/internal/webhook

$ go test ./internal/kick -run '^(TestRedisKeyIncludesBindingEventAndTransportKind|TestResubscribeJob_)' -count=1
ok github.com/twirapp/twir/apps/eventsub/internal/kick
```

Required final verification:

```text
$ go test ./apps/eventsub/...
ok github.com/twirapp/twir/apps/eventsub/internal/bus-listener
ok github.com/twirapp/twir/apps/eventsub/internal/channelbinding
ok github.com/twirapp/twir/apps/eventsub/internal/handler
ok github.com/twirapp/twir/apps/eventsub/internal/kick
ok github.com/twirapp/twir/apps/eventsub/internal/manager
ok github.com/twirapp/twir/apps/eventsub/internal/platforms
ok github.com/twirapp/twir/apps/eventsub/internal/services/user-creator
ok github.com/twirapp/twir/apps/eventsub/internal/webhook
```

Additional shared-contract verification:

```text
$ go test ./libs/bus-core/...
ok github.com/twirapp/twir/libs/bus-core/generic
all remaining bus-core packages compiled successfully

$ git diff --check
exit 0
```

## Tests Added Or Extended

- `TestSubscribeAllRoutesBindingsToTheirMatchingTransport`
- `TestSubscribeAllSkipsDisabledBindings`
- `TestSubscribeAllContinuesAfterBindingFailure`
- `TestSubscribeAllPlatformsUsesRegisteredTransportPlatform`
- `TestRedisKeyIncludesBindingEventAndTransportKind`
- Updated resubscribe-job doubles and assertions for binding-aware `Subscribe` calls.

## Self-Review

- The lifecycle manager only queries bindings for platforms with a registered transport, then verifies the extracted binding matches that platform before dispatch.
- Subscription failures are accumulated with `errors.Join` after all enabled bindings are attempted, preserving the previous continue-on-failure behavior.
- Unsubscription intentionally does not filter disabled bindings, preserving development cleanup behavior.
- The shared bus subjects and request structs remain unchanged.
- The otherwise necessary call-site changes are limited to passing full bindings to the renamed Kick methods, which is required for binding-scoped state keys.
- No generated files, database migrations, or unrelated cleanup were included.

## Concerns

- Existing Redis keys that used the old user-ID format are not read under the new format. Existing cleanup already falls back to listing provider subscriptions when a Redis key is absent, so stale old keys do not block cleanup.
- Only Kick is registered today by design. VK event transport registration and its provider protocol remain intentionally out of scope for this task.

## Commits

- This report and the implementation are included in the local Task 7 commit: `feat(eventsub): generalize transport lifecycle`.
- Base before Task 7: `ebee94ae451cc139fab6e394cf9b66d0ac940c67`.
