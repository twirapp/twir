module github.com/satont/twir/libs/sentry

go 1.21

replace github.com/satont/twir/libs/config => ../config

require github.com/getsentry/sentry-go v0.25.0

require (
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
)
