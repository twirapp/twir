module github.com/satont/twir/apps/ytsr

go 1.21

replace github.com/satont/twir/libs/config => ../../libs/config

replace github.com/satont/twir/libs/grpc => ../../libs/grpc

require (
	github.com/raitonoberu/ytsearch v0.2.0
	github.com/samber/lo v1.39.0
	github.com/satont/twir/libs/config v0.0.0-20231203205548-e635accc6b72
	github.com/satont/twir/libs/grpc v0.0.0-20231203205548-e635accc6b72
	go.uber.org/zap v1.26.0
	google.golang.org/grpc v1.59.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20231127185646-65229373498e // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231127180814-3a041ad873d4 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)
