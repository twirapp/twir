package middlewares

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	config "github.com/satont/tsuwari/libs/config"
	"strings"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

var CheckUserAuth = func(services types.Services) func(c *fiber.Ctx) error {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	return func(c *fiber.Ctx) error {
		if c.Locals("dbUser") != nil {
			return c.Next()
		}

		headers := c.GetReqHeaders()
		apiKey := headers["Api-Key"]
		dbUser := model.Users{}

		if apiKey != "" {
			err := services.DB.Where(`"apiKey" = ?`, apiKey).
				Preload("DashboardAccess").
				First(&dbUser).
				Error
			if err != nil {
				return c.JSON(fiber.Map{"message": "user with that api key not found"})
			}
			c.Locals("dbUser", dbUser)
			return c.Next()
		}

		authorizationToken := headers["Authorization"]
		if authorizationToken != "" {
			token, err := ExtractTokenFromHeader(authorizationToken)
			if err != nil {
				return fiber.NewError(401, "invalid token. Probably token is expired.")
			}

			claims := token.Claims.(jwt.MapClaims)
			userId := claims["id"]

			if userId == "" {
				logger.Error("no userId in request")
				return fiber.NewError(401, "invalid token")
			}

			err = services.DB.Where(`"id" = ?`, userId).
				Preload("DashboardAccess").
				First(&dbUser).
				Error
			if err != nil {
				return c.JSON(fiber.Map{"message": "user not found"})
			}
			c.Locals("dbUser", dbUser)
		}

		if dbUser.ID == "" {
			return c.Status(401).JSON(fiber.Map{"message": "unauthenticated"})
		}

		return c.Next()
	}
}

func ExtractTokenFromHeader(t string) (*jwt.Token, error) {
	tokenSlice := strings.Split(t, "Bearer ")
	if len(tokenSlice) < 2 {
		return nil, fiber.NewError(401, "invalid token format")
	}

	cfg := do.MustInvoke[config.Config](di.Injector)

	token, err := jwt.Parse(tokenSlice[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(cfg.JwtAccessSecret), nil
	})

	return token, err
}
