import { $ } from 'bun'

await $`cd ../.. && bun run cli bin --pwd ./libs/grpc buf generate  --include-imports --include-wkt`
