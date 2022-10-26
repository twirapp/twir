package auth

import (
	"fmt"
	"sync"
	"time"
	model "tsuwari/models"
	"tsuwari/twitch"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
	uuid "github.com/satori/go.uuid"
)

type Tokens struct {
	UserId       string `json:"userId"`
	RefreshToken string `json:"refreshToken,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
}

func handleGetToken(code string, services types.Services) (*Tokens, error) {
	resp, err := services.Twitch.Client.RequestUserAccessToken(code)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(401, "cannot get user tokens")
	}

	newClient := twitch.NewClient(&helix.Options{
		ClientID:         services.Cfg.TwitchClientId,
		ClientSecret:     services.Cfg.TwitchClientSecret,
		UserAccessToken:  resp.Data.AccessToken,
		UserRefreshToken: resp.Data.RefreshToken,
	})

	users, err := newClient.Client.GetUsers(&helix.UsersParams{})
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(401, "cannot get user tokens")
	}

	user := users.Data.Users[0]

	claims := Claims{
		ID:     user.ID,
		Scopes: resp.Data.Scopes,
		Login:  user.Login,
	}

	now := time.Now()

	accessClaims := claims
	accessClaims.ExpiresAt = jwt.NewNumericDate(now.Add(10 * time.Minute))
	accessToken, err := createToken(claims, services.Cfg.JwtAccessSecret)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(401, "cannot create JWT access token")
	}

	refreshClaims := claims
	refreshClaims.ExpiresAt = jwt.NewNumericDate(now.Add(31 * 24 * time.Hour))

	refreshToken, err := createToken(claims, services.Cfg.JwtRefreshSecret)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(401, "cannot create JWT refresh token")
	}

	return &Tokens{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

type Claims struct {
	ID     string   `json:"id"`
	Scopes []string `json:"scopes"`
	Login  string   `json:"login"`
	jwt.RegisteredClaims
}

func createToken(claims Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func handleGetProfile(user model.Users, services types.Services) (*Profile, error) {
	dbDashboards := []model.ChannelsDashboardAccess{}
	err := services.DB.Where(`"userId" = ?`, user.ID).Find(&dbDashboards).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "internal error")
	}

	if user.IsBotAdmin {
		channelsIds := lo.Map(dbDashboards, func(d model.ChannelsDashboardAccess, _ int) string {
			return d.ChannelID
		})
		channelsIds = append(channelsIds, user.ID)
		channels := []model.Channels{}
		err := services.DB.Not(channelsIds).Find(&channels).Error
		if err != nil {
			services.Logger.Sugar().Error(err)
			return nil, fiber.NewError(500, "cannot get channels")
		}

		for _, c := range channels {
			dbDashboards = append(dbDashboards, model.ChannelsDashboardAccess{
				ID:        uuid.NewV4().String(),
				ChannelID: c.ID,
				UserID:    user.ID,
			})
		}
	}

	neededUsersIds := []string{user.ID}
	neededUsersIds = append(
		neededUsersIds,
		lo.Map(dbDashboards, func(a model.ChannelsDashboardAccess, _ int) string {
			return a.ChannelID
		})...,
	)
	chunks := lo.Chunk(neededUsersIds, 100)

	twitchUsers := []helix.User{}
	twitchUsersWg := sync.WaitGroup{}
	twitchUsersWg.Add(len(chunks))

	for _, chunk := range chunks {
		go func(c []string) {
			defer twitchUsersWg.Done()
			users, err := services.Twitch.Client.GetUsers(&helix.UsersParams{IDs: c})
			if err != nil {
				return
			}

			twitchUsers = append(twitchUsers, users.Data.Users...)
		}(chunk)
	}

	twitchUsersWg.Wait()

	dashboards := []DashboardAndUser{}

	for _, d := range dbDashboards {
		twitchUser, _ := lo.Find(twitchUsers, func(user helix.User) bool {
			return user.ID == d.ChannelID
		})

		newDash := DashboardAndUser{
			ChannelsDashboardAccess: d,
			TwitchUser:              twitchUser,
		}

		dashboards = append(dashboards, newDash)
	}

	myTwitchUser, _ := lo.Find(twitchUsers, func(u helix.User) bool {
		return u.ID == user.ID
	})

	return &Profile{
		User:       myTwitchUser,
		DashBoards: dashboards,
	}, nil
}

type DashboardAndUser struct {
	model.ChannelsDashboardAccess
	TwitchUser helix.User `json:"twitchUser"`
}

type Profile struct {
	helix.User
	DashBoards []DashboardAndUser `json:"dashboards"`
}

func handleRefresh(dto *refreshDto, services types.Services) (string, error) {
	token, err := jwt.Parse(dto.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(services.Cfg.JwtAccessSecret), nil
	})
	if err != nil {
		return "", fiber.NewError(401, "invalid token. Probably token is expired.")
	}

	claims := token.Claims.(jwt.MapClaims)
	newClaims := Claims{
		ID:     claims["id"].(string),
		Scopes: claims["scopes"].([]string),
		Login:  claims["login"].(string),
	}
	newClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(31 * 24 * time.Hour))
	newToken, err := createToken(newClaims, services.Cfg.JwtRefreshSecret)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return "", fiber.NewError(401, "cannot create new access token")
	}
	return newToken, nil
}
