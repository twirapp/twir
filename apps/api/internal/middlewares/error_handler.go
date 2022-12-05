package middlewares

import (
	"encoding/json"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var ErrorHandler = func(t ut.Translator, logger *zap.Logger) func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		switch castedErr := err.(type) {
		case validator.ValidationErrors:
			errors := []string{}
			for _, e := range castedErr {
				errors = append(errors, e.Translate(t))
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"messages": errors})
		case *json.InvalidUnmarshalError:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"messages": []string{"bad request body"}})
		case *fiber.Error:
			return c.Status(castedErr.Code).JSON(fiber.Map{"messages": []string{castedErr.Message}})
		default:
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"messages": []string{err.Error()}})
		}
	}
}
