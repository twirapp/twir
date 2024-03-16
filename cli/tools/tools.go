//go:build tools
// +build tools

package tools

import (
	"github.com/bufbuild/buf/cmd/buf"
	"github.com/caddyserver/caddy/v2/cmd/caddy"
	"github.com/twitchtv/twirp/protoc-gen-twirp"
	"google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	"google.golang.org/protobuf/cmd/protoc-gen-go"
)
