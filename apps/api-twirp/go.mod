module github.com/satont/twir/apps/api-twirp

go 1.20

replace github.com/satont/twir/libs/grpc => ../../libs/grpc

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/alexedwards/scs/goredisstore v0.0.0-20230327161757-10d4299e3b24
	github.com/alexedwards/scs/v2 v2.5.1
	github.com/bytedance/sonic v1.9.1
	github.com/google/uuid v1.3.0
	github.com/guregu/null v4.0.0+incompatible
	github.com/nicklaw5/helix/v2 v2.22.2
	github.com/redis/go-redis/v9 v9.0.5
	github.com/rs/cors v1.9.0
	github.com/samber/lo v1.38.1
	github.com/satont/twir/libs/config v0.0.0-20230617211209-79e3285c6910
	github.com/satont/twir/libs/crypto v0.0.0-20230617211209-79e3285c6910
	github.com/satont/twir/libs/gomodels v0.0.0-20230617211209-79e3285c6910
	github.com/satont/twir/libs/grpc v0.0.0-20230617211209-79e3285c6910
	github.com/satont/twir/libs/twitch v0.0.0-20230617211209-79e3285c6910
	github.com/twitchtv/twirp v8.1.3+incompatible
	go.uber.org/fx v1.20.0
	go.uber.org/zap v1.24.0
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.1
)

require (
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	go.uber.org/dig v1.17.0 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/grpc v1.55.0 // indirect
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/protobuf v1.30.0
)
