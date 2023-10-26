module github.com/satont/twir/apps/api

go 1.21.0

replace (
	github.com/satont/twir/libs/config => ../../libs/config
	github.com/satont/twir/libs/crypto => ../../libs/crypto
	github.com/satont/twir/libs/gomodels => ../../libs/gomodels
	github.com/satont/twir/libs/grpc => ../../libs/grpc
	github.com/satont/twir/libs/logger => ../../libs/logger
	github.com/satont/twir/libs/pubsub => ../../libs/pubsub
	github.com/satont/twir/libs/twitch => ../../libs/twitch
	github.com/satont/twir/libs/types => ../../libs/types
	github.com/satont/twir/libs/utils => ../../libs/utils
)

replace github.com/satont/twir/libs/integrations/spotify => ../../libs/integrations/spotify

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/alexedwards/scs/goredisstore v0.0.0-20230327161757-10d4299e3b24
	github.com/alexedwards/scs/v2 v2.5.1
	github.com/bytedance/sonic v1.9.1
	github.com/getsentry/sentry-go v0.25.0
	github.com/goccy/go-json v0.10.2
	github.com/golang/protobuf v1.5.3
	github.com/google/uuid v1.3.1
	github.com/guregu/null v4.0.0+incompatible
	github.com/imroc/req/v3 v3.42.1
	github.com/lib/pq v1.10.9
	github.com/minio/minio-go/v7 v7.0.62
	github.com/nicklaw5/helix/v2 v2.25.1
	github.com/raitonoberu/ytsearch v0.2.0
	github.com/redis/go-redis/v9 v9.2.1
	github.com/rs/cors v1.9.0
	github.com/samber/lo v1.38.1
	github.com/satont/twir/libs/config v0.0.0-20230617211209-79e3285c6910
	github.com/satont/twir/libs/crypto v0.0.0-20230617211209-79e3285c6910
	github.com/satont/twir/libs/gomodels v0.0.0-20230617211209-79e3285c6910
	github.com/satont/twir/libs/grpc v0.0.0-20230617211209-79e3285c6910
	github.com/satont/twir/libs/integrations/spotify v0.0.0-20230713153539-b2fe2b3b5757
	github.com/satont/twir/libs/logger v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/pubsub v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/twitch v0.0.0-20230617211209-79e3285c6910
	github.com/satont/twir/libs/types v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/utils v0.0.0-00010101000000-000000000000
	github.com/shkh/lastfm-go v0.0.0-20191215035245-89a801c244e0
	github.com/twitchtv/twirp v8.1.3+incompatible
	go.uber.org/fx v1.20.1
	golang.org/x/sync v0.4.0
	gorm.io/driver/postgres v1.5.3
	gorm.io/gorm v1.25.5
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/cloudflare/circl v1.3.3 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/gaukas/godicttls v0.0.4 // indirect
	github.com/go-redis/redis/v9 v9.0.0-rc.2 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/pprof v0.0.0-20230901174712-0191c66da455 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/ginkgo/v2 v2.12.0 // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/qtls-go1-20 v0.3.3 // indirect
	github.com/quic-go/quic-go v0.38.1 // indirect
	github.com/refraction-networking/utls v1.5.3 // indirect
	github.com/rs/xid v1.5.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	go.uber.org/dig v1.17.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/mod v0.13.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/tools v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231012201019-e917dd12ba7a // indirect
	google.golang.org/grpc v1.58.3 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/protobuf v1.31.0
)
