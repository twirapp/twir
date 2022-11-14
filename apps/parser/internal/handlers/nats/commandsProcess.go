package natshandler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/satont/tsuwari/apps/parser/internal/permissions"
	"github.com/satont/tsuwari/apps/parser/pkg/helpers"

	"github.com/go-redis/redis/v9"
	parserproto "github.com/satont/tsuwari/libs/nats/parser"
)

const (
	cooldownPerUser = "PER_USER"
	cooldownGlobal  = "GLOBAL"
)

func (c *NatsServiceImpl) shouldCheckCooldown(badges []string) bool {
	return !helpers.Contains(badges, "BROADCASTER") &&
		!helpers.Contains(badges, "MODERATOR") &&
		!helpers.Contains(badges, "SUBSCRIBER")
}

func (c *NatsServiceImpl) HandleProcessCommand(data parserproto.Request) *parserproto.Response {
	if !strings.HasPrefix(data.Message.Text, "!") {
		return nil
	}
	data.Message.Text = data.Message.Text[1:]

	cmds, err := c.commands.GetChannelCommands(data.Channel.Id)
	if err != nil {
		return nil
	}

	cmd := c.commands.FindByMessage(data.Message.Text, cmds)

	if cmd.Cmd == nil {
		return nil
	}

	if cmd.Cmd.Cooldown.Valid && cmd.Cmd.CooldownType == cooldownGlobal &&
		cmd.Cmd.Cooldown.Int64 > 0 &&
		c.shouldCheckCooldown(data.Sender.Badges) {
		key := fmt.Sprintf("commands:%s:cooldowns:global", cmd.Cmd.ID)
		_, rErr := c.redis.Get(context.TODO(), key).
			Result()

		if rErr == redis.Nil {
			c.redis.Set(context.TODO(), key, "", time.Duration(cmd.Cmd.Cooldown.Int64)*time.Second)
		} else {
			return nil
		}
	}

	if cmd.Cmd.Cooldown.Valid && cmd.Cmd.CooldownType == cooldownPerUser &&
		cmd.Cmd.Cooldown.Int64 > 0 &&
		c.shouldCheckCooldown(data.Sender.Badges) {
		key := fmt.Sprintf("commands:%s:cooldowns:user:%s", cmd.Cmd.ID, data.Sender.Id)
		_, rErr := c.redis.Get(context.TODO(), key).
			Result()

		if rErr == redis.Nil {
			c.redis.Set(context.TODO(), key, "", time.Duration(cmd.Cmd.Cooldown.Int64)*time.Second)
		} else {
			return nil
		}
	}

	hasPerm := permissions.UserHasPermissionToCommand(data.Sender.Badges, cmd.Cmd.Permission)

	if !hasPerm {
		return nil
	}

	result := c.commands.ParseCommandResponses(cmd, data)

	return result
}
