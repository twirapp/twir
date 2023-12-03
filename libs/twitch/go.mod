module github.com/satont/twir/libs/twitch

go 1.21

require (
	github.com/nicklaw5/helix/v2 v2.25.2
	github.com/samber/lo v1.39.0
	github.com/satont/twir/libs/config v0.0.0-20231203205548-e635accc6b72
	github.com/satont/twir/libs/grpc v0.0.0-20231203205548-e635accc6b72
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	golang.org/x/exp v0.0.0-20231127185646-65229373498e // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231127180814-3a041ad873d4 // indirect
)

replace github.com/satont/twir/libs/gomodels => ../../libs/gomodels

replace github.com/satont/twir/libs/grpc => ../../libs/grpc

replace github.com/satont/twir/libs/config => ../../libs/config
