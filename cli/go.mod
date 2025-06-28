module github.com/twirapp/twir/cli

go 1.24.1

replace (
	github.com/satont/twir/libs/config => ../libs/config
	github.com/satont/twir/libs/migrations => ../libs/migrations
	github.com/twirapp/twir/libs/grpc => ../libs/grpc
)

require (
	github.com/99designs/gqlgen v0.17.57
	github.com/ClickHouse/clickhouse-go/v2 v2.37.2
	github.com/Masterminds/semver/v3 v3.3.0
	github.com/goccy/go-json v0.10.3
	github.com/lib/pq v1.10.9
	github.com/pressly/goose/v3 v3.18.0
	github.com/pterm/pterm v0.12.77
	github.com/rjeczalik/notify v0.9.3
	github.com/samber/lo v1.49.1
	github.com/satont/twir/libs/config v0.0.0-20240201110132-12475b437e7a
	github.com/satont/twir/libs/migrations v0.0.0-00010101000000-000000000000
	github.com/twirapp/twir/libs/grpc v0.0.0-00010101000000-000000000000
	github.com/urfave/cli/v2 v2.27.5
	golang.org/x/sync v0.15.0
)

require (
	buf.build/go/app v0.1.0 // indirect
	buf.build/go/interrupt v1.1.0 // indirect
	buf.build/go/protovalidate v0.13.1 // indirect
	buf.build/go/standard v0.1.0 // indirect
	gotest.tools/v3 v3.5.2 // indirect
)

require (
	atomicgo.dev/cursor v0.2.0 // indirect
	atomicgo.dev/keyboard v0.2.9 // indirect
	atomicgo.dev/schedule v0.1.0 // indirect
	buf.build/gen/go/bufbuild/bufplugin/protocolbuffers/go v1.36.6-20250121211742-6d880cc6cc8d.1 // indirect
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.6-20250613105001-9f2d3c737feb.1 // indirect
	buf.build/gen/go/bufbuild/registry/connectrpc/go v1.18.1-20250606164443-9d1800bf4ccc.1 // indirect
	buf.build/gen/go/bufbuild/registry/protocolbuffers/go v1.36.6-20250606164443-9d1800bf4ccc.1 // indirect
	buf.build/gen/go/pluginrpc/pluginrpc/protocolbuffers/go v1.36.6-20241007202033-cf42259fcbfc.1 // indirect
	buf.build/go/bufplugin v0.9.0 // indirect
	buf.build/go/protoyaml v0.6.0 // indirect
	buf.build/go/spdx v0.2.0 // indirect
	cel.dev/expr v0.24.0 // indirect
	connectrpc.com/connect v1.18.1 // indirect
	connectrpc.com/otelconnect v0.7.2 // indirect
	dario.cat/mergo v1.0.1 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/AndreasBriese/bbloom v0.0.0-20190825152654-46b345b51c96 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/ClickHouse/ch-go v0.66.1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/sprig/v3 v3.3.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/alecthomas/chroma/v2 v2.14.0 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/aryann/difflib v0.0.0-20210328193216-ff5ff6dc229b // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bufbuild/buf v1.55.1 // indirect
	github.com/bufbuild/protocompile v0.14.1 // indirect
	github.com/bufbuild/protoplugin v0.0.0-20250218205857-750e09ce93e1 // indirect
	github.com/caddyserver/caddy/v2 v2.9.1 // indirect
	github.com/caddyserver/certmagic v0.21.6 // indirect
	github.com/caddyserver/zerossl v0.1.3 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chavacava/garif v0.1.0 // indirect
	github.com/chzyer/readline v1.5.1 // indirect
	github.com/containerd/console v1.0.4 // indirect
	github.com/containerd/errdefs v1.0.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.16.3 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/dgraph-io/badger v1.6.2 // indirect
	github.com/dgraph-io/badger/v2 v2.2007.4 // indirect
	github.com/dgraph-io/ristretto v0.1.0 // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/dlclark/regexp2 v1.11.5 // indirect
	github.com/docker/cli v28.2.2+incompatible // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/docker v28.2.2+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.9.3 // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/francoispqt/gojay v1.2.13 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/go-chi/chi/v5 v5.2.1 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-jose/go-jose/v3 v3.0.3 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/gofrs/flock v0.12.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.2.4 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/cel-go v0.25.0 // indirect
	github.com/google/certificate-transparency-go v1.1.8-0.20240110162603-74a5dd331745 // indirect
	github.com/google/go-containerregistry v0.20.6 // indirect
	github.com/google/go-tpm v0.9.0 // indirect
	github.com/google/go-tspi v0.3.0 // indirect
	github.com/google/pprof v0.0.0-20250607225305-033d6d78b36a // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/pgx/v4 v4.18.3 // indirect
	github.com/jdx/go-netrc v1.0.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/klauspost/pgzip v1.2.6 // indirect
	github.com/libdns/libdns v0.2.2 // indirect
	github.com/lithammer/fuzzysearch v1.1.8 // indirect
	github.com/manifoldco/promptui v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/mgechev/dots v0.0.0-20210922191527-e955255bf517 // indirect
	github.com/mgechev/revive v1.9.0 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mholt/acmez/v3 v3.0.0 // indirect
	github.com/miekg/dns v1.1.62 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/onsi/ginkgo/v2 v2.23.4 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/paulmach/orb v0.11.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pires/go-proxyproto v0.7.1-0.20240628150027-b718e7ce4964 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.52.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/rs/xid v1.5.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/satont/twir/libs/crypto v0.0.0-20240201110132-12475b437e7a // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/segmentio/encoding v0.5.1 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/slackhq/nebula v1.6.1 // indirect
	github.com/smallstep/certificates v0.26.1 // indirect
	github.com/smallstep/go-attestation v0.4.4-0.20240109183208-413678f90935 // indirect
	github.com/smallstep/nosql v0.6.1 // indirect
	github.com/smallstep/pkcs7 v0.0.0-20231024181729-3b98ecc1ca81 // indirect
	github.com/smallstep/scep v0.0.0-20231024192529-aee96d7ad34d // indirect
	github.com/smallstep/truststore v0.13.0 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.7.0 // indirect
	github.com/spf13/cobra v1.9.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/tailscale/tscert v0.0.0-20240608151842-d3f834017e53 // indirect
	github.com/tetratelabs/wazero v1.9.0 // indirect
	github.com/twitchtv/twirp v8.1.3+incompatible // indirect
	github.com/urfave/cli v1.22.16 // indirect
	github.com/vbatts/tar-split v0.12.1 // indirect
	github.com/vektah/gqlparser/v2 v2.5.19 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	github.com/yuin/goldmark v1.7.8 // indirect; indirectdata
	github.com/yuin/goldmark-highlighting/v2 v2.0.0-20230729083705-37449abec8cc // indirect
	github.com/zeebo/assert v1.3.0 // indirect
	github.com/zeebo/blake3 v0.2.4 // indirect
	go.etcd.io/bbolt v1.3.10 // indirect
	go.lsp.dev/jsonrpc2 v0.10.0 // indirect
	go.lsp.dev/pkg v0.0.0-20210717090340-384b27a52fb2 // indirect
	go.lsp.dev/protocol v0.12.0 // indirect
	go.lsp.dev/uri v0.3.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.61.0 // indirect
	go.opentelemetry.io/contrib/propagators/autoprop v0.42.0 // indirect
	go.opentelemetry.io/contrib/propagators/aws v1.17.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.32.0 // indirect
	go.opentelemetry.io/contrib/propagators/jaeger v1.17.0 // indirect
	go.opentelemetry.io/contrib/propagators/ot v1.17.0 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.31.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/sdk v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	go.opentelemetry.io/proto/otlp v1.7.0 // indirect
	go.step.sm/cli-utils v0.9.0 // indirect
	go.step.sm/crypto v0.45.0 // indirect
	go.step.sm/linkedca v0.20.1 // indirect
	go.uber.org/automaxprocs v1.6.0 // indirect
	go.uber.org/mock v0.5.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	go.uber.org/zap/exp v0.3.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/crypto/x509roots/fallback v0.0.0-20241104001025-71ed71b4faf9 // indirect
	golang.org/x/exp v0.0.0-20250606033433-dcc06ee1d476 // indirect
	golang.org/x/mod v0.25.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/term v0.32.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	golang.org/x/time v0.7.0 // indirect
	golang.org/x/tools v0.34.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/grpc v1.72.2 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.5.1 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	howett.net/plist v1.0.0 // indirect
	pluginrpc.com/pluginrpc v0.5.0 // indirect
)

tool (
	github.com/bufbuild/buf/cmd/buf
	github.com/caddyserver/caddy/v2/cmd/caddy
	github.com/mgechev/revive
	github.com/twitchtv/twirp/protoc-gen-twirp
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	google.golang.org/protobuf/cmd/protoc-gen-go
)
