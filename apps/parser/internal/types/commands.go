package types

import (
	"context"

	model "github.com/satont/twir/libs/gomodels"
)

type CommandHandlerError struct {
	Message string
	Err     error
}

func (c *CommandHandlerError) Error() string {
	return c.Message
}

type CommandsHandlerResult struct {
	Result []string
}

type DefaultCommand struct {
	*model.ChannelsCommands

	Handler func(ctx context.Context, parseCtx *ParseContext) (*CommandsHandlerResult, error)
}
