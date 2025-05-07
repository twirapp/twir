package variables

import (
	"context"
	"errors"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/parser/pkg/executron"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/repositories/variables"
	"github.com/twirapp/twir/libs/repositories/variables/model"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	TwirBus             *buscore.Bus
	Config              config.Config
	TokensGrpc          tokens.TokensClient
	CachedTwitchClient  *twitch.CachedTwitchClient
	Gorm                *gorm.DB
	Logger              logger.Logger
	VariablesRepository variables.Repository
	Executron           executron.Executron
}

type Service struct {
	twirbus             *buscore.Bus
	config              config.Config
	tokensGrpc          tokens.TokensClient
	cachedTwitchClient  *twitch.CachedTwitchClient
	gorm                *gorm.DB
	logger              logger.Logger
	variablesRepository variables.Repository
	executron           executron.Executron
}

func New(opts Opts) *Service {
	return &Service{
		twirbus:             opts.TwirBus,
		config:              opts.Config,
		tokensGrpc:          opts.TokensGrpc,
		cachedTwitchClient:  opts.CachedTwitchClient,
		gorm:                opts.Gorm,
		logger:              opts.Logger,
		variablesRepository: opts.VariablesRepository,
		executron:           opts.Executron,
	}
}

const MaxPerChannel = 10

var ErrNotFound = errors.New("variable not found")

func (c *Service) dbToModel(m model.CustomVariable) entity.CustomVariable {
	return entity.CustomVariable{
		ID:          m.ID,
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

		result, err := c.executron.ExecuteUserCode(ctx, "javascript", preparedEvalValue.Data.Text)
		if err != nil {
			return "", fmt.Errorf("cannot evaluate script: %w", err)
		}

		var res string
		if result.Result != "" {
			res = result.Result
		} else if result.Error != "" {
			res = result.Error
		}

		return res, nil
	}

	result, err := c.executron.ExecuteUserCode(ctx, "javascript", script)
	if err != nil {
		return "", fmt.Errorf("cannot evaluate script: %w", err)
	}

	var res string
	if result.Result != "" {
		res = result.Result
	} else if result.Error != "" {
		res = result.Error
	}

	return res, nil
}
