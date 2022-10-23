package middlewares

import (
	"fmt"

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
		if err.Error() == "Unprocessable Entity" {
			return fiber.NewError(400, "data not provided")
		}
		return err
	}

	if err := v.Struct(dto); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
