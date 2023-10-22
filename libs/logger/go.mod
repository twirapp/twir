module github.com/satont/twir/libs/logger

go 1.21.0

replace github.com/satont/twir/libs/config => ../config

require github.com/getsentry/sentry-go v0.25.0

require (
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
)
