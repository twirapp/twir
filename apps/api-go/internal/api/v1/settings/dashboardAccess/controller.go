package dashboardaccess

import (
	"fmt"
	"time"
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/satont/tsuwari/apps/api-go/internal/middlewares"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func getDashboardAccessCacheKey(channelId string) string {
	return fmt.Sprintf("channels:dashboardAccess:%s", channelId)
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

	middleware.Get("", dashboardAccessList, func(c *fiber.Ctx) error {
		users, err := handleGet(c.Params("channelId"), services)
		if err != nil {
			return nil
		}

		return c.JSON(users)
	})

	limit := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			dbUser := c.Locals("dbUser").(model.Users)
			return fmt.Sprintf("fiber:limiter:feedback:%s", dbUser.ID)
		},
		LimitReached: func(c *fiber.Ctx) error {
			header := c.GetRespHeader("Retry-After", "2")
			return c.Status(429).JSON(fiber.Map{"message": fmt.Sprintf("wait %s seconds", header)})
		},
		Storage: services.RedisStorage,
	})

	middleware.Post("", limit, func(c *fiber.Ctx) error {
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

		services.RedisStorage.Delete(getDashboardAccessCacheKey(c.Params("channelId")))

		return c.JSON(entity)
	})

	middleware.Delete(":entityId", func(c *fiber.Ctx) error {
		err := handleDelete(c.Params("entityId"), services)
		if err != nil {
			return err
		}
		return c.SendStatus(200)
	})

	return middleware
}
