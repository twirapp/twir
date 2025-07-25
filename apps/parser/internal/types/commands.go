package types

import (
	"context"

	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	model "github.com/twirapp/twir/libs/gomodels"
)

type CommandHandlerError struct {
	Err     error
	Message string
}

func (c *CommandHandlerError) Error() string {
	return c.Message
}

type CommandsHandlerResult struct {
	Result []string
}

type DefaultCommand struct {
	*model.ChannelsCommands

	Handler func(ctx context.Context, parseCtx *ParseContext) (
		*CommandsHandlerResult,
		error,
	)
	Args              []command_arguments.Arg
	ArgsDelimiter     string
	SkipToxicityCheck bool
}
