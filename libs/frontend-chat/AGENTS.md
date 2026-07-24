# AGENTS.md — libs/frontend-chat

Vue 3 chat widget component library (`@twir/frontend-chat`).

## OVERVIEW

Chat UI components + composables (message rendering, fragments-to-chunks, overlay socket) consumed by the web dashboard chat overlay settings and `frontend/overlays` chat overlay.

## CONVENTIONS

- **Shipped as SFC source, not prebuilt**: `exports` points at `./src/index.ts`. Consumers MUST transpile it — `web/nuxt.config.ts` has `build.transpile: ['@twir/frontend-chat']`. When adding a new consumer, add the same transpile entry.
- Peer dependency on `vue` (catalog version) — do not bundle Vue.
- Uses `@twir/api` (openapi client) and `@twir/fontsource` workspaces.

## STRUCTURE

```
libs/frontend-chat/src/
├── index.ts        # Public exports
├── components/     # Chat message components (scoped styles inside SFCs)
├── composables/    # Socket/fragments logic
├── styles/
├── helpers.ts
└── types.ts
```

## NOTES

- Because SFCs ship raw, scoped message styles travel with the source — do not "optimize" by extracting/building them away without updating all consumers.
