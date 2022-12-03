package dashboardaccess

import (
	"fmt"
	"time"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func getDashboardAccessCacheKey(channelId string) string {
	return fmt.Sprintf("fiber:cache:settings:dashboardAccess:%s", channelId)
}

func Setup(router fiber.Router, services types.Services) fiber.Router {
	middleware := router.Group("dashboard-access")

	dashboardAccessList := cache.New(cache.Config{
		Expiration: 10 * time.Minute,
		Storage:    services.RedisStorage,
		KeyGenerator: func(c *fiber.Ctx) string {
			return getDashboardAccessCacheKey(c.Params("channelId"))
		},
	})

	middleware.Get("", dashboardAccessList, get(services))

	limit := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			dbUser := c.Locals("dbUser").(model.Users)
			return fmt.Sprintf("fiber:limiter:dashboardaccess:post:%s", dbUser.ID)
		},
		LimitReached: func(c *fiber.Ctx) error {
			header := c.GetRespHeader("Retry-After", "2")
			return c.Status(429).JSON(fiber.Map{"message": fmt.Sprintf("wait %s seconds", header)})
		},
		Storage: services.RedisStorage,
	})

	middleware.Post("", limit, post(services))

	middleware.Delete(":entityId", func(c *fiber.Ctx) error {
		err := handleDelete(c.Params("entityId"), services)
		if err != nil {
			return err
		}
		services.RedisStorage.DeleteByMethod(
			getDashboardAccessCacheKey(c.Params("channelId")),
			"GET",
		)
		return c.SendStatus(200)
	})

	return middleware
}

// Dashboard Access godoc
// @Security ApiKeyAuth
// @Summary      Get channel dashboard access list
// @Tags         Dashboard Access
// @Accept       json
// @Produce      json
// @Param        channelId   path      string  true  "ChannelId"
// @Success      200  {array}  Entity
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/dashboard-access [get]
func get(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return nil
		}

		return c.JSON(users)
	}
}

// Dashboard Access godoc
// @Security ApiKeyAuth
// @Summary      Add user to dashboard access list
// @Tags         Dashboard Access
// @Accept       json
// @Produce      json
// @Param data body addUserDto true "Data"
// @Param        channelId   path      string  true  "ID of channel"
// @Success      200  {object}  Entity
// @Failure 400 {object} types.DOCApiValidationError
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/channels/{channelId}/dashboard-access [post]
func post(services types.Services) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dto := &addUserDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return err
		}

		entity, err := handlePost(c.Params("channelId"), dto, services)
		if err != nil {
			return err
		}

		services.RedisStorage.DeleteByMethod(
			getDashboardAccessCacheKey(c.Params("channelId")),
			"GET",
		)

		return c.JSON(entity)
	}
}
