package v1_handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api-new/internal/http/helpers"
	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func (c *Handlers) getChannelCommands(channelId string) ([]model.ChannelsCommands, error) {
	var commands []model.ChannelsCommands
	err := c.gorm.
		Where(`"channelId" = ?`, channelId).
		Preload("Responses").
		Preload("Group").
		Find(&commands).Error

	return commands, err
}

func (c *Handlers) GetChannelsCommands(ctx *fiber.Ctx) error {
	channelId := ctx.Params("channelId")

	commands, err := c.getChannelCommands(channelId)
	if err != nil {
		c.logger.Error(err)
		return helpers.ErrInternalError
	}

	return ctx.JSON(commands)
}

type responsesDto struct {
	Text  string `validate:"min=1,max=1000" json:"text"`
	Order uint8  `validate:"gte=0"                                             json:"order"`
}

type commandDto struct {
	ID *string `json:"id,omitempty"`

	Name               string         `validate:"required,min=1,max=100"       json:"name"`
	Cooldown           uint32         `validate:"lte=86400"                   json:"cooldown"`
	CooldownType       string         `validate:"required"                    json:"cooldownType"`
	Description        *string        `validate:"omitempty,max=500"           json:"description,omitempty"`
	Aliases            []string       `validate:"required,dive,max=20" json:"aliases"`
	Visible            *bool          `validate:"omitempty,required"          json:"visible,omitempty"`
	Enabled            *bool          `validate:"omitempty,required"          json:"enabled,omitempty"`
	Responses          []responsesDto `validate:"dive"                        json:"responses"`
	KeepResponsesOrder *bool          `validate:"required"                    json:"keepResponsesOrder"`
	IsReply            *bool          `validate:"omitempty,required"          json:"isReply,omitempty"`
	GroupID            *string        `json:"groupId,omitempty"`
	DeniedUsersIds     []string       `json:"deniedUsersIds"`
	AllowedUsersIds    []string       `json:"allowedUsersIds"`
	RolesIDS           []string       `json:"rolesIds"`
	OnlineOnly         *bool          `validate:"required" json:"onlineOnly"`
}

func (c *Handlers) checkIsCommandWithNameExists(channelId string, incomingNames []string, exceptCommand *string) (bool, error) {
	channelCommands, err := c.getChannelCommands(channelId)
	if err != nil {
		return false, nil
	}

	busyNames := make([]string, 0, len(channelCommands))
	for _, command := range channelCommands {
		if exceptCommand != nil && command.ID == *exceptCommand {
			continue
		}

		busyNames = append(busyNames, command.Name)
		busyNames = append(busyNames, command.Aliases...)
	}

	return lo.SomeBy(busyNames, func(item string) bool {
		return lo.Contains(incomingNames, item)
	}), nil
}

func (c *Handlers) CreateCommand(ctx *fiber.Ctx) error {
	channelId := ctx.Params("channelId")

	data := &commandDto{}
	err := c.middlewares.ValidateBody(ctx, data)
	if err != nil {
		return err
	}

	data.Name = strings.TrimSpace(data.Name)
	data.Name = strings.ToLower(data.Name)
	data.Name = strings.Replace(data.Name, "!", "", 1)
	if len(data.Name) == 0 {
		return fiber.NewError(http.StatusBadRequest, "name cannot be empty")
	}

	data.Aliases = lo.Filter(lo.Map(data.Aliases, func(a string, _ int) string {
		a = strings.TrimSpace(a)
		a = strings.ToLower(a)
		a = strings.Replace(a, "!", "", 1)
		return a
	}), func(item string, _ int) bool {
		return len(item) > 0
	})

	isExists, err := c.checkIsCommandWithNameExists(
		channelId,
		append(data.Aliases, data.Name),
		nil,
	)
	if err != nil {
		c.logger.Error(err)
		return helpers.ErrInternalError
	}

	if isExists {
		return fiber.NewError(http.StatusBadRequest, "command with that name or alias already exists")
	}

	err = c.gorm.Transaction(func(tx *gorm.DB) error {
		command := &model.ChannelsCommands{
			ID:           uuid.New().String(),
			Name:         data.Name,
			Cooldown:     null.IntFrom(int64(data.Cooldown)),
			CooldownType: data.CooldownType,
			Enabled:      lo.If(data.Enabled == nil, false).Else(*data.Enabled),
			Aliases:      data.Aliases,
			Description:  null.StringFromPtr(data.Description),
			Visible:      lo.If(data.Visible == nil, false).Else(*data.Visible),
			ChannelID:    channelId,
			Module:       "CUSTOM",
			IsReply:      lo.If(data.IsReply == nil, false).Else(*data.IsReply),
			KeepResponsesOrder: lo.If(data.KeepResponsesOrder == nil, false).
				Else(*data.KeepResponsesOrder),
			GroupID:         null.StringFromPtr(data.GroupID),
			AllowedUsersIDS: data.AllowedUsersIds,
			DeniedUsersIDS:  data.DeniedUsersIds,
			RolesIDS:        data.RolesIDS,
			OnlineOnly:      *data.OnlineOnly,
		}

		err = tx.Save(command).Error
		if err != nil {
			return err
		}

		for _, r := range data.Responses {
			response := model.ChannelsCommandsResponses{
				ID:        uuid.New().String(),
				Text:      null.NewString(r.Text, true),
				Order:     int(r.Order),
				CommandID: command.ID,
			}
			err = tx.Save(&response).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.logger.Error(err)
		return helpers.ErrInternalError
	}

	c.cacheStorage.DeleteGet(c.cacheStorage.BuildKey(fmt.Sprintf("v1/channels/%s/commands", channelId)))

	return ctx.SendStatus(http.StatusCreated)
}
