module github.com/satont/twir/libs/twitch

go 1.21

replace github.com/nicklaw5/helix/v2 => ../helix

require (
	github.com/nicklaw5/helix/v2 v2.25.3
	github.com/samber/lo v1.39.0
	github.com/satont/twir/libs/config v0.0.0-20240126231400-72985ccc25a5
	github.com/twirapp/twir/libs/grpc v0.0.0-20240126231400-72985ccc25a5
	google.golang.org/grpc v1.62.0
	google.golang.org/protobuf v1.33.0
)

require (
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240304212257-790db918fca8 // indirect
)

replace github.com/satont/twir/libs/gomodels => ../../libs/gomodels

replace github.com/twirapp/twir/libs/grpc => ../../libs/grpc

replace github.com/satont/twir/libs/config => ../../libs/config
