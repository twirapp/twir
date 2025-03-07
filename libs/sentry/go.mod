module github.com/satont/twir/libs/sentry

go 1.24.1

replace github.com/satont/twir/libs/config => ../config

require (
	github.com/getsentry/sentry-go v0.29.1
	github.com/satont/twir/libs/config v0.0.0-20231203205548-e635accc6b72
)

require (
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)
