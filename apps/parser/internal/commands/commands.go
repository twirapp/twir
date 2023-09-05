package commands

import (
	"context"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/cacher"
	channel_game "github.com/satont/twir/apps/parser/internal/commands/channel/game"
	channel_title "github.com/satont/twir/apps/parser/internal/commands/channel/title"
	"github.com/satont/twir/apps/parser/internal/commands/games"
	"github.com/satont/twir/apps/parser/internal/commands/manage"
	"github.com/satont/twir/apps/parser/internal/commands/nuke"
	"github.com/satont/twir/apps/parser/internal/commands/permit"
	"github.com/satont/twir/apps/parser/internal/commands/shoutout"
	"github.com/satont/twir/apps/parser/internal/commands/song"
	sr_youtube "github.com/satont/twir/apps/parser/internal/commands/songrequest/youtube"
	"github.com/satont/twir/apps/parser/internal/commands/spam"
	"github.com/satont/twir/apps/parser/internal/commands/stats"
	"github.com/satont/twir/apps/parser/internal/commands/tts"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/internal/types/services"
	"github.com/satont/twir/apps/parser/internal/variables"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/gopool"
	"github.com/satont/twir/libs/grpc/generated/parser"
)

type Commands struct {
	DefaultCommands    map[string]*types.DefaultCommand
	parseResponsesPool *gopool.Pool

	services         *services.Services
	variablesService *variables.Variables
}

type Opts struct {
	Services         *services.Services
	VariablesService *variables.Variables
}

func New(opts *Opts) *Commands {
	commands := lo.SliceToMap(
		[]*types.DefaultCommand{
			song.CurrentSong,
			channel_game.SetCommand,
			channel_game.History,
			channel_title.SetCommand,
			channel_title.History,
			manage.AddAliaseCommand,
			manage.AddCommand,
			manage.CheckAliasesCommand,
			manage.DelCommand,
			manage.EditCommand,
			manage.RemoveAliaseCommand,
			nuke.Command,
			permit.Command,
			shoutout.ShoutOut,
			spam.Command,
			stats.TopEmotes,
			stats.TopEmotesUsers,
			stats.TopMessages,
			stats.TopPoints,
			stats.TopTime,
			stats.Uptime,
			stats.UserAge,
			stats.UserFollowSince,
			stats.UserFollowage,
			stats.UserMe,
			stats.UserWatchTime,
			tts.DisableCommand,
			tts.EnableCommand,
			tts.PitchCommand,
			tts.RateCommand,
			tts.SayCommand,
			tts.SkipCommand,
			tts.VoiceCommand,
			tts.VoicesCommand,
			tts.VolumeCommand,
			sr_youtube.SkipCommand,
			sr_youtube.SrCommand,
			sr_youtube.SrListCommand,
			sr_youtube.WrongCommand,
			games.EightBall,
		}, func(v *types.DefaultCommand) (string, *types.DefaultCommand) {
			return v.Name, v
		},
	)

	ctx := &Commands{
		DefaultCommands:    commands,
		parseResponsesPool: gopool.NewPool(100),
		services:           opts.Services,
		variablesService:   opts.VariablesService,
	}

	return ctx
}

func (c *Commands) GetChannelCommands(
	ctx context.Context,
	channelId string,
) ([]*model.ChannelsCommands, error) {
	var cmds []*model.ChannelsCommands

	err := c.services.Gorm.
		Model(&model.ChannelsCommands{}).
		Where(`"channelId" = ? AND "enabled" = ?`, channelId, true).
		Preload("Responses").
		WithContext(ctx).
		Find(&cmds).Error
	if err != nil {
		return nil, err
	}

	return cmds, nil
}

var splittedNameRegexp = regexp.MustCompile(`[^\s]+`)

type FindByMessageResult struct {
	Cmd     *model.ChannelsCommands
	FoundBy string
}

// FindByMessage
// Splitting chat message by spaces, then
// read message from end to start, and delete one word from end while message gets empty,
// or we found a command in message
func (c *Commands) FindChannelCommandInInput(
	input string,
	cmds []*model.ChannelsCommands,
) *FindByMessageResult {
	msg := strings.ToLower(input)
	splittedName := splittedNameRegexp.FindAllString(msg, -1)

	res := FindByMessageResult{}

	length := len(splittedName)

	for i := 0; i < length; i++ {
		query := strings.Join(splittedName, " ")
		for _, cmd := range cmds {
			if cmd.Name == query {
				res.FoundBy = query
				res.Cmd = cmd
				break
			}

			if lo.SomeBy(
				cmd.Aliases, func(item string) bool {
					return item == query
				},
			) {
				res.FoundBy = query
				res.Cmd = cmd
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

	// sort command responses in right order, which set from dashboard ui
	if res.Cmd != nil {
		sort.Slice(
			res.Cmd.Responses, func(a, b int) bool {
				return res.Cmd.Responses[a].Order < res.Cmd.Responses[b].Order
			},
		)
	}

	return &res
}

func (c *Commands) ParseCommandResponses(
	ctx context.Context,
	command *FindByMessageResult,
	requestData *parser.ProcessCommandRequest,
) *parser.ProcessCommandResponse {
	result := &parser.ProcessCommandResponse{
		KeepOrder: &command.Cmd.KeepResponsesOrder,
		IsReply:   command.Cmd.IsReply,
	}

	var cmdParams *string
	params := strings.TrimSpace(requestData.Message.Text[len(command.FoundBy):])
	// this shit comes from 7tv for bypass message duplicate
	params = strings.ReplaceAll(params, "\U000e0000", "")
	params = strings.TrimSpace(params)
	if len(params) > 0 {
		cmdParams = &params
	}

	var defaultCommand *types.DefaultCommand

	if command.Cmd.Default {
		cmd, ok := c.DefaultCommands[command.Cmd.DefaultName.String]
		if ok {
			defaultCommand = cmd
		}
	}

	defer c.services.Gorm.
		WithContext(ctx).
		Create(
			&model.ChannelsCommandsUsages{
				ID:        uuid.New().String(),
				UserID:    requestData.Sender.Id,
				ChannelID: requestData.Channel.Id,
				CommandID: command.Cmd.ID,
			},
		)

	parseCtxChannel := &types.ParseContextChannel{
		ID:   requestData.Channel.Id,
		Name: requestData.Channel.Name,
	}
	parseCtxSender := &types.ParseContextSender{
		ID:          requestData.Sender.Id,
		Name:        requestData.Sender.Name,
		DisplayName: requestData.Sender.DisplayName,
		Badges:      requestData.Sender.Badges,
	}

	parseCtx := &types.ParseContext{
		Channel:   parseCtxChannel,
		Sender:    parseCtxSender,
		Text:      cmdParams,
		IsCommand: true,
		Services:  c.services,
		Cacher: cacher.NewCacher(
			&cacher.CacherOpts{
				Services:        c.services,
				ParseCtxChannel: parseCtxChannel,
				ParseCtxSender:  parseCtxSender,
				ParseCtxText:    cmdParams,
			},
		),
		Emotes: lo.Map(
			requestData.Message.Emotes, func(e *parser.Message_Emote, _ int) *types.ParseContextEmote {
				return &types.ParseContextEmote{
					Name:  e.Name,
					ID:    e.Id,
					Count: e.Count,
					Positions: lo.Map(
						e.Positions,
						func(p *parser.Message_EmotePosition, _ int) *types.ParseContextEmotePosition {
							return &types.ParseContextEmotePosition{
								Start: p.Start,
								End:   p.End,
							}
						},
					),
				}
			},
		),
		Command: command.Cmd,
	}

	if command.Cmd.Default && defaultCommand != nil {
		results := defaultCommand.Handler(ctx, parseCtx)

		result.Responses = lo.
			IfF(results == nil, func() []string { return []string{} }).
			ElseF(
				func() []string {
					return results.Result
				},
			)
	} else {
		result.Responses = lo.Map(
			command.Cmd.Responses, func(r *model.ChannelsCommandsResponses, _ int) string {
				return r.Text.String
			},
		)
	}

	wg := &sync.WaitGroup{}
	for i, r := range result.Responses {
		wg.Add(1)

		index := i
		response := r
		c.parseResponsesPool.Submit(
			func() {
				result.Responses[index] = c.variablesService.ParseVariablesInText(ctx, parseCtx, response)
				wg.Done()
			},
		)
	}
	wg.Wait()

	return result
}
