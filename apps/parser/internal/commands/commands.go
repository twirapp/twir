package commands

import (
	"regexp"
	"sort"
	"strings"
	"sync"
	model "tsuwari/models"
	channel_game "tsuwari/parser/internal/commands/channel/game"
	channel_title "tsuwari/parser/internal/commands/channel/title"
	"tsuwari/parser/internal/commands/dota"
	"tsuwari/parser/internal/commands/manage"
	"tsuwari/parser/internal/commands/nuke"
	"tsuwari/parser/internal/commands/permit"
	sr_youtube "tsuwari/parser/internal/commands/songrequest/youtube"
	"tsuwari/parser/internal/commands/spam"
	"tsuwari/parser/internal/config/twitch"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables"
	"tsuwari/parser/pkg/helpers"

	usersauth "tsuwari/parser/internal/twitch/user"

	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/go-redis/redis/v9"
	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
	parserproto "github.com/satont/tsuwari/libs/nats/parser"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Commands struct {
	DefaultCommands  []types.DefaultCommand
	redis            *redis.Client
	variablesService variables.Variables
	Db               *gorm.DB
	UsersAuth        *usersauth.UsersTokensService
	Nats             *nats.Conn
	Twitch           *twitch.Twitch
}

type CommandsOpts struct {
	Redis            *redis.Client
	VariablesService variables.Variables
	Db               *gorm.DB
	UsersAuth        *usersauth.UsersTokensService
	Nats             *nats.Conn
	Twitch           *twitch.Twitch
}

func New(opts CommandsOpts) Commands {
	commands := []types.DefaultCommand{
		channel_title.Command,
		channel_game.Command,
		permit.Command,
		spam.Command,
		nuke.Command,
		dota.AddAccCommand,
		dota.DelAccCommand,
		dota.ListAccCommand,
		dota.NpAccCommand,
		dota.WlCommand,
		dota.LgCommand,
		dota.GmCommand,
		manage.AddCommand,
		manage.DelCommand,
		manage.EditCommand,
		sr_youtube.SrCommand,
	}

	ctx := Commands{
		redis:            opts.Redis,
		DefaultCommands:  commands,
		variablesService: opts.VariablesService,
		Db:               opts.Db,
		UsersAuth:        opts.UsersAuth,
		Nats:             opts.Nats,
		Twitch:           opts.Twitch,
	}

	return ctx
}

func (c *Commands) GetChannelCommands(channelId string) (*[]model.ChannelsCommands, error) {
	cmds := []model.ChannelsCommands{}
	err := c.Db.
		Model(&model.ChannelsCommands{}).
		Where(`"channelId" = ?`, channelId).
		Preload("Responses").
		Find(&cmds).Error
	if err != nil {
		return nil, err
	}

	return &cmds, nil
}

var splittedNameRegexp = regexp.MustCompile(`[^\s]+`)

type FindByMessageResult struct {
	Cmd     *model.ChannelsCommands
	FoundBy string
}

func (c *Commands) FindByMessage(input string, cmds *[]model.ChannelsCommands) FindByMessageResult {
	msg := strings.ToLower(input)
	splittedName := splittedNameRegexp.FindAllString(msg, -1)

	res := FindByMessageResult{}

	length := len(splittedName)

	for i := 0; i < length; i++ {
		query := strings.Join(splittedName, " ")
		for _, c := range *cmds {
			if c.Name == query {
				res.FoundBy = query
				res.Cmd = &c
				break
			}

			if helpers.Contains(c.Aliases, query) {
				res.FoundBy = query
				res.Cmd = &c
				break
			}
		}

		if res.Cmd != nil {
			break
		} else {
			splittedName = splittedName[:len(splittedName)-1]
			continue
		}
	}

	if res.Cmd != nil {
		sort.Slice(res.Cmd.Responses, func(a, b int) bool {
			return res.Cmd.Responses[a].Order > res.Cmd.Responses[b].Order
		})
	}

	return res
}

func (c *Commands) ParseCommandResponses(
	command FindByMessageResult,
	data parserproto.Request,
) *parserproto.Response {
	result := &parserproto.Response{
		KeepOrder: &command.Cmd.KeepResponsesOrder,
	}

	cmd := *command.Cmd
	var cmdParams *string
	params := strings.TrimSpace(data.Message.Text[len(command.FoundBy):])
	if len(params) > 0 {
		cmdParams = &params
	}

	defaultCommand, isDefaultExists := lo.Find(
		c.DefaultCommands,
		func(command types.DefaultCommand) bool {
			if cmd.DefaultName.Valid {
				return command.Name == cmd.DefaultName.String
			} else {
				return false
			}
		},
	)

	c.Db.Create(&model.ChannelsCommandsUsages{
		ID:        uuid.NewV4().String(),
		UserID:    data.Sender.Id,
		ChannelID: data.Channel.Id,
		CommandID: cmd.ID,
	})

	if cmd.Default && isDefaultExists {
		results := defaultCommand.Handler(variables_cache.ExecutionContext{
			ChannelName: data.Channel.Name,
			ChannelId:   data.Channel.Id,
			SenderId:    data.Sender.Id,
			SenderName:  data.Sender.Name,
			Text:        cmdParams,
			Services: variables_cache.ExecutionServices{
				Redis:     c.redis,
				Regexp:    nil,
				Twitch:    c.Twitch,
				Db:        c.Db,
				UsersAuth: c.UsersAuth,
				Nats:      c.Nats,
			},
			IsCommand: true,
			Command:   command.Cmd,
		})
		if results == nil {
			result.Responses = []string{}
		} else {
			result.Responses = results.Result
		}
	} else {
		result.Responses = lo.Map(cmd.Responses, func(r model.ChannelsCommandsResponses, _ int) string {
			if r.Text.Valid {
				return r.Text.String
			} else {
				return ""
			}
		})
	}

	result.IsReply = cmd.IsReply

	wg := sync.WaitGroup{}
	for i, r := range result.Responses {
		wg.Add(1)
		// TODO: concatenate all responses into one slice and use it for cache
		cacheService := variables_cache.New(variables_cache.VariablesCacheOpts{
			Text:       cmdParams,
			SenderId:   data.Sender.Id,
			SenderName: &data.Sender.DisplayName,
			ChannelId:  data.Channel.Id,
			Redis:      c.redis,
			Regexp:     variables.Regexp,
			Twitch:     c.Twitch,
			DB:         c.Db,
			Nats:       c.Nats,
			IsCommand:  true,
			Command:    command.Cmd,
		})

		go func(i int, r string) {
			defer wg.Done()

			result.Responses[i] = c.variablesService.ParseInput(cacheService, r)
		}(i, r)
	}
	wg.Wait()

	return result
}
