package middlewares

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	cfg "github.com/satont/tsuwari/libs/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

var CheckUserAuth = func(c *fiber.Ctx) error {
	db := do.MustInvoke[*gorm.DB](di.Injector)
	config := do.MustInvoke[*cfg.Config](di.Injector)
	logger := do.MustInvoke[*zap.Logger](di.Injector)

	if c.Locals("dbUser") != nil {
		return c.Next()
	}

	headers := c.GetReqHeaders()
	apiKey := headers["Api-Key"]
	dbUser := model.Users{}

	if apiKey != "" {
		err := db.Where(`"apiKey" = ?`, apiKey).
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
		tokenSlice := strings.Split(authorizationToken, "Bearer ")
		if len(tokenSlice) < 2 {
			return fiber.NewError(401, "invalid token format")
		}

		token, err := jwt.Parse(tokenSlice[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(config.JwtAccessSecret), nil
		})
		if err != nil {
			return fiber.NewError(401, "invalid token. Probably token is expired.")
		}

		claims := token.Claims.(jwt.MapClaims)
		userId := claims["id"]

		if userId == "" {
			logger.Sugar().Info()
			logger.Sugar().Error("no userId in request")
			return fiber.NewError(401, "invalid token")
		}

		err = db.Where(`"id" = ?`, userId).
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
