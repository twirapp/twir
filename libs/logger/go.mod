module github.com/satont/twir/libs/logger

go 1.21

replace github.com/satont/twir/libs/config => ../config

require (
	github.com/getsentry/sentry-go v0.26.0
	github.com/rs/zerolog v1.31.0
	github.com/samber/slog-multi v1.0.2
	github.com/samber/slog-sentry/v2 v2.4.0
	github.com/samber/slog-zerolog/v2 v2.2.0
	github.com/satont/twir/libs/config v0.0.0-20231203205548-e635accc6b72
	go.uber.org/fx v1.20.1
	go.uber.org/zap v1.26.0
)

require (
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/samber/lo v1.39.0 // indirect
	github.com/samber/slog-common v0.15.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.17.1 // indirect
	go.uber.org/goleak v1.3.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20240119083558-1b970713d09a // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
