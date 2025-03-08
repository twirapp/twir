import { $ } from 'bun'

await $`cd ../.. && bun run cli bin --pwd ./libs/api go tool github.com/bufbuild/buf/cmd/buf generate --template buf.gen.api.yaml api.proto`
await $`cd ../.. && bun run cli bin --pwd ./libs/api go tool github.com/bufbuild/buf/cmd/buf generate --path messages --template buf.gen.messages.yaml`
