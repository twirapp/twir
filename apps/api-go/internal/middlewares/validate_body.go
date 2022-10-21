package middlewares

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidateBody[T any](
	c *fiber.Ctx,
	v *validator.Validate,
	translator ut.Translator,
	dto *T,
) error {
	if err := c.BodyParser(dto); err != nil {
		return err
	}

	if err := v.Struct(dto); err != nil {
		return err
	}

	return nil
}
