# AGENTS.md — libs/cache

Caching layer with generic cacher.

## OVERVIEW

Generic caching abstraction with multiple implementations. Provides typed cache operations with automatic serialization.

## STRUCTURE

```
libs/cache/
├── generic-cacher/
│   └── db-generic-cacher.go # DB-backed cache
├── *.go                     # Cache interfaces
├── go.mod
└── ...
```

## CONVENTIONS

### Generic Cacher

```go
cacher := cache.NewGenericCacher[MyType](redisClient)

// Get or compute
val, err := cacher.Get(ctx, key, func() (MyType, error) {
    return expensiveOperation()
})
```

## ANTI-PATTERNS

- **TODO** in code: uses context.TODO() in Delete — should accept ctx

## NOTES

- Redis backend
- Automatic serialization
- TTL support
