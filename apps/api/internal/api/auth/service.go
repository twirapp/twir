package auth

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/samber/lo"

	"github.com/samber/do"
	"github.com/satont/twir/apps/api/internal/di"
	"github.com/satont/twir/apps/api/internal/interfaces"
	cfg "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/twitch"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/api/internal/types"
	uuid "github.com/satori/go.uuid"
)

type Tokens struct {
	UserId       string `json:"userId"`
	RefreshToken string `json:"refreshToken,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
}

const (
	accessLifeTime  = 10 * time.Minute
	refreshLifeTime = 31 * 24 * time.Hour
)

func handleGetToken(code string, services types.Services) (*Tokens, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

	twitchClient, err := twitch.NewAppClient(config, tokensGrpc)
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	resp, err := twitchClient.RequestUserAccessToken(code)
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusUnauthorized, "cannot get user tokens")
	}

	twitchClient.SetUserAccessToken(resp.Data.AccessToken)

	users, err := twitchClient.GetUsers(&helix.UsersParams{})
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(401, "cannot get user tokens")
	}

	if len(users.Data.Users) == 0 {
		return nil, fiber.NewError(500, "no user found")
	}

	user := users.Data.Users[0]

	claims := Claims{
		ID:     user.ID,
		Scopes: resp.Data.Scopes,
		Login:  user.Login,
	}

	now := time.Now().UTC()

	accessClaims := claims
	accessClaims.ExpiresAt = jwt.NewNumericDate(now.Add(accessLifeTime))
	accessToken, err := createToken(accessClaims, config.JwtAccessSecret)
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(401, "cannot create JWT access token")
	}

	refreshClaims := claims
	refreshClaims.ExpiresAt = jwt.NewNumericDate(now.Add(refreshLifeTime))

	refreshToken, err := createToken(refreshClaims, config.JwtRefreshSecret)
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(401, "cannot create JWT refresh token")
	}

	err = checkUser(user.ID, resp.Data, services)
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(500, "internal error")
	}

	return &Tokens{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

type Claims struct {
	jwt.RegisteredClaims
	ID     string   `json:"id"`
	Scopes []string `json:"scopes"`
	Login  string   `json:"login"`
}

func createToken(claims jwt.Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func handleGetProfile(user model.Users, services types.Services) (*Profile, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	twitchClient, err := twitch.NewAppClient(config, tokensGrpc)
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	usersReq, err := twitchClient.GetUsers(
		&helix.UsersParams{
			IDs: []string{user.ID},
		},
	)

	if err != nil || usersReq.ErrorMessage != "" {
		logger.Error(err)
		logger.Error(usersReq.ErrorMessage)

		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if len(usersReq.Data.Users) == 0 {
		return nil, fiber.NewError(http.StatusUnauthorized, "no user found on twitch")
	}

	myTwitchUser := usersReq.Data.Users[0]

	return &Profile{
		User:       myTwitchUser,
		ApiKey:     user.ApiKey,
		IsBotAdmin: user.IsBotAdmin,
	}, nil
}

type Profile struct {
	helix.User
	ApiKey     string `json:"apiKey"`
	IsBotAdmin bool   `json:"isBotAdmin"`
}

func handleRefresh(refreshToken string) (string, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	token, err := jwt.Parse(
		refreshToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(config.JwtRefreshSecret), nil
		},
	)
	if err != nil {
		return "", fiber.NewError(401, "invalid token. Probably token is expired.")
	}

	claims := token.Claims.(jwt.MapClaims)
	scopes := []string{}
	for _, item := range claims["scopes"].([]interface{}) {
		scopes = append(scopes, item.(string))
	}

	newClaims := Claims{
		ID:     claims["id"].(string),
		Scopes: scopes,
		Login:  claims["login"].(string),
	}

	newClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(accessLifeTime))
	newToken, err := createToken(newClaims, config.JwtAccessSecret)
	if err != nil {
		logger.Error(err)
		return "", fiber.NewError(401, "cannot create new access token")
	}
	return newToken, nil
}

func handleUpdateApiKey(userId string, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	user := model.Users{}
	err := services.DB.Where(`"id" = ?`, userId).First(&user).Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	user.ApiKey = uuid.NewV4().String()
	err = services.DB.Save(&user).Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}

type Dashboard struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Avatar      string   `json:"avatar"`
	Flags       []string `json:"flags"`
}

func handleGetDashboards(user model.Users, services types.Services) ([]Dashboard, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	twitchClient, err := twitch.NewAppClient(config, tokensGrpc)
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	dashboards := []Dashboard{
		{
			ID:          user.ID,
			Name:        "",
			DisplayName: "",
			Avatar:      "",
			Flags:       []string{model.RolePermissionCanAccessDashboard.String()},
		},
	}

	if user.IsBotAdmin {
		channels := []model.Channels{}

		err = services.DB.Where(`"id" != ?`, user.ID).Find(&channels).Error
		if err != nil {
			logger.Error(err)
			return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		dashboards = append(
			dashboards, lo.Map(
				channels, func(channel model.Channels, _ int) Dashboard {
					return Dashboard{
						ID:    channel.ID,
						Flags: []string{model.RolePermissionCanAccessDashboard.String()},
					}
				},
			)...,
		)
	} else {
		if err != nil {
			logger.Error(err)
			return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		for _, role := range user.Roles {
			ok := lo.Contains(role.Role.Permissions, model.RolePermissionCanAccessDashboard.String())
			if !ok {
				continue
			}

			dashboards = append(
				dashboards, Dashboard{
					ID:    role.Role.ChannelID,
					Flags: role.Role.Permissions,
				},
			)
		}
	}

	chunks := lo.Chunk(dashboards, 100)
	wg := sync.WaitGroup{}

	for _, chunk := range chunks {
		wg.Add(1)
		go func(chunk []Dashboard) {
			defer wg.Done()

			ids := lo.Map(
				chunk, func(dashboard Dashboard, _ int) string {
					return dashboard.ID
				},
			)

			usersReq, err := twitchClient.GetUsers(
				&helix.UsersParams{
					IDs: ids,
				},
			)

			if err != nil || usersReq.ErrorMessage != "" {
				logger.Error(err)
				logger.Error(usersReq.ErrorMessage)
				return
			}

			for _, user := range usersReq.Data.Users {
				for i := range chunk {
					if chunk[i].ID == user.ID {
						chunk[i].Name = user.Login
						chunk[i].DisplayName = user.DisplayName
						chunk[i].Avatar = user.ProfileImageURL
					}
				}
			}
		}(chunk)
	}

	wg.Wait()

	return dashboards, nil
}
