# AGENTS.md — libs/gomodels

Legacy Go models (deprecated, migrating to entities).

## OVERVIEW

**DEPRECATED**: Legacy database models. Being migrated to `libs/entities`. New code should use entity pattern.

## STRUCTURE

```
libs/gomodels/
├── *.go                     # Model definitions
├── go.mod
└── ...
```

## STATUS

- **DO NOT** add new models here
- **MIGRATE** existing models to `libs/entities`
- Some models still in use during migration

## MIGRATION PATH

```go
// OLD (gomodels)
type Channel struct { ... }

// NEW (entities)
package channels

type Channel struct {
    // fields
    isNil bool
}

func (c Channel) IsNil() bool { return c.isNil }
var Nil = Channel{isNil: true}
```

## NOTES

- Migration in progress
- Coordinate with team before modifying
