# AGENTS.md — libs/i18n

Internationalization utilities.

## OVERVIEW

i18n helpers for Go backend. Provides translation keys and formatting for multilingual support.

## STRUCTURE

```
libs/i18n/
├── *.go                     # i18n utilities
├── go.mod
└── ...
```

## USAGE

```go
import "libs/i18n"

msg := i18n.T(ctx, "key")
```

## NOTES

- Go backend only
- Frontend uses vue-i18n
- Locale files in apps/parser/locales
