package auth

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/libs/crypto"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/bus-core/scheduler"
)

type authBody struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type authResponse struct {
	RedirectTo string `json:"redirect_to"`
}

func (a *Auth) handleAuthPostCode(c *gin.Context) {
	ctx := c.Request.Context()

	body := authBody{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "wrong body"})
		return
	}

	if body.Code == "" {
		c.JSON(400, gin.H{"error": "no code provided"})
		return
	}

	if body.State == "" {
		c.JSON(400, gin.H{"error": "no state provided"})
		return
	}

	redirectTo, err := base64.StdEncoding.DecodeString(body.State)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("cannot decode state: %s", err)})
		return
	}

	twitchClient, err := twitch.NewAppClientWithContext(ctx, a.config, a.tokensGrpc)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("cannot create twitch client: %s", err)})
		return
	}

	tokens, err := twitchClient.RequestUserAccessToken(body.Code)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("cannot user data from twitch: %s", err)})
		return
	}

	if tokens.ErrorMessage != "" {
		c.JSON(500, gin.H{"error": tokens.ErrorMessage})
		return
	}

	twitchClient.SetUserAccessToken(tokens.Data.AccessToken)

	users, err := twitchClient.GetUsers(&helix.UsersParams{})
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Errorf("cannot get user data from twitch: %s", err)})
		return
	}
	if len(users.Data.Users) == 0 {
		c.JSON(500, gin.H{"error": "twitch user not found"})
		return
	}

	twitchUser := users.Data.Users[0]

	dbUser := &model.Users{}
	err = a.gorm.WithContext(ctx).Where("id = ?", twitchUser.ID).Preload("Token").Find(dbUser).Error
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("cannot find user: %s", err)})
		return
	}

	defaultBot := &model.Bots{}
	err = a.gorm.WithContext(ctx).Where("type = ?", "DEFAULT").Find(defaultBot).Error
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("cannot find default bot: %s", err)})
		return
	}

	if defaultBot.ID == "" {
		c.JSON(500, gin.H{"error": "no default bot found"})
		return
	}

	accessToken, err := crypto.Encrypt(tokens.Data.AccessToken, a.config.TokensCipherKey)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("cannot encrypt user access token: %s", err)})
		return
	}

	refreshToken, err := crypto.Encrypt(tokens.Data.RefreshToken, a.config.TokensCipherKey)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("cannot encrypt user refresh token: %s", err)})
		return
	}

	if dbUser.ID == "" {
		newUser := &model.Users{
			ID:         twitchUser.ID,
			IsTester:   false,
			IsBotAdmin: false,
			ApiKey:     uuid.NewString(),
			Channel: &model.Channels{
				ID:    twitchUser.ID,
				BotID: defaultBot.ID,
			},
		}

		if err := a.gorm.Create(newUser).Error; err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("cannot create user: %s", err)})
			return
		}

		dbUser = newUser
	}

	tokenData := model.Tokens{
		ID:                  uuid.New().String(),
		AccessToken:         accessToken,
		RefreshToken:        refreshToken,
		ExpiresIn:           int32(tokens.Data.ExpiresIn),
		ObtainmentTimestamp: time.Now().UTC(),
		Scopes:              tokens.Data.Scopes,
	}
	if dbUser.TokenID.Valid {
		tokenData.ID = dbUser.TokenID.String
	}

	if err := a.gorm.WithContext(ctx).Save(tokenData).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("cannot update user token: %s", err)})
		return
	}

	if err := a.gorm.WithContext(ctx).Debug().Save(&tokenData).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("cannot update db user: %s", err)})
		return
	}

	dbUser.TokenID = sql.NullString{
		String: tokenData.ID,
		Valid:  true,
	}

	if dbUser.Channel == nil || dbUser.Channel.ID == "" {
		dbUser.Channel = &model.Channels{
			ID:    twitchUser.ID,
			BotID: defaultBot.ID,
		}
	}

	if err := a.gorm.WithContext(ctx).Debug().Save(dbUser).Error; err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("cannot update db user: %s", err)})
		return
	}

	err = a.bus.Scheduler.CreateDefaultRoles.Publish(
		scheduler.CreateDefaultRolesRequest{ChannelsIDs: []string{twitchUser.ID}},
	)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("cannot create default roles: %s", err)})
		return
	}

	err = a.bus.Scheduler.CreateDefaultCommands.Publish(
		scheduler.CreateDefaultCommandsRequest{ChannelsIDs: []string{twitchUser.ID}},
	)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("cannot create default commands: %s", err)})
		return
	}

	a.sessions.Put(ctx, "dbUser", &dbUser)
	a.sessions.Put(ctx, "twitchUser", &twitchUser)
	a.sessions.Put(ctx, "dashboardId", dbUser.ID)

	if err := a.bus.EventSub.SubscribeToAllEvents.Publish(
		eventsub.EventsubSubscribeToAllEventsRequest{
			ChannelID: dbUser.ID,
		},
	); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("cannot subscribe to eventsub: %s", err)})
		return
	}

	c.JSON(
		200,
		&authResponse{
			RedirectTo: string(redirectTo),
		},
	)
}
