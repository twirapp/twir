module github.com/satont/twir/libs/sentry

go 1.21

replace github.com/satont/twir/libs/config => ../config

require (
	github.com/getsentry/sentry-go v0.25.0
	github.com/satont/twir/libs/config v0.0.0-00010101000000-000000000000
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
)
