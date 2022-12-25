package grpc_impl

import "github.com/satont/tsuwari/apps/parser/pkg/helpers"

func (c *parserGrpcServer) shouldCheckCooldown(badges []string) bool {
	return !helpers.Contains(badges, "BROADCASTER") &&
		!helpers.Contains(badges, "MODERATOR") &&
		!helpers.Contains(badges, "SUBSCRIBER")
}
