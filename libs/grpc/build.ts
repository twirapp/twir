import { $ } from 'bun'

await $`cd ../.. && bun run cli bin --pwd ./libs/grpc go tool github.com/bufbuild/buf/cmd/buf generate  --include-imports --include-wkt`
