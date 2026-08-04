package variables

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/audit"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/executron"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/cache/twitch"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/plans"
	"github.com/twirapp/twir/libs/repositories/variables"
	"github.com/twirapp/twir/libs/repositories/variables/model"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Service struct {
	twirbus             *buscore.Bus
	config              config.Config
	cachedTwitchClient  *twitch.CachedTwitchClient
	gorm                *gorm.DB
	auditRecorder       audit.Recorder
	variablesRepository variables.Repository
	plansRepository     plans.Repository
	channelService      *channelservice.ChannelService
}

func New(
	twirBus *buscore.Bus,
	cfg config.Config,
	cachedTwitchClient *twitch.CachedTwitchClient,
	db *gorm.DB,
	auditRecorder audit.Recorder,
	variablesRepository variables.Repository,
	plansRepository plans.Repository,
	channelService *channelservice.ChannelService,
) *Service {
	return &Service{
		twirbus:             twirBus,
		config:              cfg,
		cachedTwitchClient:  cachedTwitchClient,
		gorm:                db,
		auditRecorder:       auditRecorder,
		variablesRepository: variablesRepository,
		plansRepository:     plansRepository,
		channelService:      channelService,
	}
}

var ErrNotFound = errors.New("variable not found")

func (c *Service) GetBuiltIn(ctx context.Context) ([]parser.BuiltInVariable, error) {
	result, err := c.twirbus.Parser.GetBuiltInVariables.Request(ctx, struct{}{})
	if err != nil {
		return nil, fmt.Errorf("cannot get built-in variables: %w", err)
	}

	return result.Data, nil
}

func (c *Service) dbToModel(m model.CustomVariable) entity.CustomVariable {
	return entity.CustomVariable{
		ID:             m.ID,
		ChannelID:      m.ChannelID,
		Name:           m.Name,
		Description:    m.Description.Ptr(),
		Type:           entity.CustomVarType(m.Type),
		EvalValue:      m.EvalValue,
		Response:       m.Response,
		ScriptLanguage: entity.CustomVarScriptLanguage(m.ScriptLanguage),
	}
}

func (c *Service) EvaluateScript(
	ctx context.Context,
	channelID,
	script string,
	language entity.CustomVarScriptLanguage,
	testAsUserName *string,
) (string, error) {
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return "", fmt.Errorf("cannot parse channel id: %w", err)
	}

	channel, err := c.channelService.GetChannelByID(ctx, parsedChannelID)
	if err != nil {
		return "", fmt.Errorf("cannot get channel: %w", err)
	}

	if testAsUserName != nil && *testAsUserName != "" {
		twitchBinding, found := channel.Binding(platform.PlatformTwitch)
		if !found {
			return "", fmt.Errorf("channel has no twitch platform ID")
		}

		var channelUser, user helix.User
		var wg errgroup.Group

		wg.Go(
			func() error {
				u, err := c.cachedTwitchClient.GetUserById(ctx, twitchBinding.PlatformChannelID)
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

		platformSource := platform.PlatformTwitch
		preparedEvalValue, err := c.twirbus.Parser.ParseVariablesInText.Request(
			ctx, parser.ParseVariablesInTextRequest{
				ChannelID:      parsedChannelID,
				ChannelName:    channelUser.Login,
				Text:           script,
				UserID:         user.ID,
				UserLogin:      user.Login,
				UserName:       user.DisplayName,
				IsCommand:      true,
				IsInCustomVar:  true,
				PlatformSource: &platformSource,
			},
		)
		if err != nil {
			return "", fmt.Errorf("cannot parse variables in text: %w", err)
		}

		result, err := c.twirbus.Executron.Execute.Request(
			ctx,
			executron.ExecuteRequest{
				ChannelId: channelID,
				Language:  language.String(),
				Code:      preparedEvalValue.Data.Text,
			},
		)
		if err != nil {
			return "", fmt.Errorf("cannot evaluate script: %w", err)
		}

		var res string
		if result.Data.Result != "" {
			res = result.Data.Result
		} else if result.Data.Error != "" {
			res = result.Data.Error
		}

		return res, nil
	}

	result, err := c.twirbus.Executron.Execute.Request(
		ctx,
		executron.ExecuteRequest{
			ChannelId: channelID,
			Language:  language.String(),
			Code:      script,
		},
	)
	if err != nil {
		return "", fmt.Errorf("cannot evaluate script: %w", err)
	}

	var res string
	if result.Data.Result != "" {
		res = result.Data.Result
	} else if result.Data.Error != "" {
		res = result.Data.Error
	}

	return res, nil
}
