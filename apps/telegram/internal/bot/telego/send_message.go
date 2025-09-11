package telego

import (
	"context"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/twirapp/twir/apps/telegram/internal/bot"
)

func (s *Service) SendMessage(ctx context.Context, input bot.SendMessageInput) error {
	params := tu.Message(
		tu.ID(input.ChatID),
		input.Text,
	)

	if input.ParseMode != nil {
		params.ParseMode = string(*input.ParseMode)
	}

	_, err := s.b.SendMessage(
		ctx,
		params,
	)

	return err
}

func (s *Service) SendPhoto(ctx context.Context, input bot.SendPhotoInput) error {
	params := tu.Photo(
		tu.ID(input.ChatID),
		telego.InputFile{
			URL: input.PhotoURL,
		},
	)
	params.Caption = input.Text

	if input.ParseMode != nil {
		params.ParseMode = string(*input.ParseMode)
	}

	_, err := s.b.SendPhoto(
		ctx,
		params,
	)

	return err
}
