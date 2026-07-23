# VK Video Live Platform Dashboard Bindings Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Let dashboard users view, connect, disconnect, and enable normalized channel platform bindings from Bot Settings.

**Architecture:** The backend exposes options and bindings through the already-secured GraphQL contract. A dashboard feature-local API module owns one query and the three mutations; a composable joins registered options with existing bindings and initiates OAuth using the server-returned URL. Presentation is generic: a centralized metadata map supplies names and Lucide icon names, while one card component renders every platform without provider-specific template branches.

**Tech Stack:** Nuxt 4, Vue 3 Composition API, urql, generated GraphQL typed documents, shadcn-nuxt/Reka UI primitives, Nuxt Icon, Vitest, happy-dom, Bun.

## Global Constraints

- Use Bun for installs, code generation, tests, and builds; do not introduce Node.js commands.
- Do not edit `web/app/gql/*`; regenerate it with `bun run --cwd web graphql-codegen`.
- Keep all dashboard binding code in `web/layers/dashboard/features/channel-platforms/` except the Bot Settings mount, test configuration, and package scripts.
- Render platforms from `channelPlatformOptions` and `channelPlatformBindings`; only the centralized presentation map may mention a provider by name.
- Treat the `channelPlatformConnect` result as the only OAuth URL and navigate with `window.location.assign`; do not construct URLs or add callback routes.
- This OAuth redirect restriction applies only to binding OAuth. Preserve the existing generic `authLink(redirectTo:)` flow as a user-approved legacy exception outside Task 12 scope.
- Always retain Connect and Disconnect actions, and show the binding-wide enabled switch for every connected platform. The current GraphQL contract has no capability-specific mutation; render capabilities as data-driven information rather than inventing a disabled control.
- Use existing `Card`, `Badge`, `Button`, `Switch`, and `AlertDialog` components, plus `<Icon name="lucide:...">` for platform-neutral UI affordances.
- Preserve the two intended dependency changes currently unstaged in `web/package.json` and `bun.lock` while adding the test command/configuration.

---

## File Map

- `apps/api-gql/internal/delivery/gql/schema/{channels,user}.graphql`: Existing Task 12 permissions and binding contract; verify but do not change unless fresh verification exposes a regression.
- `apps/api-gql/internal/delivery/http/routes/auth/{oauth-platform,post-platform-code}.go`: Existing trusted callback redirect and callback authorization; verify but do not change unless fresh verification exposes a regression.
- `web/package.json`: Replace Node-shebang package shims with Bun-direct JavaScript entry points and add the focused test script alongside the already-added Vitest dependencies.
- `web/vitest.config.ts`: Configure Nuxt-aware Vitest with happy-dom and dashboard feature spec discovery.
- `web/layers/dashboard/features/channel-platforms/api.ts`: Define the binding/options query, connect/disconnect/enable mutations, and a single urql invalidation key.
- `web/layers/dashboard/features/channel-platforms/composables/use-channel-platforms.ts`: Merge options and bindings into ordered view models and perform mutations/OAuth navigation.
- `web/layers/dashboard/features/channel-platforms/ui/platform-binding-card.vue`: Render one binding state and emit generic user intents.
- `web/layers/dashboard/features/channel-platforms/ui/platform-bindings.vue`: Query-backed card list, loading/error presentation, and mutation feedback.
- `web/layers/dashboard/features/channel-platforms/ui/platform-binding-card.spec.ts`: Cover disconnected and connected card states, including a capability-less binding that retains its generic enabled control.
- `web/layers/dashboard/features/channel-platforms/ui/platform-bindings.spec.ts`: Cover OAuth initiation and generic card action wiring.
- `web/layers/dashboard/features/bot-settings/bot-settings.vue`: Mount the new bindings section above the existing Commands Prefix form.

### Task 1: Reverify the secured backend contract

**Files:**
- Verify: `apps/api-gql/internal/delivery/gql/schema/channels.graphql`
- Verify: `apps/api-gql/internal/delivery/gql/schema/user.graphql`
- Verify: `apps/api-gql/internal/delivery/http/routes/auth/post-platform-code.go`
- Verify: `apps/api-gql/internal/delivery/http/routes/auth/oauth-platform.go`
- Test: `apps/api-gql/internal/delivery/gql/directives/*platform*_execution_test.go`
- Test: `apps/api-gql/internal/delivery/http/routes/auth/oauth-platform_test.go`

**Consumes:** Commits `c8fa716ae` and `797fc310e`.

**Produces:** Fresh evidence that read actions require `VIEW_BOT_SETTINGS`, all binding mutations including legacy unlink require `MANAGE_BOT_SETTINGS`, and targeted OAuth rechecks Manage permission before provider work.

- [ ] **Step 1: Regenerate GraphQL resolver contracts**

Run: `bun cli build gql`

Expected: exit status 0 and no generated-source changes tracked by Git.

- [ ] **Step 2: Run the exact permission regressions**

Run: `go test -count=1 ./apps/api-gql/internal/delivery/gql/directives -run 'Test(ChannelPlatformOptionsGraphQLRequiresViewBotSettings|ChannelPlatformConnectGraphQLRequiresManageBotSettings|UnlinkPlatformAccountGraphQLRequiresManageBotSettings)$'`

Expected: PASS; a view-only collaborator cannot invoke Connect or legacy Unlink.

- [ ] **Step 3: Run the targeted callback regression**

Run: `go test -count=1 ./apps/api-gql/internal/delivery/http/routes/auth -run '^TestCompletePlatformCodeRejectsRevokedTargetedManagePermissionBeforeProviderWork$'`

Expected: PASS; the test asserts zero callback writes, provider calls, and transactions after permission revocation.

- [ ] **Step 4: Run the API-GQL suite and whitespace check**

Run: `go test -count=1 ./apps/api-gql/...`

Expected: PASS.

Run: `git diff --check`

Expected: no output.

- [ ] **Step 5: Review the task-scoped diff before frontend work**

Inspect: `git diff 6ba818a3e..HEAD -- apps/api-gql/internal/delivery/gql apps/api-gql/internal/delivery/http/routes/auth apps/api-gql/internal/services/channel_platforms`

Expected: no caller-controlled redirect in the binding OAuth path, no provider work before the callback Manage-permission check, and no binding mutation path without the Manage directive. The unchanged generic `authLink(redirectTo:)` flow is the approved legacy exception.

### Task 2: Establish Bun-native frontend generation and test harness

**Files:**
- Modify: `web/package.json`
- Modify: `bun.lock`
- Create: `web/vitest.config.ts`
- Generate: `web/app/gql/*` (ignored)

**Consumes:** Existing GraphQL schema type names: `ChannelPlatformBinding`, `ChannelPlatformOption`, `Platform`, and `PlatformCapability`.

**Produces:** A repeatable Bun command for component tests and generated typed documents that recognize the new binding contract.

- [ ] **Step 1: Write the failing card test before the component exists**

Create `web/layers/dashboard/features/channel-platforms/ui/platform-binding-card.spec.ts` with a minimal mount of the planned component:

```ts
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import PlatformBindingCard from './platform-binding-card.vue'

describe('PlatformBindingCard', () => {
	it('offers OAuth connection for a disconnected platform', () => {
		const wrapper = mount(PlatformBindingCard, {
			props: {
				platform: 'TWITCH',
				presentation: { label: 'Twitch', icon: 'lucide:radio' },
				capabilities: [{ name: 'chat.write' }],
				binding: null,
			},
		})

		expect(wrapper.get('button').text()).toContain('Connect')
	})
})
```

- [ ] **Step 2: Replace Node-shebang package scripts and add the Nuxt-aware Vitest configuration**

Create `web/vitest.config.ts`:

```ts
import { defineVitestConfig } from '@nuxt/test-utils/config'

export default defineVitestConfig({
	test: {
		environment: 'happy-dom',
		include: ['layers/dashboard/features/**/*.spec.ts'],
	},
})
```

The worktree has no `node` executable, and Bun 1.3.14 still invokes a locally resolved binary's `#!/usr/bin/env node` shebang when it is launched through `bunx --bun`. Use Bun to execute each package's JavaScript entry point directly. Set these scripts in `web/package.json`:

```json
{
  "build": "bun ./node_modules/nuxt/bin/nuxt.mjs build",
  "dev": "bun --env-file=../.env ./node_modules/nuxt/bin/nuxt.mjs dev --no-fork",
  "start": "bun .output/server/index.mjs",
  "generate": "bun ./node_modules/nuxt/bin/nuxt.mjs generate",
  "preview": "bun ./node_modules/nuxt/bin/nuxt.mjs preview",
  "nuxt-prepare": "bun ./node_modules/nuxt/bin/nuxt.mjs prepare",
  "prebuild": "bun run nuxt-prepare",
  "graphql-codegen": "bun ./node_modules/@graphql-codegen/cli/cjs/bin.js --config codegen.ts",
  "shadcn-vue": "bun ./node_modules/shadcn-vue/dist/index.js",
  "test": "bun ./node_modules/vitest/vitest.mjs run"
}
```

Keep the exact development dependencies already introduced: `vitest@4.1.10`, `@nuxt/test-utils@4.0.3`, `@vue/test-utils@2.4.11`, and `happy-dom@20.11.1`.

- [ ] **Step 3: Confirm the test fails for the missing component**

Run: `bun run --cwd web test -- layers/dashboard/features/channel-platforms/ui/platform-binding-card.spec.ts`

Expected: FAIL with module resolution for `platform-binding-card.vue`.

- [ ] **Step 4: Verify the Bun-direct Nuxt preparation and regenerate frontend GraphQL documents**

Run: `bun run --cwd web nuxt-prepare`

Expected: exit status 0, `.nuxt/tsconfig.json` exists, and the Nuxt pre-build GraphQL hook succeeds through the updated `graphql-codegen` script.

Run: `bun run --cwd web graphql-codegen`

Expected: exit status 0 and `web/app/gql/graphql.ts` includes `Platform`, `ChannelPlatformBinding`, and `ChannelPlatformOption` types.

- [ ] **Step 5: Check dependency and generated-artifact boundaries**

Run: `git status --short`

Expected: only intentional package/test configuration, implementation, and ignored generated-artifact changes are present; do not stage `web/app/gql/*` if it remains ignored.

### Task 3: Add the generic binding API and action composable

**Files:**
- Create: `web/layers/dashboard/features/channel-platforms/api.ts`
- Create: `web/layers/dashboard/features/channel-platforms/composables/use-channel-platforms.ts`
- Test: `web/layers/dashboard/features/channel-platforms/ui/platform-bindings.spec.ts`

**Consumes:** Generated `graphql` tag from `~/gql/gql.js`, `Platform` types from `~/gql/graphql.js`, `useMutation` from `~~/layers/dashboard/composables/use-mutation.js`, and the backend operations below.

**Produces:** `useChannelPlatforms()` with `cards`, `fetching`, `error`, `connect(platform)`, `disconnect(platform)`, and `setEnabled(platform, enabled)`.

- [ ] **Step 1: Write the failing OAuth action test**

In `platform-bindings.spec.ts`, mock `useChannelPlatforms` to return a disconnected Twitch card and a `connect` spy. Mount `platform-bindings.vue`, click its `Connect` button, and assert the spy receives `'TWITCH'`:

```ts
expect(connect).toHaveBeenCalledWith('TWITCH')
```

Also add a direct composable assertion using mocked mutation results: when Connect returns `{ channelPlatformConnect: 'https://id.example/authorize' }`, it calls `window.location.assign('https://id.example/authorize')`; when the result has `error`, it does not navigate.

- [ ] **Step 2: Define the complete generated GraphQL surface in `api.ts`**

Use one query:

```graphql
query ChannelPlatforms {
	channelPlatformBindings {
		id
		platform
		userId
		platformChannelId
		enabled
		platformUserId
		platformLogin
		platformDisplayName
		platformAvatar
		capabilities { name }
	}
	channelPlatformOptions {
		platform
		capabilities { name }
	}
}
```

Use these mutations:

```graphql
mutation ChannelPlatformConnect($platform: Platform!) {
	channelPlatformConnect(platform: $platform)
}

mutation ChannelPlatformDisconnect($platform: Platform!) {
	channelPlatformDisconnect(platform: $platform)
}

mutation ChannelPlatformSetEnabled($platform: Platform!, $enabled: Boolean!) {
	channelPlatformSetEnabled(platform: $platform, enabled: $enabled) {
		id
		platform
		enabled
	}
}
```

Export `channelPlatformsInvalidationKey = 'ChannelPlatforms'`. Apply it to the query context and pass it to every mutation through the existing wrapper, matching `web/layers/dashboard/api/commands-prefix.ts`.

- [ ] **Step 3: Implement the data-to-card composable**

Define the centralized presentation map in `use-channel-platforms.ts`:

```ts
const platformPresentation = {
	TWITCH: { label: 'Twitch', icon: 'lucide:radio' },
	KICK: { label: 'Kick', icon: 'lucide:circle-play' },
	VK_VIDEO_LIVE: { label: 'VK Video Live', icon: 'lucide:video' },
} satisfies Record<Platform, { label: string; icon: string }>
```

Build `cards` by iterating `channelPlatformOptions` in server order and locating a same-platform binding. Each card contains `platform`, `presentation`, option capabilities, and `binding | null`. Do not iterate bindings alone because a disconnected registered platform must be visible.

Implement action behavior:

```ts
async function connect(platform: Platform) {
	const result = await connectMutation.executeMutation({ platform })
	if (result.error || !result.data?.channelPlatformConnect) return result.error
	window.location.assign(result.data.channelPlatformConnect)
}
```

Implement `disconnect` and `setEnabled` with their generated mutation variables. On a successful non-navigation mutation, call the query's `executeQuery({ requestPolicy: 'network-only' })` so the cards immediately reflect the authoritative server state. Return the urql error to the UI rather than swallowing it.

- [ ] **Step 4: Run the focused action test**

Run: `bun run --cwd web test -- layers/dashboard/features/channel-platforms/ui/platform-bindings.spec.ts`

Expected: PASS for button-to-connect wiring, successful OAuth navigation, and no navigation on a mutation error.

- [ ] **Step 5: Commit the test harness and data layer**

```bash
git add web/package.json bun.lock web/vitest.config.ts web/layers/dashboard/features/channel-platforms/api.ts web/layers/dashboard/features/channel-platforms/composables/use-channel-platforms.ts web/layers/dashboard/features/channel-platforms/ui/platform-bindings.spec.ts
git commit -m "feat(web): add platform binding data"
```

### Task 4: Render data-driven binding cards and mount them

**Files:**
- Create: `web/layers/dashboard/features/channel-platforms/ui/platform-binding-card.vue`
- Create: `web/layers/dashboard/features/channel-platforms/ui/platform-bindings.vue`
- Modify: `web/layers/dashboard/features/bot-settings/bot-settings.vue`
- Test: `web/layers/dashboard/features/channel-platforms/ui/platform-binding-card.spec.ts`
- Test: `web/layers/dashboard/features/channel-platforms/ui/platform-bindings.spec.ts`

**Consumes:** `useChannelPlatforms()` from Task 3 and its card view model.

**Produces:** A visible Platform Bindings section that handles all registered providers uniformly.

- [ ] **Step 1: Expand the card test with connected and capability-independent states**

Add these expectations:

```ts
it('shows the connected account and a disconnect action', () => {
	const wrapper = mount(PlatformBindingCard, {
		props: {
			platform: 'TWITCH',
			presentation: { label: 'Twitch', icon: 'lucide:radio' },
			capabilities: [{ name: 'chat.write' }],
			binding: connectedTwitchBinding,
		},
	})

	expect(wrapper.text()).toContain(connectedTwitchBinding.platformDisplayName)
	expect(wrapper.text()).toContain('Disconnect')
})

it('keeps the generic enabled switch for a capability-less binding', () => {
	const wrapper = mount(PlatformBindingCard, {
		props: {
			platform: 'KICK',
			presentation: { label: 'Kick', icon: 'lucide:circle-play' },
			capabilities: [],
			binding: connectedKickBinding,
		},
	})

	expect(wrapper.text()).toContain('Enable bot')
})
```

Use an `AlertDialog` confirmation before emitting `disconnect`; assert that confirming emits `disconnect` with the platform. For any connected card, assert that the switch is rendered and its model update emits `set-enabled` with the new boolean. Also assert that capabilities render as badges without a provider-name conditional.

- [ ] **Step 2: Confirm the focused card test still fails**

Run: `bun run --cwd web test -- layers/dashboard/features/channel-platforms/ui/platform-binding-card.spec.ts`

Expected: FAIL until the card provides connected state, confirmation, capability badges, and the generic enabled status control.

- [ ] **Step 3: Implement `platform-binding-card.vue` as a generic intent emitter**

Use explicit typed props and emits:

```ts
const props = defineProps<{
	platform: Platform
	presentation: { label: string; icon: string }
	capabilities: { name: string }[]
	binding: ChannelPlatformBinding | null
	busy?: boolean
}>()

const emit = defineEmits<{
	connect: [platform: Platform]
	disconnect: [platform: Platform]
	setEnabled: [platform: Platform, enabled: boolean]
}>()
```

Render `Card`, `CardHeader`, `CardContent`, and `CardFooter`. Use `<Icon :name="presentation.icon" />` for the platform mark, use the avatar when present, and render the profile display name/login when connected. Render capability names as `Badge` elements. For a disconnected binding, emit Connect. For every connected binding, place Disconnect behind existing `AlertDialog` primitives and render the `Switch` labelled `Enable bot`. No `v-if` may inspect `platform` values or individual capability names.

- [ ] **Step 4: Implement `platform-bindings.vue` as the feature container**

Call `useChannelPlatforms()` once. Render a `Platform bindings` section title, loading state, mutation/query error text, and a responsive grid of `PlatformBindingCard` instances. Wire card events to `connect`, `disconnect`, and `setEnabled`. Use `toast.error('Unable to update platform binding')` only after an action returns an error; do not show success toasts immediately before OAuth navigation.

- [ ] **Step 5: Mount the section in Bot Settings**

Add this import to `web/layers/dashboard/features/bot-settings/bot-settings.vue`:

```ts
import PlatformBindings from '~~/layers/dashboard/features/channel-platforms/ui/platform-bindings.vue'
```

Add `<PlatformBindings />` before `<CommandsPrefix />` inside the existing `flex flex-col gap-12` content wrapper.

- [ ] **Step 6: Run focused frontend tests**

Run: `bun run --cwd web test -- layers/dashboard/features/channel-platforms`

Expected: PASS for disconnected/connected cards, generic enabled state for a capability-less binding, disconnect confirmation, and OAuth initiation.

- [ ] **Step 7: Commit the generic UI**

```bash
git add web/layers/dashboard/features/channel-platforms/ui/platform-binding-card.vue web/layers/dashboard/features/channel-platforms/ui/platform-bindings.vue web/layers/dashboard/features/channel-platforms/ui/platform-binding-card.spec.ts web/layers/dashboard/features/channel-platforms/ui/platform-bindings.spec.ts web/layers/dashboard/features/bot-settings/bot-settings.vue
git commit -m "feat(web): manage platform bindings"
```

### Task 5: Verify, review, and record residual environment limits

**Files:**
- Verify: all Task 12 paths above
- Modify if needed: `.superpowers/sdd/task-12-report.md` (ignored evidence only)

**Consumes:** Completed backend and frontend implementation.

**Produces:** Fresh, task-scoped verification evidence and a reviewed local commit set.

- [ ] **Step 1: Run generation and frontend verification from clean generated output**

Run: `bun run --cwd web graphql-codegen`

Expected: exit status 0.

Run: `bun run --cwd web test -- layers/dashboard/features/channel-platforms`

Expected: PASS.

Run: `bun run --cwd web build`

Expected: exit status 0. If an existing unrelated build failure occurs, preserve its full command output and distinguish it from Task 12 paths.

- [ ] **Step 2: Re-run backend verification after frontend work**

Run: `bun cli build gql`

Expected: exit status 0.

Run: `go test -count=1 ./apps/api-gql/...`

Expected: PASS.

- [ ] **Step 3: Check format, worktree scope, and diff health**

Run: `git diff --check`

Expected: no output.

Run: `git status --short`

Expected: only Task 12 source, test, package/lock, and plan/report changes; leave ignored generated artifacts untracked.

- [ ] **Step 4: Conduct a task-scoped review**

Inspect: `git diff 6ba818a3e..HEAD -- web apps/api-gql/internal/delivery/gql apps/api-gql/internal/delivery/http/routes/auth apps/api-gql/internal/services/channel_platforms`

Review for these regressions:

```text
1. A provider name appears in a UI conditional instead of the centralized presentation map.
2. The client constructs an OAuth callback URL or accepts a redirect URL from user input.
3. A capability-less platform loses the generic Enable bot switch or receives a synthetic unsupported operation.
4. Any disconnect path lacks the backend Manage permission gate.
5. A focused test passes only because it bypasses the component/composable action under test.
```

- [ ] **Step 5: Commit verification evidence if source changes remain**

```bash
git add web apps/api-gql/internal/delivery/gql apps/api-gql/internal/delivery/http/routes/auth apps/api-gql/internal/services/channel_platforms bun.lock docs/superpowers/plans/2026-07-23-vk-video-live-platform-task-12.md
git commit -m "test(web): verify platform bindings"
```

## Plan Self-Review

- Spec coverage: Task 1 rechecks the security correction; Task 2 supplies Bun-native tooling and generated contracts; Task 3 builds the data/OAuth layer; Task 4 provides generic rendered management; Task 5 verifies and reviews the complete surface.
- Placeholder scan: no unfinished requirements or unspecified test cases remain.
- Type consistency: all mutations accept generated `Platform`; every card action emits that same value; option and binding capabilities use `{ name: string }`; the enabled action remains binding-wide; OAuth navigation consumes `channelPlatformConnect`.
