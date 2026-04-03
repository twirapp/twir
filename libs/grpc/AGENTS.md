# AGENTS.md — libs/grpc

gRPC definitions and generated code.

## OVERVIEW

Protocol Buffers definitions and generated Go/TypeScript code for inter-service communication. Services communicate via gRPC for type-safe RPC calls.

## STRUCTURE

```
libs/grpc/
├── *.proto                  # Protobuf definitions
├── */                       # Generated code per service
│   ├── *.pb.go             # Go generated
│   └── *_grpc.pb.go        # gRPC generated
├── package.json            # JS package
├── go.mod
└── ...
```

## CONVENTIONS

### Adding New Service

1. Create `{service}.proto`:

```protobuf
syntax = "proto3";
package service;

service MyService {
  rpc Method (Request) returns (Response);
}

message Request { ... }
message Response { ... }
```

2. Generate code:

```bash
# Go
go generate ./...

# TypeScript (if needed)
bun run codegen
```

## ANTI-PATTERNS

- **DO NOT** edit `*.pb.go` files — modify `.proto` and regenerate
- **DO NOT** edit `*_grpc.pb.go` files — modify `.proto` and regenerate

## NOTES

- Generated files have "DO NOT EDIT" headers
- Services in `apps/*` implement these interfaces
- Some frontend libs use generated TS
