module github.com/twirapp/twir/libs/cache

go 1.23.0

toolchain go1.23.2

replace (
	github.com/satont/twir/libs/config => ../config
	github.com/satont/twir/libs/gomodels => ../gomodels
	github.com/satont/twir/libs/twitch => ../twitch
	github.com/twirapp/twir/libs/grpc => ../grpc
	github.com/twirapp/twir/libs/integrations => ../integrations
)

require (
	github.com/goccy/go-json v0.10.3
	github.com/nicklaw5/helix/v2 v2.30.0
	github.com/redis/go-redis/v9 v9.6.1
	github.com/samber/lo v1.47.0
	github.com/satont/twir/libs/config v0.0.0-20241105005614-f68dbf11ade1
	github.com/satont/twir/libs/gomodels v0.0.0-20240208154120-fc098a9e20a2
	github.com/satont/twir/libs/twitch v0.0.0-00010101000000-000000000000
	github.com/twirapp/twir/libs/grpc v0.0.0-20241105005614-f68dbf11ade1
	github.com/twirapp/twir/libs/integrations v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.28.0
	go.opentelemetry.io/otel/trace v1.28.0
	golang.org/x/sync v0.8.0
	gorm.io/gorm v1.25.12
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudflare/circl v1.5.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.1 // indirect
	github.com/google/pprof v0.0.0-20241101162523-b92577c0c142 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/imroc/req/v3 v3.48.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/onsi/ginkgo/v2 v2.21.0 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.48.1 // indirect
	github.com/refraction-networking/utls v1.6.7 // indirect
	github.com/satont/twir/libs/types v0.0.0-20241201233432-6c61559e1688 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	go.uber.org/mock v0.5.0 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	golang.org/x/tools v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241104194629-dd2ea8efbc28 // indirect
	google.golang.org/grpc v1.67.1 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)
