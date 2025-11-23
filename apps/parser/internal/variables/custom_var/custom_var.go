package custom_var

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/bus-core/parser"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
)

var CustomVar = &types.Variable{
	Name:                     "customvar",
	Description:              lo.ToPtr("Custom variable"),
	Visible:                  lo.ToPtr(false),
	DisableInCustomVariables: true,
	NotCachable:              true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		if variableData.Params == nil {
			return result, nil
		}

		v := &model.ChannelsCustomvars{}
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND "name" = ?`, parseCtx.Channel.ID, variableData.Params).
			Find(v).Error
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			return result, nil
		}

		if v.ID == "" || (v.Response == "" && v.EvalValue == "") {
			return result, nil
		}

		if parseCtx.Text != nil &&
			len(*parseCtx.Text) > 0 &&
			v.Type != model.CustomVarScript &&
			slices.ContainsFunc(
				parseCtx.Sender.Roles, func(item model.ChannelRole) bool {
					return item.Type == model.ChannelRoleTypeBroadcaster || item.Type == model.ChannelRoleTypeModerator
				},
			) {
			if v.Type == model.CustomVarNumber {
				parsed, err := strconv.Atoi(*parseCtx.Text)
				if err != nil {
					return nil, fmt.Errorf(
						i18n.GetCtx(
							ctx,
							locales.Translations.Variables.CustomVar.Errors.WrongNumbers.
								SetVars(locales.KeysVariablesCustomVarErrorsWrongNumbersVars{Reason: err.Error()}),
						),
					)
				}

				v.Response = fmt.Sprint(parsed)
			} else {
				v.Response = *parseCtx.Text
			}

			if err := parseCtx.Services.Gorm.Save(&v).Error; err != nil {
				return nil, fmt.Errorf(
					i18n.GetCtx(
						ctx,
						locales.Translations.Variables.CustomVar.Errors.UpdateCustomVar,
					),
				)
			}
		}

		if v.Type == model.CustomVarScript {
			requestCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			text := v.EvalValue
			if parseCtx.Text != nil {
				text = strings.ReplaceAll(text, "$(command.param)", *parseCtx.Text)
			} else {
				text = strings.ReplaceAll(text, "$(command.param)", "")
			}

			filledWithVariablesValue, err := parseCtx.Services.Bus.Parser.ParseVariablesInText.Request(
				requestCtx,
				parser.ParseVariablesInTextRequest{
					ChannelID:     parseCtx.Channel.ID,
					ChannelName:   parseCtx.Channel.Name,
					Text:          text,
					UserID:        parseCtx.Sender.ID,
					UserLogin:     parseCtx.Sender.Name,
					UserName:      parseCtx.Sender.DisplayName,
					IsCommand:     true,
					IsInCustomVar: true,
				},
			)
			if err != nil {
				return nil, err
			}

			res, err := parseCtx.Services.Executron.ExecuteUserCode(
				ctx,
				parseCtx.Channel.ID,
				v.ScriptLanguage,
				filledWithVariablesValue.Data.Text,
			)
			if err != nil {
				parseCtx.Services.Logger.Sugar().Error(err)

				return nil, errors.New(
					i18n.GetCtx(ctx, locales.Translations.Variables.CustomVar.Errors.EvaluateVariable),
				)
			}

			if res.Result != "" {
				result.Result = res.Result
			} else if res.Error != "" {
				result.Result = res.Error
			}
		}

		if v.Type == model.CustomVarText || v.Type == model.CustomVarNumber {
			result.Result = v.Response
		}

		result.Result = lo.Substring(result.Result, 0, 1000)

		return result, nil
	},
}
