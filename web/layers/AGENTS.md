# AGENTS.md — web/layers

Nuxt layer modules for the public website.

## OVERVIEW

Feature-isolated Nuxt layers that compose the public website. Each layer is a self-contained Nuxt application that gets merged at build time.

## STRUCTURE

```
web/layers/
├── landing/                 # Marketing site
│   ├── pages/
│   └── nuxt.config.ts
├── public/                  # Public utilities
│   ├── pages/
│   └── nuxt.config.ts
├── url-shortener/           # Link shortener UI
│   ├── pages/
│   └── nuxt.config.ts
├── pastebin/                # Code sharing UI
│   ├── pages/
│   └── nuxt.config.ts
└── overlays/                # Overlay previews
    ├── pages/
    └── nuxt.config.ts
```

## LAYERS

| Layer           | Route         | Purpose                  |
| --------------- | ------------- | ------------------------ |
| `landing`       | `/`           | Marketing site           |
| `url-shortener` | `/s/*`        | Short link redirect page |
| `pastebin`      | `/paste/*`    | Code sharing UI          |
| `public`        | `/public/*`   | Public utilities         |
| `overlays`      | `/overlays/*` | Overlay previews         |

## CONVENTIONS

Each layer:

- Has its own `nuxt.config.ts`
- Can define pages, components, composables
- Gets merged by parent `web/nuxt.config.ts`

## USAGE

Layers registered in `web/nuxt.config.ts`:

```typescript
export default defineNuxtConfig({
	extends: [
		"./layers/landing",
		"./layers/url-shortener",
		// ...
	],
});
```

## NOTES

- Layers don't run independently
- Parent config takes precedence
- Shared components go in `web/app/components`
