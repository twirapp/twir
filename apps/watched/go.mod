module github.com/satont/tsuwari/apps/watched

go 1.19

replace github.com/satont/tsuwari/libs/config => ../../libs/config

replace github.com/satont/tsuwari/libs/nats => ../../libs/nats

replace github.com/satont/tsuwari/libs/gomodels => ../../libs/gomodels

replace github.com/satont/tsuwari/libs/twitch => ../../libs/twitch

require (
	github.com/golang-jwt/jwt/v4 v4.4.2 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/nats-io/nats.go v1.19.0 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/satont/go-helix/v2 v2.7.22 // indirect
	github.com/satont/tsuwari/libs/config v0.0.0-20221125194658-5cb70dbdbf2a // indirect
	github.com/satont/tsuwari/libs/gomodels v0.0.0-20221125194658-5cb70dbdbf2a // indirect
	github.com/satont/tsuwari/libs/nats v0.0.0-20221125194658-5cb70dbdbf2a // indirect
	github.com/satont/tsuwari/libs/twitch v0.0.0-20221125194658-5cb70dbdbf2a // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	golang.org/x/crypto v0.1.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gorm.io/gorm v1.24.1 // indirect
)
