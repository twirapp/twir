module github.com/satont/twir/apps/eventsub

go 1.21

replace (
	github.com/dnsge/twitch-eventsub-bindings => ../../libs/twitch-eventsub-bindings
	github.com/dnsge/twitch-eventsub-framework => ../../libs/twitch-eventsub-framework
	github.com/satont/twir/libs/config => ../../libs/config
	github.com/satont/twir/libs/gomodels => ../../libs/gomodels
	github.com/satont/twir/libs/pubsub => ../../libs/pubsub
	github.com/satont/twir/libs/twitch => ../../libs/twitch
	github.com/satont/twir/libs/types => ../../libs/types
	github.com/twirapp/twir/libs/grpc => ../../libs/grpc
)

require (
	github.com/dnsge/twitch-eventsub-bindings v1.2.0
	github.com/dnsge/twitch-eventsub-framework v1.2.4
	github.com/google/uuid v1.5.0
	github.com/lib/pq v1.10.9
	github.com/localtunnel/go-localtunnel v0.0.0-20170326223115-8a804488f275
	github.com/nicklaw5/helix/v2 v2.25.2
	github.com/redis/go-redis/v9 v9.3.0
	github.com/samber/lo v1.39.0
	github.com/satont/twir/libs/config v0.0.0-20231203205548-e635accc6b72
	github.com/satont/twir/libs/gomodels v0.0.0-20231203205548-e635accc6b72
	github.com/satont/twir/libs/pubsub v0.0.0-20231203205548-e635accc6b72
	github.com/satont/twir/libs/twitch v0.0.0-20231203205548-e635accc6b72
	github.com/satont/twir/libs/types v0.0.0-20231203205548-e635accc6b72
	github.com/twirapp/twir/libs/grpc v0.0.0-20231203205548-e635accc6b72
	go.uber.org/zap v1.26.0
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgx/v5 v5.5.1 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/mozillazg/go-httpheader v0.4.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/exp v0.0.0-20231214170342-aacd6d4b4611 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231212172506-995d672761c0 // indirect
)
