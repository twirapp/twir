module github.com/satont/twir/libs/sentry

go 1.21

replace github.com/satont/twir/libs/config => ../config

require github.com/getsentry/sentry-go v0.25.0

require (
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
)
