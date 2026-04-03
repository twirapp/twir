# AGENTS.md — frontend/overlays

Browser overlay renderer for Twitch streams.

## OVERVIEW

Vue 3 application that runs in browser OBS overlays. Displays alerts, chat dudes (avatars), and interactive elements. Communicates via WebSocket with the websockets service.

## STRUCTURE

```
frontend/overlays/
├── src/
│   ├── main.ts              # Entry point
│   ├── app.vue              # Root component
│   ├── composables/         # Vue composables
│   │   └── dudes/           # Dudes (avatar) system
│   │       ├── use-dudes.ts
│   │       └── use-dudes-socket.ts
│   └── ...
├── vite.config.ts
├── package.json
└── Dockerfile
```

## ENTRY POINTS

| Type        | Path             | Purpose             |
| ----------- | ---------------- | ------------------- |
| Main        | `src/main.ts`    | Overlay bootstrap   |
| Vite Config | `vite.config.ts` | Build configuration |

## KEY COMMANDS

```bash
# Development (runs on :3008)
bun dev                    # vite dev

# Build
bun run build             # vite build
```

## CONVENTIONS

Same as dashboard:

- `<script setup lang="ts">`
- `lucide-vue-next` for icons
- Tailwind CSS for styling
- urql for data fetching

## DUDES SYSTEM

Interactive avatar system:

- `use-dudes.ts` — Dudes rendering and animation
- `use-dudes-socket.ts` — WebSocket communication

## ANTI-PATTERNS

- **TODO**: refactor to new api proxy (in use-dudes.ts)
- **TODO**: rename and deprecate `eyes_color`, `cosmetics_color`

## DEPENDENCIES

- websockets service (Socket.io)
- api-gql (initial data load)

## PORTS

| Service    | Port |
| ---------- | ---- |
| Dev Server | 3008 |

## NOTES

- Loaded in OBS browser source
- Real-time via Socket.io
- Dudes use canvas rendering
