module github.com/satont/tsuwari/apps/parser

go 1.18

require (
	github.com/Masterminds/squirrel v1.5.3
	github.com/getsentry/sentry-go v0.18.0
	github.com/go-redis/redis/v9 v9.0.0-rc.2
	github.com/google/uuid v1.3.0
	github.com/guregu/null v4.0.0+incompatible
	github.com/jmoiron/sqlx v1.3.5
	github.com/nicklaw5/helix/v2 v2.22.0
	github.com/prometheus/client_golang v1.14.0
	github.com/samber/do v1.6.0
	github.com/satont/tsuwari/apps/bots v0.0.0-20230306140052-60f1293647c4
	github.com/satont/tsuwari/libs/gomodels v0.0.0-20221114143619-e5e207524b96
	github.com/satont/tsuwari/libs/grpc v0.0.0-00010101000000-000000000000
	github.com/satont/tsuwari/libs/integrations/spotify v0.0.0-00010101000000-000000000000
	github.com/satont/tsuwari/libs/twitch v0.0.0-20230112134844-218713e75233
	github.com/satont/tsuwari/libs/types v0.0.0-20230103005447-f8437299a3a0
	github.com/satori/go.uuid v1.2.0
	github.com/shkh/lastfm-go v0.0.0-20191215035245-89a801c244e0
	github.com/tidwall/gjson v1.14.4
	github.com/valyala/fasttemplate v1.2.2
	go.uber.org/zap v1.24.0
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
	gorm.io/driver/postgres v1.4.8
	gorm.io/gorm v1.24.6
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.3 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/pprof v0.0.0-20230207041349-798e818bf904 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgx/v5 v5.3.0 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/onsi/ginkgo/v2 v2.8.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/qtls-go1-18 v0.2.0 // indirect
	github.com/quic-go/qtls-go1-19 v0.2.1 // indirect
	github.com/quic-go/qtls-go1-20 v0.1.1 // indirect
	github.com/quic-go/quic-go v0.32.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/mod v0.8.0 // indirect
	golang.org/x/tools v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/imroc/req/v3 v3.32.0
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/leighmacdonald/steamid/v2 v2.2.0
	github.com/lib/pq v1.10.7
	github.com/samber/lo v1.37.0
	github.com/satont/tsuwari/libs/config v0.0.0
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/exp v0.0.0-20230213192124-5e25df0256eb // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
)

replace github.com/satont/tsuwari/libs/integrations/spotify => ../../libs/integrations/spotify

replace github.com/satont/tsuwari/libs/config => ../../libs/config

replace github.com/satont/tsuwari/libs/gomodels => ../../libs/gomodels

replace github.com/satont/tsuwari/libs/types => ../../libs/types

replace github.com/satont/tsuwari/libs/grpc => ../../libs/grpc

replace github.com/satont/tsuwari/libs/twitch => ../../libs/twitch
