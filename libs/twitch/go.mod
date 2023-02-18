module github.com/satont/tsuwari/libs/twitch

go 1.19

require (
	github.com/satont/go-helix/v2 v2.7.28
	google.golang.org/grpc v1.52.3
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/golang-jwt/jwt/v4 v4.4.3 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6 // indirect
)

replace github.com/satont/tsuwari/libs/gomodels => ../../libs/gomodels
