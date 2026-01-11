package request

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
)

var supportedContentType = "text/plain"

const requestTemplate = `
const req = await fetch("%s");
if (!req.ok) {
	return "Request failed: " + req.status;
}

if (req.headers.get("content-type") !== "%s") {
	return "Unsupported content type: " + req.headers.get("content-type");
}

const response = await req.text();
return response;
`

var Request = &types.Variable{
	Name:                     "request",
	Description:              lo.ToPtr("Request third party api"),
	Example:                  lo.ToPtr("request|https://decapi.me/youtube/latest_video?id=UCjerlCIbLPQwSnYlClkjDXg"),
	DisableInCustomVariables: true,
	CommandsOnly:             true,
	CanBeUsedInRegistry:      false,
	Handler: func(
		ctx context.Context,
		parseCtx *types.VariableParseContext,
		variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		param := ""

		result := &types.VariableHandlerResult{}

		if variableData.Params != nil {
			param = *variableData.Params
		}
		if param == "" {
			return result, nil
		}

		script := fmt.Sprintf(requestTemplate, param, supportedContentType)

		req, err := parseCtx.Services.Executron.ExecuteUserCode(
			ctx,
			parseCtx.Channel.ID,
			"javascript",
			script,
		)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Request.Errors.ExecuteRequest)
			return result, nil
		}

		if req.Result != "" {
			result.Result = req.Result
		} else if req.Error != "" {
			result.Result = req.Error
		}

		return result, nil
	},
}
