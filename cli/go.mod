module github.com/twirapp/twir/cli

go 1.21.5

replace (
	github.com/satont/twir/libs/config => ../libs/config
	github.com/satont/twir/libs/migrations => ../libs/migrations
)

require (
	github.com/lib/pq v1.10.9
	github.com/pressly/goose/v3 v3.18.0
	github.com/pterm/pterm v0.12.77
	github.com/rjeczalik/notify v0.9.3
	github.com/samber/lo v1.39.0
	github.com/satont/twir/libs/config v0.0.0-20240201110132-12475b437e7a
	github.com/satont/twir/libs/migrations v0.0.0-20240201110132-12475b437e7a
	github.com/urfave/cli/v2 v2.27.1
	golang.org/x/sync v0.6.0
)

require (
	atomicgo.dev/cursor v0.2.0 // indirect
	atomicgo.dev/keyboard v0.2.9 // indirect
	atomicgo.dev/schedule v0.1.0 // indirect
	github.com/containerd/console v1.0.4-0.20230313162750-1ae8d489ac81 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.6 // indirect
	github.com/lithammer/fuzzysearch v1.1.8 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/rivo/uniseg v0.4.6 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/satont/twir/libs/crypto v0.0.0-20240201110132-12475b437e7a // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	github.com/xrash/smetrics v0.0.0-20231213231151-1d8dd44e695e // indirect
	go.opentelemetry.io/otel v1.23.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/exp v0.0.0-20240119083558-1b970713d09a // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/term v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
