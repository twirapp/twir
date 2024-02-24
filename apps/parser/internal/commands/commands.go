package commands

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
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
	busparser "github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"go.uber.org/zap"
)

type Commands struct {
	DefaultCommands  map[string]*types.DefaultCommand
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
		DefaultCommands:  commands,
		services:         opts.Services,
		variablesService: opts.VariablesService,
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
	requestData twitch.TwitchChatMessage,
) *busparser.CommandParseResponse {
	result := &busparser.CommandParseResponse{
		KeepOrder: command.Cmd.KeepResponsesOrder,
		IsReply:   command.Cmd.IsReply,
	}

	var cmdParams *string
	params := strings.TrimSpace(requestData.Message.Text[1:][len(command.FoundBy):])
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
				UserID:    requestData.ChatterUserId,
				ChannelID: requestData.BroadcasterUserId,
				CommandID: command.Cmd.ID,
			},
		)

	parseCtxChannel := &types.ParseContextChannel{
		ID:   requestData.BroadcasterUserId,
		Name: requestData.BroadcasterUserLogin,
	}
	parseCtxSender := &types.ParseContextSender{
		ID:          requestData.ChatterUserId,
		Name:        requestData.ChatterUserLogin,
		DisplayName: requestData.ChatterUserName,
		Badges:      []string{},
		Color:       requestData.Color,
	}

	mentions := make(
		[]twitch.ChatMessageMessageFragmentMention,
		0,
		len(requestData.Message.Fragments),
	)
	for _, f := range requestData.Message.Fragments {
		if f.Type != twitch.FragmentType_MENTION {
			continue
		}
		mentions = append(mentions, *f.Mention)
	}

	emotes := make([]*types.ParseContextEmote, 0, len(requestData.Message.Fragments))
	for _, f := range requestData.Message.Fragments {
		if f.Type != twitch.FragmentType_EMOTE {
			continue
		}
		emotes = append(
			emotes, &types.ParseContextEmote{
				Name:  f.Text,
				ID:    f.Emote.Id,
				Count: 1,
				Positions: []*types.ParseContextEmotePosition{
					{
						Start: int64(f.Position.Start),
						End:   int64(f.Position.End),
					},
				},
			},
		)
	}

	parseCtx := &types.ParseContext{
		MessageId: requestData.MessageId,
		Channel:   parseCtxChannel,
		Sender:    parseCtxSender,
		Text:      cmdParams,
		RawText:   requestData.Message.Text[1:],
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
		Emotes:   emotes,
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
					zap.String("id", requestData.BroadcasterUserId),
					zap.String("name", requestData.BroadcasterUserLogin),
				),
				zap.Dict(
					"sender",
					zap.String("id", requestData.ChatterUserId),
					zap.String("name", requestData.ChatterUserLogin),
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
		go func() {
			defer wg.Done()
			result.Responses[index] = c.variablesService.ParseVariablesInText(ctx, parseCtx, response)
		}()
	}
	wg.Wait()

	return result
}

func (c *Commands) ProcessChatMessage(ctx context.Context, data twitch.TwitchChatMessage) (
	*busparser.CommandParseResponse,
	error,
) {
	if data.Message.Text[0] != '!' {
		return nil, nil
	}

	cmds, err := c.GetChannelCommands(ctx, data.BroadcasterUserId)
	if err != nil {
		return nil, err
	}

	cmd := c.FindChannelCommandInInput(data.Message.Text[1:], cmds)
	if cmd.Cmd == nil {
		return nil, nil
	}

	if cmd.Cmd.OnlineOnly {
		stream := &model.ChannelsStreams{}
		err = c.services.Gorm.
			WithContext(ctx).
			Where(`"userId" = ?`, data.BroadcasterUserId).
			Find(stream).Error
		if err != nil {
			return nil, err
		}
		if stream == nil || stream.ID == "" {
			return nil, nil
		}
	}

	convertedBadges := make([]string, 0, len(data.Badges))
	for _, badge := range data.Badges {
		convertedBadges = append(convertedBadges, badge.Id)
	}

	dbUser, _, userRoles, commandRoles, err := c.prepareCooldownAndPermissionsCheck(
		ctx,
		data.ChatterUserId,
		data.BroadcasterUserId,
		convertedBadges,
		cmd.Cmd,
	)
	if err != nil {
		return nil, err
	}

	shouldCheckCooldown := c.shouldCheckCooldown(convertedBadges, cmd.Cmd, userRoles)
	if cmd.Cmd.CooldownType == "GLOBAL" && cmd.Cmd.Cooldown.Int64 > 0 && shouldCheckCooldown {
		key := fmt.Sprintf("commands:%s:cooldowns:global", cmd.Cmd.ID)
		rErr := c.services.Redis.Get(ctx, key).Err()

		if errors.Is(rErr, redis.Nil) {
			c.services.Redis.Set(ctx, key, "", time.Duration(cmd.Cmd.Cooldown.Int64)*time.Second)
		} else if rErr != nil {
			c.services.Logger.Sugar().Error(rErr)
			return nil, errors.New("error while setting redis cooldown for command")
		} else {
			return nil, nil
		}
	}

	if cmd.Cmd.CooldownType == "PER_USER" && cmd.Cmd.Cooldown.Int64 > 0 && shouldCheckCooldown {
		key := fmt.Sprintf("commands:%s:cooldowns:user:%s", cmd.Cmd.ID, data.ChatterUserId)
		rErr := c.services.Redis.Get(ctx, key).Err()

		if rErr == redis.Nil {
			c.services.Redis.Set(ctx, key, "", time.Duration(cmd.Cmd.Cooldown.Int64)*time.Second)
		} else if rErr != nil {
			zap.S().Error(rErr)
			return nil, errors.New("error while setting redis cooldown for command")
		} else {
			return nil, nil
		}
	}

	hasPerm := c.isUserHasPermissionToCommand(
		data.ChatterUserId,
		data.BroadcasterUserId,
		cmd.Cmd,
		dbUser,
		userRoles,
		commandRoles,
	)

	if !hasPerm {
		return nil, nil
	}

	go func() {
		gCtx := context.Background()

		c.services.GrpcClients.Events.CommandUsed(
			// this should be background, because we don't want to wait for response
			gCtx,
			&events.CommandUsedMessage{
				BaseInfo:           &events.BaseInfo{ChannelId: data.BroadcasterUserId},
				CommandId:          cmd.Cmd.ID,
				CommandName:        cmd.Cmd.Name,
				CommandInput:       strings.TrimSpace(data.Message.Text[len(cmd.FoundBy):]),
				UserName:           data.ChatterUserLogin,
				UserDisplayName:    data.ChatterUserName,
				UserId:             data.ChatterUserId,
				IsDefault:          cmd.Cmd.Default,
				DefaultCommandName: cmd.Cmd.DefaultName.String,
			},
		)

		alert := model.ChannelAlert{}
		if err := c.services.Gorm.Where(
			"channel_id = ? AND command_ids && ?",
			data.BroadcasterUserId,
			pq.StringArray{cmd.Cmd.ID},
		).Find(&alert).Error; err != nil {
			zap.S().Error(err)
			return
		}

		if alert.ID == "" {
			return
		}
		c.services.GrpcClients.WebSockets.TriggerAlert(
			gCtx,
			&websockets.TriggerAlertRequest{
				ChannelId: data.BroadcasterUserId,
				AlertId:   alert.ID,
			},
		)
	}()

	// TODO: refactor parsectx to new chat message struct
	result := c.ParseCommandResponses(
		ctx,
		cmd,
		data,
	)

	return result, nil
}
