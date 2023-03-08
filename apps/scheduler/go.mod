module github.com/satont/tsuwari/apps/scheduler

go 1.20

replace (
	github.com/satont/tsuwari/libs/config => ../../libs/config
	github.com/satont/tsuwari/libs/gomodels => ../../libs/gomodels
	github.com/satont/tsuwari/libs/grpc => ../../libs/grpc
	github.com/satont/tsuwari/libs/twitch => ../../libs/twitch
)

require (
	github.com/google/uuid v1.3.0
	github.com/guregu/null v4.0.0+incompatible
	github.com/lib/pq v1.10.7
	github.com/samber/lo v1.37.0
	github.com/satont/go-helix/v2 v2.7.28
	github.com/satont/tsuwari/libs/config v0.0.0-20230302140714-704247d5bf81
	github.com/satont/tsuwari/libs/gomodels v0.0.0-00010101000000-000000000000
	github.com/satont/tsuwari/libs/grpc v0.0.0-20230302140714-704247d5bf81
	github.com/satont/tsuwari/libs/twitch v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.24.0
	google.golang.org/protobuf v1.28.1
	gorm.io/driver/postgres v1.4.8
	gorm.io/gorm v1.24.6
)

require (
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.3 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
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
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/grpc v1.53.0 // indirect
)
