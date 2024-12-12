package variables

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	config "github.com/satont/twir/libs/config"
	dbmodels "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/eval"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	TwirBus            *buscore.Bus
	Config             config.Config
	TokensGrpc         tokens.TokensClient
	CachedTwitchClient *twitch.CachedTwitchClient
	Gorm               *gorm.DB
	Logger             logger.Logger
}

type Service struct {
	twirbus            *buscore.Bus
	config             config.Config
	tokensGrpc         tokens.TokensClient
	cachedTwitchClient *twitch.CachedTwitchClient
	gorm               *gorm.DB
	logger             logger.Logger
}

func New(opts Opts) *Service {
	return &Service{
		twirbus:            opts.TwirBus,
		config:             opts.Config,
		tokensGrpc:         opts.TokensGrpc,
		cachedTwitchClient: opts.CachedTwitchClient,
		gorm:               opts.Gorm,
		logger:             opts.Logger,
	}
}

func (c *Service) dbToModel(m dbmodels.ChannelsCustomvars) entity.CustomVariable {
	return entity.CustomVariable{
		ID:          uuid.MustParse(m.ID),
		ChannelID:   m.ChannelID,
		Name:        m.Name,
		Description: m.Description.Ptr(),
		Type:        entity.CustomVarType(m.Type),
		EvalValue:   m.EvalValue,
		Response:    m.Response,
	}
}

func (c *Service) EvaluateScript(
	ctx context.Context,
	channelID, script string,
	testAsUserName *string,
) (string, error) {
	if testAsUserName != nil && *testAsUserName != "" {
		var channelUser, user helix.User
		var wg errgroup.Group

		wg.Go(
			func() error {
				u, err := c.cachedTwitchClient.GetUserById(ctx, channelID)
				if err != nil {
					return fmt.Errorf("cannot get channel user: %w", err)
				}

				channelUser = u.User
				return nil
			},
		)

		wg.Go(
			func() error {
				u, err := c.cachedTwitchClient.GetUserByName(ctx, *testAsUserName)
				if err != nil {
					return fmt.Errorf("cannot get user: %w", err)
				}

				user = u.User

				return nil
			},
		)

		if err := wg.Wait(); err != nil {
			return "", err
		}

		preparedEvalValue, err := c.twirbus.Parser.ParseVariablesInText.Request(
			ctx, parser.ParseVariablesInTextRequest{
				ChannelID:     channelUser.ID,
				ChannelName:   channelUser.Login,
				Text:          script,
				UserID:        user.ID,
				UserLogin:     user.Login,
				UserName:      user.DisplayName,
				IsCommand:     true,
				IsInCustomVar: true,
			},
		)
		if err != nil {
			return "", fmt.Errorf("cannot parse variables in text: %w", err)
		}

		result, err := c.twirbus.Eval.Evaluate.Request(
			ctx,
			eval.EvalRequest{
				Expression: preparedEvalValue.Data.Text,
			},
		)
		if err != nil {
			return "", fmt.Errorf("cannot evaluate script: %w", err)
		}

		return result.Data.Result, nil
	}

	result, err := c.twirbus.Eval.Evaluate.Request(
		ctx,
		eval.EvalRequest{
			Expression: script,
		},
	)
	if err != nil {
		return "", fmt.Errorf("cannot evaluate script: %w", err)
	}

	return result.Data.Result, nil
}
