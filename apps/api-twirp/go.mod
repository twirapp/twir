module github.com/satont/tsuwari/apps/api-twirp

go 1.20

replace github.com/satont/tsuwari/libs/grpc => ../../libs/grpc

require (
	github.com/rs/cors v1.9.0
	github.com/satont/tsuwari/libs/grpc v0.0.0-20230617211209-79e3285c6910
	github.com/twitchtv/twirp v8.1.3+incompatible
)

require google.golang.org/protobuf v1.30.0 // indirect
