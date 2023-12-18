module github.com/satont/twir/libs/logger

go 1.21

replace github.com/satont/twir/libs/config => ../config

require (
	github.com/getsentry/sentry-go v0.25.0
	github.com/rs/zerolog v1.30.0
	github.com/samber/slog-multi v1.0.2
	github.com/samber/slog-sentry/v2 v2.2.1
	github.com/samber/slog-zerolog/v2 v2.1.0
	github.com/satont/twir/libs/config v0.0.0-20231203205548-e635accc6b72
)

require (
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/samber/lo v1.39.0 // indirect
	github.com/samber/slog-common v0.11.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/exp v0.0.0-20231214170342-aacd6d4b4611 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
