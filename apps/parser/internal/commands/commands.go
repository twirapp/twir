package commands

import (
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/satont/tsuwari/apps/parser/internal/commands/dota"
	"github.com/satont/tsuwari/apps/parser/internal/commands/manage"
	"github.com/satont/tsuwari/apps/parser/internal/commands/nuke"
	"github.com/satont/tsuwari/apps/parser/internal/commands/permit"
	"github.com/satont/tsuwari/apps/parser/internal/commands/spam"
	"github.com/satont/tsuwari/apps/parser/internal/config/twitch"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	"github.com/satont/tsuwari/apps/parser/internal/variables"
	"github.com/satont/tsuwari/apps/parser/pkg/helpers"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/eval"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"

	uuid "github.com/satori/go.uuid"

	channel_game "github.com/satont/tsuwari/apps/parser/internal/commands/channel/game"
	channel_title "github.com/satont/tsuwari/apps/parser/internal/commands/channel/title"

	usersauth "github.com/satont/tsuwari/apps/parser/internal/twitch/user"

	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"

	"github.com/go-redis/redis/v9"
	dotaGrpc "github.com/satont/tsuwari/libs/grpc/generated/dota"
	"gorm.io/gorm"
)

type Commands struct {
	DefaultCommands  []types.DefaultCommand
	redis            *redis.Client
	variablesService variables.Variables
	Db               *gorm.DB
	UsersAuth        *usersauth.UsersTokensService
	Twitch           *twitch.Twitch
	DotaGrpc         dotaGrpc.DotaClient
	BotsGrpc         bots.BotsClient
	EvalGrpc         eval.EvalClient
}

type CommandsOpts struct {
	Redis            *redis.Client
	VariablesService variables.Variables
	Db               *gorm.DB
	UsersAuth        *usersauth.UsersTokensService
	Twitch           *twitch.Twitch
	DotaGrpc         dotaGrpc.DotaClient
	BotsGrpc         bots.BotsClient
	EvalGrpc         eval.EvalClient
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
		manage.AddAliaseCommand,
		manage.RemoveAliaseCommand,
		manage.CheckAliasesCommand,
	}

	ctx := Commands{
		redis:            opts.Redis,
		DefaultCommands:  commands,
		variablesService: opts.VariablesService,
		Db:               opts.Db,
		UsersAuth:        opts.UsersAuth,
		Twitch:           opts.Twitch,
		BotsGrpc:         opts.BotsGrpc,
		DotaGrpc:         opts.DotaGrpc,
		EvalGrpc:         opts.EvalGrpc,
	}

	return ctx
}

func (c *Commands) GetChannelCommands(channelId string) (*[]model.ChannelsCommands, error) {
	cmds := []model.ChannelsCommands{}

	err := c.Db.
		Model(&model.ChannelsCommands{}).
		Where(`"channelId" = ? AND "enabled" = ?`, channelId, true).
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
			return res.Cmd.Responses[a].Order < res.Cmd.Responses[b].Order
		})
	}

	return res
}

func (c *Commands) ParseCommandResponses(
	command FindByMessageResult,
	data *parser.ProcessCommandRequest,
) *parser.ProcessCommandResponse {
	result := &parser.ProcessCommandResponse{
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
				DotaGrpc:  c.DotaGrpc,
				BotsGrpc:  c.BotsGrpc,
				EvalGrpc:  c.EvalGrpc,
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
			BotsGrpc:   c.BotsGrpc,
			DotaGrpc:   c.DotaGrpc,
			EvalGrpc:   c.EvalGrpc,
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
