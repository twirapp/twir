package commands

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/cacher"
	channel_game "github.com/satont/twir/apps/parser/internal/commands/channel/game"
	channel_title "github.com/satont/twir/apps/parser/internal/commands/channel/title"
	"github.com/satont/twir/apps/parser/internal/commands/dudes"
	"github.com/satont/twir/apps/parser/internal/commands/games"
	"github.com/satont/twir/apps/parser/internal/commands/manage"
	"github.com/satont/twir/apps/parser/internal/commands/nuke"
	"github.com/satont/twir/apps/parser/internal/commands/overlays/brb"
	"github.com/satont/twir/apps/parser/internal/commands/overlays/kappagen"
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
	"github.com/twirapp/twir/libs/grpc/parser"
	"go.uber.org/zap"
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
			games.RussianRoulette,
			kappagen.Kappagen,
			brb.Start,
			brb.Stop,
			games.Duel,
			games.DuelAccept,
			games.DuelStats,
			dudes.Jump,
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
	params := strings.TrimSpace(requestData.GetMessage().GetText()[len(command.FoundBy):])
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

	go c.services.Gorm.
		WithContext(context.TODO()).
		Create(
			&model.ChannelsCommandsUsages{
				ID:        uuid.New().String(),
				UserID:    requestData.GetSender().GetId(),
				ChannelID: requestData.GetChannel().GetId(),
				CommandID: command.Cmd.ID,
			},
		)

	parseCtxChannel := &types.ParseContextChannel{
		ID:   requestData.GetChannel().GetId(),
		Name: requestData.GetChannel().GetName(),
	}
	parseCtxSender := &types.ParseContextSender{
		ID:          requestData.GetSender().GetId(),
		Name:        requestData.GetSender().GetName(),
		DisplayName: requestData.GetSender().GetDisplayName(),
		Badges:      requestData.GetSender().GetBadges(),
		Color:       requestData.GetSender().GetColor(),
	}
	mentions := make([]types.ParseContextMention, 0, len(requestData.GetMessage().GetMentions()))
	for _, m := range requestData.GetMessage().GetMentions() {
		mentions = append(
			mentions,
			types.ParseContextMention{
				UserId:    m.GetUserId(),
				UserName:  m.GetUserName(),
				UserLogin: m.GetUserLogin(),
			},
		)
	}

	parseCtx := &types.ParseContext{
		MessageId: requestData.GetMessage().GetId(),
		Channel:   parseCtxChannel,
		Sender:    parseCtxSender,
		Text:      cmdParams,
		RawText:   requestData.GetMessage().GetText(),
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
			requestData.GetMessage().GetEmotes(), func(
				e *parser.Message_Emote,
				_ int,
			) *types.ParseContextEmote {
				return &types.ParseContextEmote{
					Name:  e.GetName(),
					ID:    e.GetId(),
					Count: e.GetCount(),
					Positions: lo.Map(
						e.GetPositions(),
						func(p *parser.Message_EmotePosition, _ int) *types.ParseContextEmotePosition {
							return &types.ParseContextEmotePosition{
								Start: p.GetStart(),
								End:   p.GetEnd(),
							}
						},
					),
				}
			},
		),
		Mentions: mentions,
		Command:  command.Cmd,
	}

	if command.Cmd.Default && defaultCommand != nil {
		results, err := defaultCommand.Handler(ctx, parseCtx)
		if err != nil {
			c.services.Logger.Sugar().Error(
				"error happened on default command execution",
				zap.Error(err),
				zap.Dict(
					"channel",
					zap.String("id", requestData.Channel.Id),
					zap.String("name", requestData.Channel.Name),
				),
				zap.Dict(
					"sender",
					zap.String("id", requestData.Sender.Id),
					zap.String("name", requestData.Sender.Name),
				),
				zap.String("message", requestData.Message.Text),
				zap.Dict("command", zap.String("id", command.Cmd.ID), zap.String("name", command.Cmd.Name)),
			)

			var commandErr *types.CommandHandlerError

			if errors.As(err, &commandErr) {
				results = &types.CommandsHandlerResult{
					Result: []string{
						fmt.Sprintf("[Twir error]: %s", commandErr.Message),
					},
				}
			} else {
				results = &types.CommandsHandlerResult{
					Result: []string{"[Twir error]: unknown error happened. Please contact developers."},
				}
			}
		}

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
				defer wg.Done()
				result.Responses[index] = c.variablesService.ParseVariablesInText(ctx, parseCtx, response)
			},
		)
	}
	wg.Wait()

	return result
}
