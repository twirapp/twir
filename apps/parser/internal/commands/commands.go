package commands

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/cacher"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	seventv "github.com/twirapp/twir/apps/parser/internal/commands/7tv"
	"github.com/twirapp/twir/apps/parser/internal/commands/categories_aliases"
	channel_game "github.com/twirapp/twir/apps/parser/internal/commands/channel/game"
	channel_title "github.com/twirapp/twir/apps/parser/internal/commands/channel/title"
	"github.com/twirapp/twir/apps/parser/internal/commands/chat_wall"
	"github.com/twirapp/twir/apps/parser/internal/commands/clip"
	"github.com/twirapp/twir/apps/parser/internal/commands/dudes"
	"github.com/twirapp/twir/apps/parser/internal/commands/games"
	"github.com/twirapp/twir/apps/parser/internal/commands/manage"
	"github.com/twirapp/twir/apps/parser/internal/commands/marker"
	"github.com/twirapp/twir/apps/parser/internal/commands/nuke"
	"github.com/twirapp/twir/apps/parser/internal/commands/overlays/brb"
	"github.com/twirapp/twir/apps/parser/internal/commands/overlays/kappagen"
	"github.com/twirapp/twir/apps/parser/internal/commands/permit"
	"github.com/twirapp/twir/apps/parser/internal/commands/predictions"
	"github.com/twirapp/twir/apps/parser/internal/commands/prefix"
	"github.com/twirapp/twir/apps/parser/internal/commands/shorturl"
	"github.com/twirapp/twir/apps/parser/internal/commands/shoutout"
	"github.com/twirapp/twir/apps/parser/internal/commands/song"
	sr_youtube "github.com/twirapp/twir/apps/parser/internal/commands/songrequest/youtube"
	"github.com/twirapp/twir/apps/parser/internal/commands/spam"
	"github.com/twirapp/twir/apps/parser/internal/commands/stats"
	"github.com/twirapp/twir/apps/parser/internal/commands/subage"
	"github.com/twirapp/twir/apps/parser/internal/commands/tracer"
	"github.com/twirapp/twir/apps/parser/internal/commands/tts"
	"github.com/twirapp/twir/apps/parser/internal/commands/utility"
	"github.com/twirapp/twir/apps/parser/internal/commands/vips"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/types/services"
	"github.com/twirapp/twir/apps/parser/internal/variables"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/bus-core/events"
	busparser "github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/i18n"
	channelscommandsusages "github.com/twirapp/twir/libs/repositories/channels_commands_usages"
	commandsmodel "github.com/twirapp/twir/libs/repositories/commands/model"
	commandresponsemodel "github.com/twirapp/twir/libs/repositories/commands_response/model"
	commandswithgroupsandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
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
			song.Playlist,
			song.History,
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
			kappagen.Kappagen,
			brb.Start,
			brb.Stop,
			games.EightBall,
			games.RussianRoulette,
			games.Voteban,
			games.Duel,
			games.DuelAccept,
			games.DuelStats,
			games.Seppuku,
			dudes.Jump,
			dudes.Grow,
			dudes.Color,
			dudes.Sprite,
			dudes.Leave,
			dudes.List,
			seventv.Profile,
			seventv.EmoteFind,
			seventv.EmoteRename,
			seventv.EmoteDelete,
			seventv.EmoteAdd,
			seventv.EmoteCopy,
			clip.MakeClip,
			marker.Marker,
			prefix.SetPrefix,
			categories_aliases.Add,
			categories_aliases.List,
			categories_aliases.Remove,
			vips.Add,
			vips.Remove,
			vips.List,
			vips.SetExpire,
			chat_wall.Delete,
			chat_wall.Ban,
			chat_wall.Timeout,
			chat_wall.Stop,
			shorturl.Command,
			utility.FirstFollowers,
			predictions.Resolve,
			predictions.Cancel,
			predictions.Lock,
			predictions.Start,
			subage.SubAge,
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
) ([]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses, error) {
	return c.services.CommandsCache.Get(ctx, channelId)
}

type FindByMessageResult struct {
	Cmd     *commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses
	FoundBy string
}

// FindByMessage
// Splitting chat message by spaces, then
// read message from end to start, and delete one word from end while message gets empty,
// or we found a command in message
func (c *Commands) FindChannelCommandInInput(
	input string,
	cmds []commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses,
) *FindByMessageResult {
	input = strings.ToLower(input)
	splitName := strings.Fields(input)

	res := FindByMessageResult{}

	length := len(splitName)

	for i := 0; i < length; i++ {
		query := strings.Join(splitName, " ")
		for _, cmd := range cmds {
			if strings.ToLower(cmd.Name) == query {
				res.FoundBy = query
				res.Cmd = &cmd
				break
			}

			if lo.SomeBy(
				cmd.Aliases, func(item string) bool {
					return strings.ToLower(item) == query
				},
			) {
				res.FoundBy = query
				res.Cmd = &cmd
				break
			}
		}

		if res.Cmd != nil {
			break
		} else {
			splitName = splitName[:len(splitName)-1]
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
	userRoles []model.ChannelRole,
	userChannelStats *model.UsersStats,
	dbUser *model.Users,
) *busparser.CommandParseResponse {
	newCtx, span := tracer.CommandsTracer.Start(ctx, "ParseCommandResponses")
	defer span.End()
	ctx = newCtx

	commandsPrefix := requestData.EnrichedData.ChannelCommandPrefix

	result := &busparser.CommandParseResponse{
		KeepOrder: command.Cmd.KeepResponsesOrder,
		IsReply:   command.Cmd.IsReply,
	}

	var cmdParams *string
	params := strings.TrimSpace(requestData.Message.Text[len(commandsPrefix):][len(command.FoundBy):])
	// this shit comes from 7tv for bypass message duplicate
	params = strings.ReplaceAll(params, "\U000e0000", "")
	params = strings.TrimSpace(params)
	if params != "" {
		cmdParams = &params
	}

	var defaultCommand *types.DefaultCommand

	if command.Cmd.Default && command.Cmd.DefaultName != nil {
		cmd, ok := c.DefaultCommands[*command.Cmd.DefaultName]
		if ok {
			defaultCommand = cmd
		}

		result.SkipToxicityCheck = cmd.SkipToxicityCheck
	}

	go c.services.ChannelsCommandsUsagesRepo.Create(
		ctx,
		channelscommandsusages.CreateInput{
			ChannelID: requestData.BroadcasterUserId,
			UserID:    requestData.ChatterUserId,
			CommandID: command.Cmd.ID,
		},
	)

	parseCtxChannel := &types.ParseContextChannel{
		ID:   requestData.BroadcasterUserId,
		Name: requestData.BroadcasterUserLogin,
	}

	parseCtxSender := &types.ParseContextSender{
		ID:               requestData.ChatterUserId,
		Name:             requestData.ChatterUserLogin,
		DisplayName:      requestData.ChatterUserName,
		Color:            requestData.Color,
		Roles:            userRoles,
		DbUser:           dbUser,
		UserChannelStats: userChannelStats,
		IsBroadcaster:    requestData.IsChatterBroadcaster(),
		IsModerator:      requestData.IsChatterModerator(),
		IsVip:            requestData.IsChatterVip(),
		IsSubscriber:     requestData.IsChatterSubscriber(),
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
		RawText:   requestData.Message.Text[len(commandsPrefix):],
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
		ChannelStream: requestData.EnrichedData.ChannelStream,
		Emotes:        emotes,
		Mentions:      mentions,
		Command:       command.Cmd,
	}

	if command.Cmd.Default && defaultCommand != nil {
		argsParser, err := command_arguments.NewParser(
			command_arguments.Opts{
				Args:          defaultCommand.Args,
				Input:         params,
				ArgsDelimiter: defaultCommand.ArgsDelimiter,
			},
		)
		if err != nil {
			if argsParser == nil {
				c.services.Logger.Sugar().Error(
					"error happened on args parser creation for default command",
					zap.Error(err),
					zap.Dict(
						"channel",
						zap.String("id", requestData.BroadcasterUserId),
						zap.String("name", requestData.BroadcasterUserLogin),
					),
				)

				return result
			}

			usage := argsParser.BuildUsageString(ctx, defaultCommand.Args, defaultCommand.Name)

			results := &busparser.CommandParseResponse{
				Responses: []string{fmt.Sprintf("[Usage]: %s", usage)},
				IsReply:   command.Cmd.IsReply,
			}
			return results
		}
		parseCtx.ArgsParser = argsParser

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
				zap.Dict(
					"command",
					zap.String("id", command.Cmd.ID.String()),
					zap.String("name", command.Cmd.Name),
				),
			)

			var commandErr *types.CommandHandlerError

			if errors.As(err, &commandErr) {
				results = &types.CommandsHandlerResult{
					Result: []string{
						fmt.Sprintf("[Twir error]: %s", commandErr.Message),
					},
				}
				fmt.Println(commandErr.Err)
			} else {
				results = &types.CommandsHandlerResult{
					Result: []string{i18n.GetCtx(ctx, locales.Translations.Errors.Generic.TwirError)},
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
		responsesForCategory := make(
			[]commandresponsemodel.Response,
			0,
			len(command.Cmd.Responses),
		)
		for _, r := range command.Cmd.Responses {
			if len(r.TwitchCategoryIDs) > 0 && requestData.EnrichedData.ChannelStream != nil {
				if !lo.ContainsBy(
					r.TwitchCategoryIDs,
					func(categoryId string) bool {
						return categoryId == requestData.EnrichedData.ChannelStream.GameId
					},
				) {
					continue
				}
			}

			if r.OnlineOnly && requestData.EnrichedData.ChannelStream == nil {
				continue
			}

			if r.OfflineOnly && requestData.EnrichedData.ChannelStream != nil {
				continue
			}

			responsesForCategory = append(responsesForCategory, r)
		}

		result.Responses = lo.Map(
			responsesForCategory,
			func(r commandresponsemodel.Response, _ int) string {
				return *r.Text
			},
		)
	}

	wg := &sync.WaitGroup{}

	responses := make([]string, 0, len(result.Responses))

	for _, r := range result.Responses {
		wg.Add(1)

		go func() {
			defer wg.Done()
			variablesResult := c.variablesService.ParseVariablesInText(
				ctx,
				parseCtx,
				r,
			)

			for _, res := range variablesResult {
				responses = append(responses, res)
			}
		}()
	}
	wg.Wait()

	result.Responses = responses

	// TODO: we need some subtype of cmd.Response, with parsed text property, and then we can properly handle sorting by index here
	slices.Reverse(result.Responses)

	return result
}

func (c *Commands) ProcessChatMessage(ctx context.Context, data twitch.TwitchChatMessage) (
	*busparser.CommandParseResponse,
	error,
) {
	if data.EnrichedData.DbUser == nil || data.EnrichedData.DbUserChannelStat == nil {
		return nil, fmt.Errorf("db user or user channel stats is nil")
	}

	commandsPrefix := data.EnrichedData.ChannelCommandPrefix

	if !strings.HasPrefix(data.Message.Text, commandsPrefix) {
		return nil, nil
	}

	cmds, err := c.GetChannelCommands(ctx, data.BroadcasterUserId)
	if err != nil {
		return nil, err
	}

	cmd := c.FindChannelCommandInInput(data.Message.Text[len(commandsPrefix):], cmds)
	if cmd.Cmd == nil || !cmd.Cmd.Enabled {
		return nil, nil
	}

	if cmd.Cmd.ExpiresAt != nil && cmd.Cmd.ExpiresType != nil && cmd.Cmd.ExpiresAt.Before(time.Now().UTC()) {
		if *cmd.Cmd.ExpiresType == commandsmodel.ExpireTypeDisable && cmd.Cmd.Enabled {
			err = c.services.Gorm.
				WithContext(ctx).
				Where(`"id" = ?`, cmd.Cmd.ID).
				Model(&model.ChannelsCommands{}).
				Updates(
					map[string]interface{}{
						"enabled": false,
					},
				).Error
			if err != nil {
				c.services.Logger.Sugar().Error(err)
				return nil, err
			}

			if err := c.services.CommandsCache.Invalidate(ctx, data.BroadcasterUserId); err != nil {
				c.services.Logger.Sugar().Error(err)
				return nil, err
			}
		} else if *cmd.Cmd.ExpiresType == commandsmodel.ExpireTypeDelete && !cmd.Cmd.Default {
			err = c.services.Gorm.
				WithContext(ctx).
				Where(`"id" = ?`, cmd.Cmd.ID).
				Delete(&model.ChannelsCommands{}).Error
			if err != nil {
				c.services.Logger.Sugar().Error(err)
				return nil, err
			}

			if err := c.services.CommandsCache.Invalidate(ctx, data.BroadcasterUserId); err != nil {
				c.services.Logger.Sugar().Error(err)
				return nil, err
			}
		}

		return nil, nil
	}

	if cmd.Cmd.OnlineOnly {
		stream := data.EnrichedData.ChannelStream
		if stream == nil || stream.ID == "" {
			return nil, nil
		}
	}

	if cmd.Cmd.OfflineOnly {
		stream := data.EnrichedData.ChannelStream
		if stream != nil && stream.ID != "" {
			return nil, nil
		}
	}

	if len(cmd.Cmd.EnabledCategories) != 0 {
		stream := &model.ChannelsStreams{}
		err = c.services.Gorm.
			WithContext(ctx).
			Where(`"userId" = ?`, data.BroadcasterUserId).
			Find(stream).Error
		if err != nil {
			return nil, err
		}

		if stream.ID != "" {
			if !lo.ContainsBy(
				cmd.Cmd.EnabledCategories,
				func(category string) bool {
					return category == stream.GameId
				},
			) {
				return nil, nil
			}
		}
	}

	_, userRoles, commandRoles, err := c.prepareCooldownAndPermissionsCheck(
		ctx,
		data.ChatterUserId,
		data.BroadcasterUserId,
		data,
		cmd.Cmd,
	)
	if err != nil {
		return nil, err
	}

	shouldCheckCooldown := c.shouldCheckCooldown(data, cmd.Cmd, userRoles)

	if shouldCheckCooldown {
		isCooldown, err := c.isCooldown(ctx, *cmd.Cmd, data.ChatterUserId, userRoles)
		if err != nil {
			return nil, err
		}

		if isCooldown {
			return nil, nil
		}
	}

	dbUser := &model.Users{
		ID:         data.EnrichedData.DbUser.ID,
		IsBotAdmin: data.EnrichedData.DbUser.IsBotAdmin,
		ApiKey:     data.EnrichedData.DbUser.ApiKey,
		Stats: &model.UsersStats{
			ID:                data.EnrichedData.DbUserChannelStat.ID.String(),
			UserID:            data.EnrichedData.DbUserChannelStat.UserID,
			ChannelID:         data.EnrichedData.DbUserChannelStat.ChannelID,
			Messages:          data.EnrichedData.DbUserChannelStat.Messages,
			Watched:           data.EnrichedData.DbUserChannelStat.Watched,
			UsedChannelPoints: data.EnrichedData.DbUserChannelStat.UsedChannelPoints,
			IsMod:             data.EnrichedData.DbUserChannelStat.IsMod,
			IsVip:             data.EnrichedData.DbUserChannelStat.IsVip,
			IsSubscriber:      data.EnrichedData.DbUserChannelStat.IsSubscriber,
			Reputation:        data.EnrichedData.DbUserChannelStat.Reputation,
			Emotes:            data.EnrichedData.DbUserChannelStat.Emotes,
			CreatedAt:         data.EnrichedData.DbUserChannelStat.CreatedAt,
			UpdatedAt:         data.EnrichedData.DbUserChannelStat.UpdatedAt,
		},
		IsBanned:          data.EnrichedData.DbUser.IsBanned,
		CreatedAt:         data.EnrichedData.DbUser.CreatedAt,
		HideOnLandingPage: data.EnrichedData.DbUser.HideOnLandingPage,
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
		withoutCancel := context.WithoutCancel(ctx)

		c.services.Bus.Events.CommandUsed.Publish(
			withoutCancel,
			events.CommandUsedMessage{
				BaseInfo: events.BaseInfo{
					ChannelID:   data.BroadcasterUserId,
					ChannelName: data.BroadcasterUserLogin,
				},
				CommandID:          cmd.Cmd.ID.String(),
				CommandName:        cmd.Cmd.Name,
				CommandInput:       strings.TrimSpace(data.Message.Text[len(cmd.FoundBy)+1:]),
				UserName:           data.ChatterUserLogin,
				UserDisplayName:    data.ChatterUserName,
				UserID:             data.ChatterUserId,
				IsDefault:          cmd.Cmd.Default,
				DefaultCommandName: cmd.Cmd.DefaultName,
				MessageID:          data.MessageId,
			},
		)

		alert := model.ChannelAlert{}
		if err := c.services.Gorm.WithContext(withoutCancel).Where(
			"channel_id = ? AND command_ids && ?",
			data.BroadcasterUserId,
			pq.StringArray{cmd.Cmd.ID.String()},
		).Find(&alert).Error; err != nil {
			zap.S().Error(err)
			return
		}

		if alert.ID == "" {
			return
		}
		c.services.GrpcClients.WebSockets.TriggerAlert(
			withoutCancel,
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
		userRoles,
		dbUser.Stats,
		dbUser,
	)

	return result, nil
}
