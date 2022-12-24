package middlewares

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
)

func ValidateBody[T any](
	c *fiber.Ctx,
	dto *T,
) error {
	v := do.MustInvoke[*validator.Validate](di.Injector)
	//translator := do.MustInvoke[ut.Translator](di.Injector)

	if err := c.BodyParser(dto); err != nil {
		if err.Error() == "Unprocessable Entity" {
			return fiber.NewError(400, "data not provided")
		}
		return err
	}

	if err := v.Struct(dto); err != nil {
		return err
	}

	return nil
}
