package commands

import (
	"context"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"
	"github.com/satont/tsuwari/apps/parser-new/internal/types/services"
	"github.com/satont/tsuwari/apps/parser-new/internal/variables"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/gopool"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"regexp"
	"sort"
	"strings"
	"sync"
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

func New(opts *Opts) Commands {
	commands := make(map[string]*types.DefaultCommand)

	ctx := Commands{
		DefaultCommands:    commands,
		parseResponsesPool: gopool.NewPool(100),
		services:           opts.Services,
		variablesService:   opts.VariablesService,
	}

	return ctx
}

func (c *Commands) GetChannelCommands(ctx context.Context, channelId string) ([]*model.ChannelsCommands, error) {
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
func (c *Commands) FindChannelCommandInInput(input string, cmds []*model.ChannelsCommands) *FindByMessageResult {
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

			if lo.SomeBy(cmd.Aliases, func(item string) bool {
				return item == query
			}) {
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
		sort.Slice(res.Cmd.Responses, func(a, b int) bool {
			return res.Cmd.Responses[a].Order < res.Cmd.Responses[b].Order
		})
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
		Create(&model.ChannelsCommandsUsages{
			ID:        uuid.New().String(),
			UserID:    requestData.Sender.Id,
			ChannelID: requestData.Channel.Id,
			CommandID: command.Cmd.ID,
		})

	parseCtx := &types.ParseContext{
		Channel: &types.ParseContextChannel{
			ID:   requestData.Channel.Id,
			Name: requestData.Channel.Name,
		},
		Sender: &types.ParseContextSender{
			ID:          requestData.Sender.Id,
			Name:        requestData.Sender.Name,
			DisplayName: requestData.Sender.DisplayName,
			Badges:      requestData.Sender.Badges,
		},
		Text:      cmdParams,
		IsCommand: true,
		Services:  c.services,
	}

	if command.Cmd.Default && defaultCommand != nil {
		results := defaultCommand.Handler(ctx, parseCtx)

		result.Responses = lo.If(results == nil, []string{}).Else(results.Result)
	} else {
		result.Responses = lo.Map(command.Cmd.Responses, func(r *model.ChannelsCommandsResponses, _ int) string {
			return r.Text.String
		})
	}

	wg := &sync.WaitGroup{}
	for i, r := range result.Responses {
		wg.Add(1)

		index := i
		response := r

		c.parseResponsesPool.Submit(func() {
			result.Responses[index] = c.variablesService.ParseVariablesInText(ctx, parseCtx, response)
			wg.Done()
		})
	}
	wg.Wait()

	return result
}
