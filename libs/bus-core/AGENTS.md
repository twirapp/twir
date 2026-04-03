# AGENTS.md — libs/bus-core

Message bus abstractions for inter-service communication.

## OVERVIEW

Core message bus interfaces and implementations. Provides pub/sub patterns for async communication between services.

## STRUCTURE

```
libs/bus-core/
├── *.go                     # Bus interfaces and implementations
├── go.mod
└── ...
```

## CONVENTIONS

### Bus Interface

```go
// Publisher/Subscriber pattern
type Publisher interface {
    Publish(ctx context.Context, topic string, msg interface{}) error
}

type Subscriber interface {
    Subscribe(ctx context.Context, topic string, handler Handler) error
}
```

## DEPENDENCIES

- Redis (pub/sub backend)
- libs/pubsub (JS side)

## NOTES

- Used for cross-service events
- Supports multiple backends
- Go services use this directly
