package types

import (
	"context"
	model "github.com/satont/twir/libs/gomodels"
)

type CommandsHandlerResult struct {
	Result []string
}

type DefaultCommand struct {
	*model.ChannelsCommands

	Handler func(ctx context.Context, parseCtx *ParseContext) *CommandsHandlerResult
}
