# AGENTS.md — apps/executron

Sandboxed JavaScript execution service (Bun/TypeScript, NOT Go).

## OVERVIEW

Executes user-provided JavaScript (custom scripts/commands) in an `isolated-vm` sandbox. Receives jobs via NATS (`twirBus.Executron.Execute`), fetches per-channel secrets from DB, returns `{ result, error }`.

## STRUCTURE

```
apps/executron/
├── src/
│   ├── index.ts            # Entry: NATS subscriber group 'executron'
│   └── libs/
│       ├── executor.ts     # Sandbox core (~13K, isolated-vm)
│       ├── host-fetch.ts   # fetch exposed to sandbox
│       ├── host-storage.ts # key-value storage exposed to sandbox
│       ├── host-timers.ts  # timers exposed to sandbox
│       ├── url-validation.ts # SSRF guard for host-fetch
│       ├── db.ts           # channel secrets lookup
│       └── twirbus.ts      # NATS bus instance
├── build.ts                # bun build script
└── Dockerfile
```

## CONVENTIONS

- Only `language === 'javascript'` is supported; other languages return an error payload.
- Anything reachable from user code goes through a `host-*.ts` module — add new capabilities there, never expose Node/Bun APIs directly.
- External requests from scripts must pass `url-validation.ts` (SSRF protection).

## COMMANDS

```bash
bun run dev     # watch mode (loads root .env)
bun run build   # tsc --noEmit + bun build.ts
bun test        # bun:test (src/libs/executor.test.ts); no package.json test script
```

## NOTES

- Sandbox crashes must not kill the service: root handlers on `uncaughtException`/`unhandledRejection` only log — keep it that way.
