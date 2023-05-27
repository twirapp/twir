package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func (c *Middlewares) ValidateBody(
	ctx *fiber.Ctx,
	dto any,
) error {
	if err := ctx.BodyParser(dto); err != nil {
		if err.Error() == "Unprocessable Entity" {
			return fiber.NewError(400, "data not provided")
		}
		return err
	}

	if err := c.validator.Struct(dto); err != nil {
		return err
	}

	return nil
}
