# AGENTS.md — web/layers

Nuxt layer modules composing the whole web app (public website + dashboard).

## OVERVIEW

Feature-isolated Nuxt layers. Each layer is a self-contained Nuxt application that gets merged at build time.

## STRUCTURE

```
web/layers/
├── dashboard/               # Dashboard SPA (own AGENTS.md)
├── landing/                 # Marketing site
├── public/                  # Public utilities
├── url-shortener/           # Link shortener UI
├── pastebin/                # Code sharing UI
├── overlays/                # Overlay previews
└── widgets/                 # Public widgets
```

## LAYERS

| Layer           | Route          | Purpose                  |
| --------------- | -------------- | ------------------------ |
| `dashboard`     | `/dashboard/*` | Dashboard SPA            |
| `landing`       | `/`            | Marketing site           |
| `url-shortener` | `/s/*`         | Short link redirect page |
| `pastebin`      | `/paste/*`     | Code sharing UI          |
| `public`        | `/public/*`    | Public utilities         |
| `overlays`      | `/overlays/*`  | Overlay previews         |
| `widgets`       | `/w/*`         | Public widgets           |

## CONVENTIONS

Each layer:

- Has its own `nuxt.config.ts`
- Can define pages, components, composables
- Gets merged by parent `web/nuxt.config.ts`

## USAGE

Layers under `layers/` are auto-registered by Nuxt — no manual `extends` entries needed.

## NOTES

- Layers don't run independently
- Parent config takes precedence
- Shared components go in `web/app/components`
