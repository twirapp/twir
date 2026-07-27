package quotes

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"gorm.io/gorm"
)

const quoteTextArgName = "text"

var AddQuote = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "quote add",
		Aliases:     pq.StringArray{"quote +"},
		Description: null.StringFrom("Add a quote"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MANAGE",
		IsReply:     true,
	},
	Args: []command_arguments.Arg{
		command_arguments.VariadicString{
			Name:     quoteTextArgName,
			Optional: true,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (*types.CommandsHandlerResult, error) {
		result := &types.CommandsHandlerResult{}
		if !parseCtx.ArgsParser.IsExists(quoteTextArgName) {
			result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.EmptyText)}
			return result, nil
		}

		text := strings.TrimSpace(parseCtx.ArgsParser.Get(quoteTextArgName).String())
		if text == "" {
			result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.EmptyText)}
			return result, nil
		}

		now := time.Now()
		quoteID, err := uuid.NewV7()
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.CannotAdd),
				Err:     err,
			}
		}
		quote := model.ChannelsQuotes{
			ID:          quoteID.String(),
			ChannelID:   parseCtx.Channel.DBChannelID,
			Text:        text,
			CreatorID:   null.StringFrom(parseCtx.Sender.ID),
			CreatorName: null.StringFrom(parseCtx.Sender.DisplayName),
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if parseCtx.ChannelStream != nil {
			quote.GameID = null.StringFrom(parseCtx.ChannelStream.GameId)
			quote.GameName = null.StringFrom(parseCtx.ChannelStream.GameName)
		}

		err = parseCtx.Services.Gorm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec("SELECT pg_advisory_xact_lock(hashtext(?))", quote.ChannelID).Error; err != nil {
				return err
			}

			if err := tx.
				Model(&model.ChannelsQuotes{}).
				Where(`"channelId" = ?`, quote.ChannelID).
				Select("COALESCE(MAX(number), 0) + 1").
				Scan(&quote.Number).
				Error; err != nil {
				return err
			}

			return tx.Create(&quote).Error
		})
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.CannotAdd),
				Err:     err,
			}
		}

		result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Add.Added.SetVars(
			locales.KeysCommandsQuotesAddAddedVars{Number: quote.Number},
		))}
		return result, nil
	},
}
