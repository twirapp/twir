module github.com/satont/tsuwari/libs/twitch

go 1.20

require (
	github.com/nicklaw5/helix/v2 v2.22.1
	github.com/satont/tsuwari/libs/config v0.0.0-00010101000000-000000000000
	github.com/satont/tsuwari/libs/grpc v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
)

replace github.com/satont/tsuwari/libs/gomodels => ../../libs/gomodels

replace github.com/satont/tsuwari/libs/grpc => ../../libs/grpc

replace github.com/satont/tsuwari/libs/config => ../../libs/config
