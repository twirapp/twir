module github.com/satont/twir/libs/migrations

go 1.21

replace github.cim/satont/twir/libs/config => ./../config

replace github.cim/satont/twir/libs/crypto => ./../crypto

require (
	github.com/google/uuid v1.3.0
	github.com/lib/pq v1.10.9
	github.com/pressly/goose/v3 v3.13.4
	github.com/satont/twir/libs/config v0.0.0-20230713153539-b2fe2b3b5757
	github.com/satont/twir/libs/crypto v0.0.0-20230713153539-b2fe2b3b5757
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/tools v0.12.0 // indirect
)
