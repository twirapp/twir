package stats

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router) fiber.Router {
	middleware := router.Group("stats")
	middleware.Get("", get)

	return middleware
}

// Stats godoc
// @Security ApiKeyAuth
// @Summary      Get some bot statistic
// @Tags         Stats
// @Produce      json
// @Success      200  {array}  statsItem
// @Failure 500 {object} types.DOCApiInternalError
// @Router       /v1/stats [get]
var get = func(c *fiber.Ctx) error {
	stats, err := handleGet()
	if err != nil {
		return err
	}
	return c.JSON(stats)
}
