module tsuwari/parser

go 1.18

require (
	github.com/getsentry/sentry-go v0.14.0
	github.com/go-redis/redis/v9 v9.0.0-rc.1
	github.com/guregu/null v4.0.0+incompatible
	github.com/nats-io/nats.go v1.18.0
	github.com/satont/go-helix/v2 v2.7.14
	github.com/satont/tsuwari/nats/bots v0.0.0-00010101000000-000000000000
	github.com/satont/tsuwari/nats/dota v0.0.0-00010101000000-000000000000
	github.com/satont/tsuwari/nats/eval v0.0.0-00010101000000-000000000000
	github.com/satori/go.uuid v1.2.0
	github.com/tidwall/gjson v1.14.3
	go.uber.org/zap v1.23.0
	google.golang.org/protobuf v1.28.1
	gorm.io/driver/postgres v1.4.4
	gorm.io/gorm v1.24.0
)

require (
	github.com/cheekybits/genny v1.0.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.2 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/lucas-clemente/quic-go v0.29.2 // indirect
	github.com/marten-seemann/qpack v0.3.0 // indirect
	github.com/marten-seemann/qtls-go1-16 v0.1.5 // indirect
	github.com/marten-seemann/qtls-go1-17 v0.1.2 // indirect
	github.com/marten-seemann/qtls-go1-18 v0.1.3 // indirect
	github.com/marten-seemann/qtls-go1-19 v0.1.1 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/mod v0.6.0 // indirect
	golang.org/x/tools v0.1.12 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/imroc/req/v3 v3.25.0
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.12.0 // indirect
	github.com/jackc/pgx/v4 v4.17.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/leighmacdonald/steamid/v2 v2.2.0
	github.com/lib/pq v1.10.7
	github.com/nats-io/nats-server/v2 v2.8.4 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/samber/lo v1.33.0
	github.com/satont/tsuwari/nats/parser v0.0.0
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/exp v0.0.0-20221019170559-20944726eadf // indirect
	golang.org/x/net v0.1.0 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	tsuwari/config v0.0.0
	tsuwari/models v0.0.0
)

replace github.com/satont/tsuwari/nats/parser => ../../libs/nats/parser

replace github.com/satont/tsuwari/nats/bots => ../../libs/nats/bots

replace github.com/satont/tsuwari/nats/dota => ../../libs/nats/dota

replace github.com/satont/tsuwari/nats/eval => ../../libs/nats/eval

replace tsuwari/config => ../../libs/config

replace tsuwari/models => ../../libs/gomodels
