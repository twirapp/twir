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
		return c.Status(400).JSON(fiber.Map{
			"message": "bad request body",
		})
	}

	if err := v.Struct(dto); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := []string{}
		for _, e := range validationErrors {
			errors = append(
				errors,
				fmt.Sprintf(
					"%s",
					e.Translate(translator),
				),
			)
		}
		c.Status(fiber.StatusBadRequest).JSON(errors)
		return nil
	}

	return nil
}
