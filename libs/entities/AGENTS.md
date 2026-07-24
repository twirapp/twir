# AGENTS.md — libs/entities

Domain entities for Go backend.

## OVERVIEW

Pure domain models used across backend services. Contains business logic and validation rules. No dependencies on repositories or services.

## STRUCTURE

```
libs/entities/
├── {domain}/                # One package per domain
│   └── entity.go           # Entity definition
├── go.mod
└── ...

# Current domains include:
channel/ channel_platform/ platform/ kick_bot/ vk_integration/ dashboard_widget/
faceit_integration/ streamlabs_integration/ timers/ plan/ secret/ voteban/ obs/
scheduled_vips/ short_links_custom_domain/ song_request_overlay_settings/
webhook_notifications/ channels_giveaways*/ command_with_relations/ custom_overlay/ ...
```

## CONVENTIONS

### Entity Definition

```go
package channels

type Channel struct {
    ID        string
    Name      string
    IsEnabled bool

    isNil     bool
}

func (c Channel) IsNil() bool {
    return c.isNil
}

var Nil = Channel{isNil: true}
```

### Rules

- Entities contain **only domain logic**
- **NO** repository dependencies
- **NO** service dependencies
- Include validation methods
- Use Nil pattern for empty results

## ANTI-PATTERNS

- **NEVER** import repositories in entities
- **NEVER** import services in entities
- **NEVER** use database types in entities

## USAGE

```go
// In services
import "libs/entities/channels"

func (s *Service) GetChannel(ctx context.Context, id string) (channels.Channel, error) {
    // Repository returns entity
    return s.repo.GetByID(ctx, id)
}
```

## NOTES

- Data flow: Model → Entity → DTO
- Entities are shared across all Go apps
- Keep entity logic pure
