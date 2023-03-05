module github.com/satont/tsuwari/apps/eventsub

go 1.20

replace (
	github.cim/satont/tsuwari/libs/config => ../../libs/config
	github.com/dnsge/twitch-eventsub-framework => ../../vendor/twitch-eventsub-framework
	github.com/satont/tsuwari/libs/grpc => ../../libs/grpc
)

require (
	github.com/dnsge/twitch-eventsub-bindings v0.0.0-20211025032511-9c30a4c90402
	github.com/dnsge/twitch-eventsub-framework v1.0.2
	github.com/samber/lo v1.37.0
	github.com/satont/tsuwari/libs/config v0.0.0-20230305122358-17fcb584d9ed
	github.com/satont/tsuwari/libs/grpc v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.24.0
	golang.ngrok.com/ngrok v1.0.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/inconshreveable/log15 v3.0.0-testing.3+incompatible // indirect
	github.com/inconshreveable/log15/v3 v3.0.0-testing.5 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mozillazg/go-httpheader v0.3.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/term v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/grpc v1.53.0 // indirect
)
