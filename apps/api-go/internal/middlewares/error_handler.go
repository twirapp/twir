package middlewares

import (
	"encoding/json"
	"fmt"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var ErrorHandler = func(t ut.Translator) func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		switch castedErr := err.(type) {
		case validator.ValidationErrors:
			errors := []string{}
			for _, e := range castedErr {
				errors = append(errors, fmt.Sprintf("%s", e.Translate(t)))
			}
			return c.Status(fiber.StatusBadRequest).JSON(errors)
		case *json.InvalidUnmarshalError, *json.UnmarshalFieldError, *json.UnmarshalTypeError:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "bad request body"})
		case *fiber.Error:
			return c.Status(castedErr.Code).JSON(fiber.Map{"message": castedErr.Message})
		default:
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": err.Error()})
		}
	}
}
