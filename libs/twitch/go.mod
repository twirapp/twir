module github.com/satont/twir/libs/twitch

go 1.21

require (
	github.com/nicklaw5/helix/v2 v2.25.1
	github.com/samber/lo v1.38.1
	github.com/satont/twir/libs/config v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/grpc v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.58.3
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231012201019-e917dd12ba7a // indirect
)

replace github.com/satont/twir/libs/gomodels => ../../libs/gomodels

replace github.com/satont/twir/libs/grpc => ../../libs/grpc

replace github.com/satont/twir/libs/config => ../../libs/config
