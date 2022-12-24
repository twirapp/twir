package feedback

import (
	"bytes"
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	cfg "github.com/satont/tsuwari/libs/config"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
)

var cannotSendFeedbackError = fiber.NewError(
	http.StatusInternalServerError,
	"cannot send feedback. Please contact developers.",
)

func handlePost(
	fromId string,
	text string,
	files []*multipart.FileHeader,
) error {
	botApi, err := do.Invoke[*tgbotapi.BotAPI](di.Injector)

	if err != nil {
		return fiber.NewError(
			400,
			"cannot send feedback because we are not currently configured this feature. Please contact bot developers",
		)
	}

	config := do.MustInvoke[*cfg.Config](di.Injector)
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	myText := fmt.Sprintf("New feedback from %s\n%s", fromId, text)

	userId, _ := strconv.Atoi(*config.FeedbackTelegramUserID)

	if len(files) == 0 {
		msg := tgbotapi.NewMessage(int64(userId), myText)
		_, err := botApi.Send(msg)
		if err != nil {
			logger.Error(err)
			return cannotSendFeedbackError
		}
	} else {
		media := make([]interface{}, 0, len(files))

		for _, f := range files {
			file, _ := f.Open()
			defer file.Close()
			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, file); err != nil {
				return fiber.NewError(http.StatusInternalServerError, "cannot read file")
			}
			newMedia := tgbotapi.NewInputMediaPhoto(tgbotapi.FileBytes{Name: f.Filename, Bytes: buf.Bytes()})
			media = append(media, newMedia)
		}

		mediaGroup := tgbotapi.NewMediaGroup(int64(userId), media)
		_, err := botApi.SendMediaGroup(mediaGroup)
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, "cannot send feedback due internal error")
		}
		botApi.Send(tgbotapi.NewMessage(int64(userId), myText))
	}

	return nil
}
