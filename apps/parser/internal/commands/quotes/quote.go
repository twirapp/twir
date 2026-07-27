package quotes

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/guregu/null"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"gorm.io/gorm"
)

const quoteIDArgName = "id"

var Quote = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "quote",
		Description: null.StringFrom("Show a quote"),
		Module:      "MANAGE",
		IsReply:     true,
	},
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name:     quoteIDArgName,
			Optional: true,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (*types.CommandsHandlerResult, error) {
		result := &types.CommandsHandlerResult{}
		quote := model.ChannelsQuotes{}
		db := parseCtx.Services.Gorm.WithContext(ctx).Where(`"channelId" = ?`, parseCtx.Channel.DBChannelID)

		if parseCtx.ArgsParser.IsExists(quoteIDArgName) {
			number, ok := parseQuoteNumber(parseCtx.ArgsParser.Get(quoteIDArgName).String())
			if !ok {
				result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.NotFound)}
				return result, nil
			}

			err := db.Where("number = ?", number).First(&quote).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.NotFound)}
				return result, nil
			}
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.CannotGet),
					Err:     err,
				}
			}
		} else {
			err := db.Order("RANDOM()").First(&quote).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.RandomEmpty)}
				return result, nil
			}
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.CannotGet),
					Err:     err,
				}
			}
		}

		result.Result = []string{fmt.Sprintf("#%d %s", quote.Number, quote.Text)}
		return result, nil
	},
}

func parseQuoteNumber(input string) (int, bool) {
	number, err := strconv.Atoi(strings.TrimPrefix(strings.TrimSpace(input), "#"))
	if err != nil || number < 1 {
		return 0, false
	}

	return number, true
}
