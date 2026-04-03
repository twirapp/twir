# AGENTS.md — libs/pubsub

JavaScript pub/sub for frontend and Node services.

## OVERVIEW

TypeScript pub/sub implementation for JavaScript services. Provides message bus functionality for `apps/integrations` and frontend libraries.

## STRUCTURE

```
libs/pubsub/
├── src/
│   └── *.ts                 # Pub/sub implementation
├── package.json
├── tsconfig.json
└── ...
```

## USAGE

```typescript
import { createPubSub } from "@twir/pubsub";

const pubsub = createPubSub();

// Subscribe
pubsub.subscribe("channel", (msg) => {
	console.log(msg);
});

// Publish
pubsub.publish("channel", { data: "value" });
```

## NOTES

- Used by integrations service
- Redis-backed in production
- Type-safe messages
