package commands

import (
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/commands/shoutout"
	"github.com/satont/tsuwari/apps/parser/internal/commands/tts"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/satont/tsuwari/apps/parser/internal/commands/dota"
	"github.com/satont/tsuwari/apps/parser/internal/commands/manage"
	"github.com/satont/tsuwari/apps/parser/internal/commands/nuke"
	"github.com/satont/tsuwari/apps/parser/internal/commands/permit"
	"github.com/satont/tsuwari/apps/parser/internal/commands/songrequest/youtube"
	"github.com/satont/tsuwari/apps/parser/internal/commands/spam"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	"github.com/satont/tsuwari/apps/parser/internal/variables"
	"github.com/satont/tsuwari/apps/parser/pkg/helpers"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"

	uuid "github.com/satori/go.uuid"

	"github.com/satont/tsuwari/apps/parser/internal/commands/channel/game"
	"github.com/satont/tsuwari/apps/parser/internal/commands/channel/title"

	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"

	"gorm.io/gorm"
)

type Commands struct {
	DefaultCommands []types.DefaultCommand
}

func New() Commands {
	commands := []types.DefaultCommand{
		channel_title.SetCommand,
		channel_title.History,
		channel_game.SetCommand,
		channel_game.History,
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
		sr_youtube.SrCommand,
		sr_youtube.WrongCommand,
		sr_youtube.SrListCommand,
		sr_youtube.SkipCommand,
		shoutout.ShoutOut,
		tts.SayCommand,
		tts.SkipCommand,
		tts.VoicesCommand,
		tts.VoiceCommand,
		tts.RateCommand,
		tts.PitchCommand,
		tts.VolumeCommand,
	}

	ctx := Commands{
		DefaultCommands: commands,
	}

	return ctx
}

func (c *Commands) GetChannelCommands(channelId string) (*[]model.ChannelsCommands, error) {
	db := do.MustInvoke[gorm.DB](di.Provider)

	cmds := []model.ChannelsCommands{}

	err := db.
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
	db := do.MustInvoke[gorm.DB](di.Provider)
	variablesService := do.MustInvoke[variables.Variables](di.Provider)

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

	db.Create(&model.ChannelsCommandsUsages{
		ID:        uuid.NewV4().String(),
		UserID:    data.Sender.Id,
		ChannelID: data.Channel.Id,
		CommandID: cmd.ID,
	})

	if cmd.Default && isDefaultExists {
		results := defaultCommand.Handler(variables_cache.ExecutionContext{
			ChannelName:       data.Channel.Name,
			ChannelId:         data.Channel.Id,
			SenderId:          data.Sender.Id,
			SenderName:        data.Sender.Name,
			SenderDisplayName: data.Sender.DisplayName,
			SenderBadges:      data.Sender.Badges,
			Text:              cmdParams,
			IsCommand:         true,
			Command:           command.Cmd,
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
			Text:              cmdParams,
			SenderId:          data.Sender.Id,
			SenderName:        &data.Sender.DisplayName,
			SenderDisplayName: &data.Sender.DisplayName,
			ChannelId:         data.Channel.Id,
			IsCommand:         true,
			Command:           command.Cmd,
			SenderBadges:      data.Sender.Badges,
		})

		go func(i int, r string) {
			defer wg.Done()

			result.Responses[i] = variablesService.ParseInput(cacheService, r)
		}(i, r)
	}
	wg.Wait()

	return result
}
