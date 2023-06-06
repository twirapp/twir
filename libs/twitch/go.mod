module github.com/satont/tsuwari/libs/twitch

go 1.20

require (
	github.com/nicklaw5/helix/v2 v2.22.1
	github.com/samber/lo v1.38.1
	github.com/satont/tsuwari/libs/config v0.0.0-00010101000000-000000000000
	github.com/satont/tsuwari/libs/grpc v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.55.0
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230530153820-e85fd2cbaebc // indirect
)

replace github.com/satont/tsuwari/libs/gomodels => ../../libs/gomodels

replace github.com/satont/tsuwari/libs/grpc => ../../libs/grpc

replace github.com/satont/tsuwari/libs/config => ../../libs/config
