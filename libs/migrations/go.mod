module github.com/satont/twir/libs/migrations

go 1.24.1

replace github.cim/satont/twir/libs/config => ./../config

replace github.cim/satont/twir/libs/crypto => ./../crypto

replace github.com/twirapp/twir/libs/grpc => ./../grpc

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/pressly/goose/v3 v3.24.3
	github.com/satont/twir/libs/config v0.0.0-20250723210134-6e95e974f9e4
	github.com/satont/twir/libs/crypto v0.0.0-20250723210134-6e95e974f9e4

)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20250606033433-dcc06ee1d476 // indirect
	golang.org/x/sync v0.16.0 // indirect
)
