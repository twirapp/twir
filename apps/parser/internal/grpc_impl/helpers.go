package grpc_impl

import (
	"strconv"
	"time"

	"github.com/satont/tsuwari/apps/parser/pkg/helpers"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"go.uber.org/zap"
)

func (c *parserGrpcServer) shouldCheckCooldown(badges []string) bool {
	return !helpers.Contains(badges, "BROADCASTER") &&
		!helpers.Contains(badges, "MODERATOR") &&
		!helpers.Contains(badges, "SUBSCRIBER")
}

func (c *parserGrpcServer) isUserRestrictedToUse(
	data *parser.ProcessCommandRequest,
	command *model.ChannelsCommands,
) bool {
	if command.Restrictions == nil || len(command.Restrictions) == 0 {
		return false
	}

	userStats := model.UsersStats{}
	err := c.commands.Db.
		Where(`"userId" = ? AND "channelId" = ?`, data.Sender.Id, data.Channel.Id).
		Find(&userStats).Error
	if err != nil {
		zap.L().Sugar().Error(err)
		return true
	}

	if userStats.ID == "" {
		return true
	}

	result := false
	for _, r := range command.Restrictions {
		if r.Type == "WATCHED" {
			watched, err := strconv.ParseInt(r.Value, 10, 64)
			if err != nil {
				zap.L().Sugar().Error(err)
				result = false
				continue
			}
			result = userStats.Watched*int64(time.Millisecond) < watched
			continue
		}

		if r.Type == "MESSAGES" {
			messages, err := strconv.ParseInt(r.Value, 10, 64)
			if err != nil {
				zap.L().Sugar().Error(err)
				result = false
				continue
			}
			result = int64(userStats.Messages)*int64(time.Millisecond) < messages
			continue
		}
	}

	return result
}
