module github.com/twirapp/twir/libs/bus-core

go 1.24.1

replace (
	github.com/satont/twir/libs/config => ../config
)

require (
	github.com/google/uuid v1.6.0
	github.com/nats-io/nats.go v1.37.0
	github.com/satont/twir/libs/config v0.0.0-00010101000000-000000000000
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)
