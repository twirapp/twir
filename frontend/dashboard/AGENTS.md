# AGENTS.md — frontend/dashboard

Vue 3 dashboard application for streamers. Main administrative interface for Twir.

## OVERVIEW

SPA built with Vue 3 + Vite. Provides streamers with tools for chat moderation, alerts, overlays, integrations, and channel management. Communicates with backend via GraphQL (urql).

## STRUCTURE

```
frontend/dashboard/
├── src/
│   ├── main.ts              # Entry point
│   ├── api/                 # GraphQL queries/mutations
│   │   ├── integrations/    # Unified integrations query
│   │   └── ...
│   ├── components/          # Shared UI components
│   │   └── ui/              # shadcn-vue components
│   ├── features/            # Feature-based modules
│   │   ├── admin-panel/     # Admin functionality
│   │   ├── chat-alerts/     # Alert configuration
│   │   ├── overlays/        # Overlay settings
│   │   └── ...
│   ├── pages/               # Route components
│   └── layouts/             # Layout components
├── vite.config.ts           # Vite configuration
├── package.json
└── Dockerfile
```

## ENTRY POINTS

| Type        | Path             | Purpose                  |
| ----------- | ---------------- | ------------------------ |
| Main        | `src/main.ts`    | Vue app bootstrap        |
| Router      | `src/router.ts`  | Route definitions        |
| Vite Config | `vite.config.ts` | Build config, dev server |

## KEY COMMANDS

```bash
# Development (runs on :3006)
bun dev                    # vite dev server

# GraphQL codegen (auto-watches in dev)
bun run codegen           # Regenerate GraphQL types

# Build for production
bun run build             # codegen + vue-tsc + vite build
```

## CONVENTIONS

### Components

```vue
<script setup lang="ts">
import { computed } from "vue";
import { User } from "lucide-vue-next";

interface Props {
	title: string;
}

const props = defineProps<Props>();
const isActive = defineModel<boolean>("modelValue");

const emit = defineEmits<{
	(e: "update", value: string): void;
}>();
</script>
```

- **ALWAYS** use `<script setup lang="ts">`
- **ALWAYS** include `.vue` extension in imports
- Use `defineModel` for v-model bindings
- Use `lucide-vue-next` for icons

### Forms (vee-validate + zod)

```vue
<script setup lang="ts">
import { useForm } from "vee-validate";
import { toTypedSchema } from "@vee-validate/zod";
import { z } from "zod";

const formSchema = z.object({
	name: z.string().min(2),
	enabled: z.boolean(),
});

const { handleSubmit } = useForm({
	validationSchema: toTypedSchema(formSchema),
});

const onSubmit = handleSubmit((values) => {
	// API call
});
</script>

<template>
	<form @submit="onSubmit">
		<!-- Use v-slot="{ componentField }" for inputs -->
		<!-- Use v-slot="{ value, handleChange }" for switches -->
	</form>
</template>
```

- Use `useForm` hook (NOT `<Form>` component)
- Bind to native `<form @submit>`
- Always include `<FormMessage />` for errors

### Integrations Page Pattern

All integrations data fetched via unified query:

```typescript
// src/api/integrations/integrations-page.ts
const { useIntegrationsPageData } = from './integrations-page';

const pageData = useIntegrationsPageData();
// Access: pageData.discordData, pageData.spotifyData, etc.
```

## ANTI-PATTERNS

- **DO NOT** use `emit('update:modelValue')` — use `defineModel()`
- **DO NOT** use vee-validate's `<Form>` component — use `useForm` hook
- **DO NOT** write custom CSS in `<style>` — use Tailwind utilities
- **DO NOT** create new delete confirmation dialogs — use existing component

## STYLING

- Tailwind CSS utility classes only
- Use theme colors: `bg-primary`, `text-accent`, `border-destructive`
- No arbitrary values: prefer `p-4` over `p-[15px]`

## PORTS

| Service       | Port |
| ------------- | ---- |
| Dev Server    | 3006 |
| HMR WebSocket | 3006 |

## NOTES

- GraphQL endpoint: `http://localhost:3009/query` (api-gql)
- Auto-generates GraphQL types on schema changes
- Uses shadcn-vue for UI components
- Supports i18n via vue-i18n
