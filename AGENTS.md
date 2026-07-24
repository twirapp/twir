### **Project Development Guidelines for AI Assistants (GitHub Copilot)**

> **Navigation**: This is the root AGENTS.md. For module-specific guidelines, see:
>
> **Apps:**
>
> - [`apps/api-gql/AGENTS.md`](apps/api-gql/AGENTS.md) — GraphQL API service
> - [`apps/bots/AGENTS.md`](apps/bots/AGENTS.md) — Twitch bot service
> - [`apps/emotes-cacher/AGENTS.md`](apps/emotes-cacher/AGENTS.md) — Emote caching
> - [`apps/events/AGENTS.md`](apps/events/AGENTS.md) — Event workflows
> - [`apps/eventsub/AGENTS.md`](apps/eventsub/AGENTS.md) — Twitch EventSub
> - [`apps/integrations/AGENTS.md`](apps/integrations/AGENTS.md) — 3rd party integrations
> - [`apps/parser/AGENTS.md`](apps/parser/AGENTS.md) — Command parser
> - [`apps/scheduler/AGENTS.md`](apps/scheduler/AGENTS.md) — Task scheduler
> - [`apps/timers/AGENTS.md`](apps/timers/AGENTS.md) — Chat timers
> - [`apps/tokens/AGENTS.md`](apps/tokens/AGENTS.md) — OAuth tokens
> - [`apps/websockets/AGENTS.md`](apps/websockets/AGENTS.md) — WebSocket server
> - [`apps/executron/AGENTS.md`](apps/executron/AGENTS.md) — Sandboxed JS execution (Bun/TS)
> - [`apps/twitch-mock/AGENTS.md`](apps/twitch-mock/AGENTS.md) — Local Twitch API mock (dev only)
>
> **Frontend:**
>
> - [`web/AGENTS.md`](web/AGENTS.md) — Nuxt website + dashboard (layers)
> - [`web/layers/dashboard/AGENTS.md`](web/layers/dashboard/AGENTS.md) — Dashboard SPA (Nuxt layer)
> - [`frontend/overlays/AGENTS.md`](frontend/overlays/AGENTS.md) — Browser overlays
>
> **Tools & Libraries:**
>
> - [`cli/AGENTS.md`](cli/AGENTS.md) — Custom CLI tool
> - [`libs/repositories/AGENTS.md`](libs/repositories/AGENTS.md) — Data access layer
> - [`libs/entities/AGENTS.md`](libs/entities/AGENTS.md) — Domain entities
> - [`libs/grpc/AGENTS.md`](libs/grpc/AGENTS.md) — gRPC definitions
> - [`libs/migrations/AGENTS.md`](libs/migrations/AGENTS.md) — Database migrations
> - [`libs/integrations/AGENTS.md`](libs/integrations/AGENTS.md) — External APIs
> - [`libs/twitch/AGENTS.md`](libs/twitch/AGENTS.md) — Twitch client
> - [`libs/gomodels/AGENTS.md`](libs/gomodels/AGENTS.md) — Legacy models
> - [`libs/cache/AGENTS.md`](libs/cache/AGENTS.md) — Caching layer
> - [`libs/bus-core/AGENTS.md`](libs/bus-core/AGENTS.md) — Message bus
> - [`libs/config/AGENTS.md`](libs/config/AGENTS.md) — Configuration
> - [`libs/i18n/AGENTS.md`](libs/i18n/AGENTS.md) — Internationalization
> - [`libs/logger/AGENTS.md`](libs/logger/AGENTS.md) — Logging
> - [`libs/types/AGENTS.md`](libs/types/AGENTS.md) — TypeScript types
> - [`libs/api/AGENTS.md`](libs/api/AGENTS.md) — API client
> - [`libs/pubsub/AGENTS.md`](libs/pubsub/AGENTS.md) — Pub/sub messaging
> - [`libs/frontend-chat/AGENTS.md`](libs/frontend-chat/AGENTS.md) — Chat widget Vue library

This document outlines the core conventions, technologies, and patterns used in this project. Please
adhere to these guidelines strictly to maintain code consistency and quality.

### **1. General Project Context**

- **Structure:** This is a monorepo.
  - **`web`**: Nuxt 3 app — public website **and** dashboard (as Nuxt layers, see
    `web/layers/dashboard`). The old `frontend/dashboard` (Vue 3 + Vite) was removed.
  - **`frontend/overlays`**: Browser overlays (Vue 3 + Vite).
  - **`apps/api-gql`**: The main backend service (Go) serving GraphQL and HTTP APIs.
  - **`libs`**: Shared Go libraries and TS packages.
- **Package Manager & Runtime:** We use **Bun** for all JavaScript/TypeScript package management,
  script execution, and as the runtime. Use `bun install`, `bun add`, and `bun run` commands.
- **Primary Technologies:**
  - **Web (site + dashboard):** Nuxt 3, TypeScript, Tailwind CSS, Reka UI / shadcn-vue, Pinia,
    Urql, vee-validate, zod, `@nuxt/icon`.
  - **Backend:** Go (Golang), pgx (PostgreSQL driver), gqlgen (GraphQL).
  - **Tooling:** Bun, Docker, oxlint (`bun lint`).
- **MCP Tools Usage**
  - **CodeGraph MCP** — use for exploring project structure, understanding architecture, finding symbol definitions/callers/callees, and analyzing impact of changes. Always prefer codegraph over manual file searching when codegraph is configured.
  - **Context7 MCP** — use for reading up-to-date library/framework documentation. If you are unsure about some library (e.g., you know version 3 but the project uses version 4), ask context7 for the correct version docs. Avoid searching the internet for library documentation — always use context7 if available.
  - **Postgres MCP** — use for inspecting the database schema, running read-only queries to understand data structure, and verifying assumptions about tables/columns before writing code.
  - **grepApp MCP** — use for searching real-world code examples from public GitHub repositories when you need to see how an API or pattern is used in practice.
  - Always inform the user when you are using MCP tools. Let me know which MCP you are using and why.
- **Skills**
  - When working with postgres read releated skill in skills directory
  - When working on the **Nuxt public website (`web/`)**, load the **`nuxt`** skill for Nuxt 3
    framework guidance (routing, data fetching, server routes, SSR/hydration, etc.).
  - When building UI with **Reka UI** primitives, load the **`reka-ui`** skill for component APIs,
    accessibility patterns, composition (`asChild`), controlled state, and SSR guidance.
  - When working with **shadcn-vue** components, load the **`shadcn-vue`** skill for component
    usage, styling conventions, CLI commands, and theming.
  - Skill files are located under `.agents/skills/` and can be loaded with the `skill` tool.

---

### **2. Vue.js Frontend Development (Dashboard & Web)**

#### **2.1. Component Structure & Syntax**

- **Composition API:** Always use the Composition API.
- **Script Setup:** All Single File Components (SFCs) **must** use the `<script setup lang="ts">`
  syntax. Do not use the `setup()` function within the `export default` block.
- **Type Definitions:** Use `defineProps`, `defineEmits`, and `defineSlots` with explicit TypeScript
  types for clear, type-safe component interfaces.
- **Model instead of emit update** Use new vue syntax for shorter code, for example instaed of `emit('open:update', value)` use `defineModel`

**Example:**

```vue
<script setup lang="ts">
import { computed } from "vue";

// Use interface or type for props definition
interface Props {
	title: string;
	items: string[];
}

const props = withDefaults(defineProps<Props>(), {});

const isActive = defineModel<boolean>({ default: false });

const emit = defineEmits<{
	(e: "itemSelected", item: string): void;
	(e: "closed"): void;
}>();

const handleItemClick = (item: string) => {
	emit("itemSelected", item);
	isActive.value = false;
};

const titleDisplay = computed(() => props.title.toUpperCase());
</script>

<template>
	<!-- Component template here -->
</template>
```

#### **2.2. Component Imports**

- **File Extension:** Always include the `.vue` file extension when importing Vue components. This
  improves clarity and avoids potential bundler configuration issues.

**Correct:**

```typescript
import UserProfile from "./components/UserProfile.vue";
import AppHeader from "@/components/layout/AppHeader.vue";
```

**Incorrect:**

```typescript
import UserProfile from "./components/UserProfile";
```

#### **2.3. Existing Components**

- **Confirmations:** The project already has a standardized confirmation component. **Do not create
  a new one.** Use `ActionConfirm` from `web/app/components/ui/action-confirm/ActionConfirm.vue`
  (`v-model:open`, `confirm`/`cancel` emits) for delete/destructive confirmations.
- **Responsive modals:** Use `web/layers/dashboard/components/dialog-or-sheet.vue` (`DialogOrSheet`)
  for dialog-on-desktop / sheet-on-mobile UX.

---

### **3. Forms with `vee-validate` and `zod`**

We use `vee-validate` with `zod` for schema-based validation. Follow this pattern precisely.

- **Schema Library:** Use **zod** to define validation schemas.
- **Adapter:** Use the `@vee-validate/zod` library to connect `zod` schemas to `vee-validate`.
- **Hook, not Component:** Use the `useForm` hook from `vee-validate`.
- **Native `<form>` Element:** Bind your submission logic to a native HTML `<form>` element's
  `@submit` event. **DO NOT use the `<Form>` component provided by `vee-validate`**.
- **Fields:**
  - For standard text inputs, textareas, etc., use `v-slot="{ componentField }"` on your field
    wrapper and bind `v-bind="componentField"` to the input element.
  - For switches, checkboxes, and custom toggle components, use `v-slot="{ value, handleChange }"`
    to manage state.
- **Error Messages:** Always include the `FormMessage` component immediately after a form field to
  display validation errors.

**Example Form Structure:**

```vue
<script setup lang="ts">
import { useForm } from "vee-validate";
import { z } from "zod";
import { FormMessage, FormItem, FormField, FormControl, FormLabel } from "@/components/ui/form";
import Input from "@/components/ui/input/Input.vue"; // Example custom input
import Switch from "@/components/ui/switch/Switch.vue";

const formSchema = z.object({
	name: z.string().min(2, "Name must be at least 2 characters."),
	enabled: z.boolean(),
});

const { handleSubmit, defineField } = useForm({
	validationSchema: formSchema,
});

const onSubmit = handleSubmit((values) => {
	console.log("Form submitted:", values);
	// API call logic here
});
</script>

<template>
	<form @submit="onSubmit" class="space-y-4">
		<div>
			<FormField v-slot="{ value, handleChange }" name="enabled">
				<FormItem class="flex gap-2 space-y-0 items-center">
					<FormLabel>...</FormLabel>
					<FormControl>
						<Switch :model-value="value" @update:model-value="handleChange" />
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>
		</div>

		<div>
			<FormField v-slot="{ componentField }" name="name">
				<FormItem>
					<FormLabel>...</FormLabel>
					<FormControl>
						<Input v-bind="componentField" type="number" :max="86400" />
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>
		</div>

		<button type="submit">Submit</button>
	</form>
</template>
```

---

### **4. Iconography**

- **Primary Library:** **Lucide icons via the Nuxt `<Icon />` component** (`name="lucide:icon-name"`).
  This is the standard for UI chrome everywhere in `web` (site and dashboard). **Do not** import
  `lucide-vue-next` components in `web`.
- **Brand / platform logos (Twitch, Kick, VK, etc.):** use `simple-icons:*` via the same `<Icon />`
  component — the collection is installed locally. Color them with Tailwind `text-[#...]` classes.
- **Integration cards & artwork:** local SVG collections via `<Icon />` (`twir-integrations:*`,
  `twir-overlays:*`, `twir-compare:*`) or `Svgo*` components from `nuxt-svgo` where already used.
- **Fallback:** Only if a specific icon is not available in Lucide/simple-icons should you add a
  local SVG file. This should be a rare exception.

**Example:**

```vue
<template>
	<Icon name="lucide:user" class="h-4 w-4 mr-2" />
	<Icon name="simple-icons:twitch" class="h-4 w-4 text-[#9146FF]" />
</template>
```

---

### **5. Styling with Tailwind CSS**

- **Utility-First:** All styling must be done using Tailwind CSS utility classes directly in the
  `<template>` block. Avoid writing custom CSS in `<style>` blocks unless absolutely necessary for a
  complex, non-reusable scenario.
- **Project Configuration:** Adhere strictly to the project's `tailwind.config.js`.
  - **Colors:** Use the defined theme colors (e.g., `bg-primary`, `text-accent`,
    `border-destructive`). Do not use arbitrary hex codes or default Tailwind colors if custom ones
    are defined.
  - **Spacing & Sizing:** Use the defined spacing scale (e.g., `p-4`, `m-8`, `w-32`) instead of
    arbitrary values like `p-[15px]`.
  - **Component Classes:** If we use a library like `shadcn-vue` or have our own `@apply` directives
    for component base styles, be aware of and use them.

---

### **6. Dashboard Integrations Page Architecture**

The integrations page uses a **unified GraphQL query** pattern to fetch all integration data in a
single request, optimizing network usage and improving user experience.

#### **6.1. Unified Query Pattern**

- **Single Query File:** All integrations page data is fetched via a unified query in
  `web/layers/dashboard/api/integrations/integrations-page.ts`.
- **Why:** This approach allows fetching data for all integrations (Discord, Spotify, LastFM,
  Valorant, DonationAlerts, etc.) in a single GraphQL request, which is significantly more efficient
  than making separate requests per integration.
- **Composable:** Use `useIntegrationsPageData()` to access the unified data. It provides computed
  refs for each integration's data.

#### **6.2. Adding or Refactoring Integrations**

When creating a new integration or refactoring an existing one to use GraphQL:

1.  **Add fields to the unified query** in `integrations-page.ts`:

        ````typescript
         const IntegrationsPageQuery = graphql(`
        query IntegrationsPageData { # ... existing fields ...

        													 # New integration
        													 myNewIntegrationData {
        														 enabled
        														 userName
        														 avatar
        													 }
        													 myNewIntegrationAuthLink
        												 }
        											 `)
        											 ```

        ````

2.  **Add computed refs** for the new integration data:
    `typescript
// MyNewIntegration
const myNewIntegrationData = computed(() => query.data.value?.myNewIntegrationData ?? null)
const myNewIntegrationAuthLink = computed(() => query.data.value?.myNewIntegrationAuthLink ?? null)
`

3.  **Export the new computed refs** in the return statement.

4.  **Use the unified data in components** instead of creating separate queries:
    `typescript
const integrationsPage = useIntegrationsPageData()
// Access via integrationsPage.myNewIntegrationData
`

#### **6.3. Mutations**

- Mutations (login, logout, update, etc.) should still be defined separately in dedicated files or
  in `integrations.ts`.
- Use `integrationsPageCacheKey` to invalidate the unified query cache after mutations:
  ```typescript
  const myMutation = () =>
  	useMutation(graphql(`mutation MyMutation { ... }`), [integrationsPageCacheKey]);
  ```

---

### **7. Go (Golang) Backend**

- **Migrations**
  - If need to create new database migration for your task, use:
  - command `bun cli m create --name value --db postgres|clickhouse --type sql|go`

- **Code Style:** Follow standard Go formatting (`gofmt`/`goimports`).
- **Project Structure:**
  - **`apps/api-gql`**: Main API service.
    - `internal/delivery/gql`: GraphQL resolvers.
    - `internal/delivery/http`: HTTP handlers.
    - `internal/services`: Business logic layer.
  - **`libs/repositories`**: Data access layer.
- **Entities**
  - Write entities in `libs/entities/{entity_name}/entity.go` file.
  - Entities should contain only domain logic and validation.
    - Avoid dependencies on other layers (e.g., repositories, services).
  - use Nil thing
- **Repositories:**
  - Always use **pgx** implementations.
  - Located in `libs/repositories/{repository_name}/pgx/pgx.go`.
  - **NEVER** use GORM or other ORMs.
  - Repository should return entity written in `libs/entities/{entity_name}/entity.go` file.
    - For new models created in repositories, or when editing some repository, you should
      create/update model inside repository, do not create separate file for model. to include
      `isNil` property, and `IsNil` method to check if the model is
      empty, also create `var Nil = &Model{}` to represent an empty model. Example:

```go
type SomeModel struct {
	ID            string
	ChannelID     string

	isNil bool
}

func (c SomeModel) IsNil() bool {
	return c.isNil
}

var Nil = SomeModel{
	isNil: true,
}

```

- **Mappers:**
  - When creating new services (e.g., in `api-gql`), always create an entity mapper.
  - Data flow: `Model (DB)` -> `Entity (Domain)` -> `DTO (GraphQL/HTTP)`.
- **GraphQL Generation:**
  - After updating GraphQL schemas (`.graphql` files), run `bun cli build gql` to regenerate
    resolvers.
  - After regeneration, refresh your data (re-read Golang files) to pick up changes.
- **Error Handling:**
  - Use `fmt.Errorf` with `%w` for wrapping errors.
  - Create custom error types if needed for specific domain error handling.
- **Loging**
  - Use \*Context where we have `ctx` field, so for example InfoContext(ctx, ...)
  - For errors use `logger.Error(err)`, which is available under logger lib. It's a shortcut.
  - Use `logger.String` and other methods, instead of inline arguments.

---

### **8. Build & CI**

#### **8.1. Local Development**

```bash
# Start infrastructure (Postgres, Redis, etc.)
docker compose -f docker-compose.dev.yml up -d

# Start all services in dev mode
bun dev

# Run custom CLI commands
bun cli <command>
```

#### **8.2. Build Commands**

```bash
# Build everything
bun cli build

# Build specific app
bun cli b app <app-name>

# Regenerate GraphQL resolvers
bun cli build gql

# Run linting
bun lint
```

#### **8.3. CI/CD (GitHub Actions)**

- **Primary Workflow**: `.github/workflows/dockerv3.yml`
  - Trigger: tags `v*` + manual dispatch
  - Matrix builds for all apps
  - Change detection (only builds changed apps)
  - Builds Docker images → `registry.twir.app/twirapp/<app>:latest`
- **Branch Checks**: `.github/workflows/build-and-lint.yml`
  - Trigger: push to any branch except `main` + manual dispatch
  - Runs `bun cli deps`, `bun cli build` and `bun lint`

#### **8.4. Key Files**

| File                       | Purpose                            |
| -------------------------- | ---------------------------------- |
| `package.json`             | Root workspace config, Bun scripts |
| `.bun-version`             | Pinned Bun version                 |
| `go.work`                  | Go workspace definition            |
| `docker-compose.dev.yml`   | Local infrastructure               |
| `docker-compose.stack.yml` | Production stack (Swarm)           |
