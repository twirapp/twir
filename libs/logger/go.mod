module github.com/satont/twir/libs/logger

go 1.23.0

replace (
	github.com/satont/twir/libs/config => ../config
	github.com/satont/twir/libs/gomodels => ../gomodels
	github.com/satont/twir/libs/pubsub/audit-logs => ../pubsub
	github.com/twirapp/twir/libs/bus-core => ../bus-core
)

require (
	github.com/getsentry/sentry-go v0.26.0
	github.com/goccy/go-json v0.10.2
	github.com/google/uuid v1.6.0
	github.com/guregu/null v4.0.0+incompatible
	github.com/rs/zerolog v1.31.0
	github.com/samber/slog-common v0.17.1
	github.com/samber/slog-multi v1.0.2
	github.com/samber/slog-sentry/v2 v2.4.0
	github.com/samber/slog-zerolog/v2 v2.2.0
	github.com/satont/twir/libs/gomodels v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/pubsub/audit-logs v0.0.0-00010101000000-000000000000
	github.com/twirapp/twir/libs/bus-core v0.0.0-20240225024146-742838c78cea
	go.uber.org/fx v1.21.0
	go.uber.org/zap v1.27.0
	gorm.io/gorm v1.25.12
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/redis/go-redis/v9 v9.5.3 // indirect
	github.com/samber/lo v1.44.0 // indirect
	github.com/satont/twir/libs/types v0.0.0-20241026093228-a5284254dcfe // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	go.uber.org/dig v1.17.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.16.0 // indirect
)
