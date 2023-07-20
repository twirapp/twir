module github.com/satont/twir/libs/migrations

go 1.20

replace github.cim/satont/twir/libs/config => ./../config

replace github.cim/satont/twir/libs/crypto => ./../crypto

require github.com/pressly/goose/v3 v3.13.4

require (
	github.cim/satont/twir/libs/config v0.0.0-00010101000000-000000000000 // indirect
	github.cim/satont/twir/libs/crypto v0.0.0-00010101000000-000000000000 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
)
