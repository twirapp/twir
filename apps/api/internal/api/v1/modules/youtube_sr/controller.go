package youtube_sr

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
	youtube "github.com/satont/tsuwari/libs/types/types/api/modules"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("youtube-sr")
	middleware.Get("", get(services))
	middleware.Post("", post(services))
	middleware.Post("/blacklist/:type", postBlacklist(services))
	middleware.Get("search", getSearch(services))

	return middleware
}

func get(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		settings, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return err
		}

		return c.JSON(settings)
	}
}

func getSearch(service types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		results, err := handleSearch(c.Query("query"), c.Query("type"))
		if err != nil {
			return err
		}

		return c.JSON(results)
	}
}

func post(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := youtube.YoutubeSettings{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			&dto,
		)
		if err != nil {
			return err
		}

		err = handlePost(c.Params("channelId"), &dto, services)
		if err != nil {
			return err
		}

		return c.SendStatus(204)
	}
}

var invalidPatchBodyError = fiber.Map{
	"message": "invalid incoming body",
}

var blackListTypes = []string{"users", "channels", "songs"}

func postBlacklist(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		blacklistType := c.Params("type")
		_, ok := lo.Find(blackListTypes, func(i string) bool {
			return i == blacklistType
		})

		if !ok {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "unknown type",
			})
		}

		var data any

		body := c.Body()

		switch blacklistType {
		case "users":
			d, err := parseBlackListBody(body, youtube.YoutubeBlacklistSettingsUsers{})
			if err != nil {
				return c.JSON(invalidPatchBodyError)
			}
			data = *d
		case "channels":
			d, err := parseBlackListBody(body, youtube.YoutubeBlacklistSettingsChannels{})
			if err != nil {
				return c.JSON(invalidPatchBodyError)
			}
			data = *d
		case "songs":
			d, err := parseBlackListBody(body, youtube.YoutubeBlacklistSettingsSongs{})
			if err != nil {
				return c.JSON(invalidPatchBodyError)
			}
			data = *d
		}

		err := handlePatch(c.Params("channelId"), blacklistType, data, services)
		if err != nil {
			return nil
		}

		return c.JSON(201)
	}
}

var blackListValidator = validator.New()

func parseBlackListBody[T any](body []byte, v T) (*T, error) {
	if err := json.Unmarshal(body, &v); err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err := blackListValidator.Struct(v); err != nil {
		return nil, err
	}

	return &v, nil
}
