module github.com/satont/twir/libs/sentry

go 1.24.1

replace github.com/satont/twir/libs/config => ../config

require (
	github.com/getsentry/sentry-go v0.34.1
	github.com/satont/twir/libs/config v0.0.0-20250723210134-6e95e974f9e4
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
)
