package request

import (
	"context"
	"fmt"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var supportedContentType = "text/plain"

var Request = &types.Variable{
	Name:                     "request",
	Description:              lo.ToPtr("Request third party api"),
	Example:                  lo.ToPtr("request|https://decapi.me/youtube/latest_video?id=UCjerlCIbLPQwSnYlClkjDXg"),
	DisableInCustomVariables: true,
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

		request, err := req.Get(param)
		if err != nil {
			result.Result = fmt.Sprintf(
				`Cannot fetch %s %s`,
				param,
				"network error, probably url is wrong or server is down",
			)
			return result, nil
		}
		if !request.IsSuccessState() {
			result.Result = fmt.Sprintf(`Request to %s failed: %s`, param, request.String())
			return result, nil
		}

		responseContentType := request.GetContentType()
		if !strings.HasPrefix(responseContentType, supportedContentType) {
			result.Result = fmt.Sprintf(
				`%s responded with "%s", must respond with "%s" content`,
				param,
				responseContentType,
				supportedContentType,
			)
			return result, nil
		}

		result.Result = request.String()

		return result, nil
	},
}
