module github.com/satont/tsuwari/apps/eventsub

go 1.20

replace (
	github.com/satont/tsuwari/libs/config => ../../libs/config
	github.com/satont/tsuwari/libs/gomodels => ../../libs/gomodels
	github.com/satont/tsuwari/libs/grpc => ../../libs/grpc
	github.com/satont/tsuwari/libs/pubsub => ../../libs/pubsub
	github.com/satont/tsuwari/libs/twitch => ../../libs/twitch
	github.com/satont/tsuwari/libs/types => ../../libs/types
)

require (
	github.com/dnsge/twitch-eventsub-bindings v1.1.0
	github.com/dnsge/twitch-eventsub-framework v1.1.1
	github.com/google/uuid v1.3.0
	github.com/lib/pq v1.10.7
	github.com/nicklaw5/helix/v2 v2.22.0
	github.com/samber/lo v1.37.0
	github.com/satont/tsuwari/libs/config v0.0.0-20230305122358-17fcb584d9ed
	github.com/satont/tsuwari/libs/gomodels v0.0.0-00010101000000-000000000000
	github.com/satont/tsuwari/libs/grpc v0.0.0-00010101000000-000000000000
	github.com/satont/tsuwari/libs/pubsub v0.0.0-00010101000000-000000000000
	github.com/satont/tsuwari/libs/twitch v0.0.0-00010101000000-000000000000
	github.com/satont/tsuwari/libs/types v0.0.0-20230305190758-592d5ff92f0d
	go.uber.org/zap v1.24.0
	golang.ngrok.com/ngrok v1.0.0
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
	gorm.io/driver/postgres v1.4.8
	gorm.io/gorm v1.24.6
)

require (
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v9 v9.0.0-rc.2 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.3 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/inconshreveable/log15 v3.0.0-testing.3+incompatible // indirect
	github.com/inconshreveable/log15/v3 v3.0.0-testing.5 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mozillazg/go-httpheader v0.3.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/goleak v1.1.12 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/exp v0.0.0-20230213192124-5e25df0256eb // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/term v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
)
