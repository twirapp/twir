module github.com/satont/twir/libs/twitch

go 1.21

require (
	github.com/nicklaw5/helix/v2 v2.25.2
	github.com/samber/lo v1.39.0
	github.com/satont/twir/libs/config v0.0.0-20231203205548-e635accc6b72
	github.com/twirapp/twir/libs/grpc v0.0.0-20231203205548-e635accc6b72
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
)

require (
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	golang.org/x/exp v0.0.0-20231214170342-aacd6d4b4611 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231212172506-995d672761c0 // indirect
)

replace github.com/satont/twir/libs/gomodels => ../../libs/gomodels

replace github.com/twirapp/twir/libs/grpc => ../../libs/grpc

replace github.com/satont/twir/libs/config => ../../libs/config
