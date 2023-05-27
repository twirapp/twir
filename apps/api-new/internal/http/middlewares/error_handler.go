package middlewares

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-new/internal/http/helpers"
)

func (c *Middlewares) ErrorHandler(ctx *fiber.Ctx, err error) error {
	switch castedErr := err.(type) {
	case validator.ValidationErrors:
		var errors []string
		for _, e := range castedErr {
			errors = append(errors, e.Translate(c.translator))
		}
		return ctx.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"messages": errors})
	case *json.InvalidUnmarshalError:
		c.logger.Error(err)
		return ctx.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"messages": []string{"bad request body"}})
	case helpers.BusinessError:
		return ctx.
			Status(castedErr.StatusCode).
			JSON(fiber.Map{"messages": castedErr.Messages})
	case *fiber.Error:
		return ctx.
			Status(castedErr.Code).
			JSON(fiber.Map{"messages": []string{castedErr.Message}})
	default:
		return ctx.
			Status(fiber.StatusBadGateway).
			JSON(fiber.Map{"messages": []string{err.Error()}})
	}
}
