package shared

import (
	"context"

	"github.com/twirapp/twir/apps/parser/internal/types"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

const (
	PlatformTwitch = platformentity.PlatformTwitch
	PlatformKick   = platformentity.PlatformKick

	defaultPlatform = platformentity.Platform("*")
)

func UnsupportedVariableHandler(
	ctx context.Context,
	parseCtx *types.VariableParseContext,
	variableData *types.VariableData,
) (*types.VariableHandlerResult, error) {
	return &types.VariableHandlerResult{Result: "not supported on this platform"}, nil
}

func DefaultPlatformHandler(handler types.VariableHandler) (platformentity.Platform, types.VariableHandler) {
	return defaultPlatform, handler
}

func HandlerByPlatform(handlers map[platformentity.Platform]types.VariableHandler) types.VariableHandler {
	return func(
		ctx context.Context,
		parseCtx *types.VariableParseContext,
		variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		if handler, ok := handlers[parseCtx.Platform]; ok {
			return handler(ctx, parseCtx, variableData)
		}

		if handler, ok := handlers[defaultPlatform]; ok {
			return handler(ctx, parseCtx, variableData)
		}

		return UnsupportedVariableHandler(ctx, parseCtx, variableData)
	}
}
