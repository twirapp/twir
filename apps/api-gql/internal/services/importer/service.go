package importer

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	commandsservice "github.com/twirapp/twir/apps/api-gql/internal/services/commands"
	"github.com/twirapp/twir/apps/api-gql/internal/services/roles"
	timersservice "github.com/twirapp/twir/apps/api-gql/internal/services/timers"
	commandentity "github.com/twirapp/twir/libs/entities/command_with_relations"
	timersentity "github.com/twirapp/twir/libs/entities/timers"
	apperrors "github.com/twirapp/twir/libs/errors"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	CommandsService *commandsservice.Service
	TimersService   *timersservice.Service
	RolesService    *roles.Service
}

func New(opts Opts) *Service {
	return newService(opts.CommandsService, opts.TimersService, opts.RolesService)
}

type commandCreator interface {
	Create(context.Context, commandsservice.CreateInput) (commandentity.Command, error)
}

type timerCreator interface {
	Create(context.Context, timersservice.CreateInput) (timersentity.Timer, error)
}

type roleLookup interface {
	GetManyByChannelID(context.Context, string) ([]entity.ChannelRole, error)
}

type Service struct {
	commands commandCreator
	timers   timerCreator
	roles    roleLookup
}

func newService(commands commandCreator, timers timerCreator, roles roleLookup) *Service {
	return &Service{commands: commands, timers: timers, roles: roles}
}

func (s *Service) ImportCommands(
	ctx context.Context,
	channelID, actorID string,
	input []Command,
) (Report, error) {
	report := Report{Failures: []Failure{}}
	if len(input) == 0 {
		return report, nil
	}

	roleIDs, err := s.resolveRoleIDs(ctx, channelID)
	if err != nil {
		return Report{}, fmt.Errorf("resolve channel roles: %w", err)
	}

	for _, command := range input {
		if !validCommand(command) {
			report.addFailure(command.Name, FailureInvalidRecord)
			continue
		}

		requiredRoleTypes, ok := requiredRoleTypes(command.Role)
		if !ok {
			report.addFailure(command.Name, FailureUnsupportedRole)
			continue
		}
		rolesIDs, err := commandRoleIDs(requiredRoleTypes, roleIDs)
		if err != nil {
			return Report{}, fmt.Errorf("resolve required roles for command %q: %w", command.Name, err)
		}

		if _, err := s.commands.Create(ctx, commandCreateInput(channelID, actorID, command, rolesIDs)); err != nil {
			if reason, ok := expectedFailureReason(err); ok {
				report.addFailure(command.Name, reason)
				continue
			}

			return Report{}, fmt.Errorf("create command %q: %w", command.Name, err)
		}

		report.ImportedCount++
	}

	return report, nil
}

func (s *Service) ImportTimers(
	ctx context.Context,
	channelID, actorID string,
	input []Timer,
) (Report, error) {
	report := Report{Failures: []Failure{}}

	for _, timer := range input {
		if !validTimer(timer) {
			report.addFailure(timer.Name, FailureInvalidRecord)
			continue
		}

		if _, err := s.timers.Create(ctx, timerCreateInput(channelID, actorID, timer)); err != nil {
			if reason, ok := expectedFailureReason(err); ok {
				report.addFailure(timer.Name, reason)
				continue
			}

			return Report{}, fmt.Errorf("create timer %q: %w", timer.Name, err)
		}

		report.ImportedCount++
	}

	return report, nil
}

func (s *Service) resolveRoleIDs(ctx context.Context, channelID string) (map[entity.ChannelRoleEnum]string, error) {
	channelRoles, err := s.roles.GetManyByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	roleIDs := make(map[entity.ChannelRoleEnum]string, len(channelRoles))
	for _, role := range channelRoles {
		roleIDs[role.Type] = role.ID.String()
	}

	return roleIDs, nil
}

func requiredRoleTypes(role RoleRequirement) ([]entity.ChannelRoleEnum, bool) {
	switch role {
	case RoleEveryone:
		return []entity.ChannelRoleEnum{}, true
	case RoleSubscriber:
		return []entity.ChannelRoleEnum{
			entity.ChannelRoleTypeBroadcaster,
			entity.ChannelRoleTypeModerator,
			entity.ChannelRoleTypeSubscriber,
		}, true
	case RoleVip:
		return []entity.ChannelRoleEnum{
			entity.ChannelRoleTypeBroadcaster,
			entity.ChannelRoleTypeModerator,
			entity.ChannelRoleTypeVip,
		}, true
	case RoleModerator:
		return []entity.ChannelRoleEnum{
			entity.ChannelRoleTypeBroadcaster,
			entity.ChannelRoleTypeModerator,
		}, true
	case RoleBroadcaster:
		return []entity.ChannelRoleEnum{entity.ChannelRoleTypeBroadcaster}, true
	default:
		return nil, false
	}
}

func commandRoleIDs(required []entity.ChannelRoleEnum, roleIDs map[entity.ChannelRoleEnum]string) ([]string, error) {
	ids := make([]string, 0, len(required))
	for _, roleType := range required {
		id, ok := roleIDs[roleType]
		if !ok || id == uuid.Nil.String() {
			return nil, fmt.Errorf("required %s role is not configured", roleType)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func validCommand(command Command) bool {
	return strings.TrimSpace(command.Name) != "" &&
		strings.TrimSpace(command.Response) != "" &&
		!(command.OnlineOnly && command.OfflineOnly)
}

func validTimer(timer Timer) bool {
	return strings.TrimSpace(timer.Name) != "" &&
		strings.TrimSpace(timer.Message) != "" &&
		timer.TimeInterval > 0 &&
		timer.MessageInterval > 0 &&
		(timer.OnlineEnabled || timer.OfflineEnabled)
}

func commandCreateInput(channelID, actorID string, command Command, roleIDs []string) commandsservice.CreateInput {
	response := command.Response
	return commandsservice.CreateInput{
		ChannelID:                 channelID,
		ActorID:                   actorID,
		Name:                      command.Name,
		Cooldown:                  command.Cooldown,
		CooldownType:              "GLOBAL",
		Enabled:                   command.Enabled,
		Aliases:                   command.Aliases,
		Description:               "",
		Visible:                   command.Visible,
		IsReply:                   command.IsReply,
		KeepResponsesOrder:        true,
		DeniedUsersIDS:            []string{},
		AllowedUsersIDS:           []string{},
		RolesIDS:                  roleIDs,
		OnlineOnly:                command.OnlineOnly,
		EnabledCategories:         []string{},
		RequiredWatchTime:         0,
		RequiredMessages:          0,
		RequiredUsedChannelPoints: 0,
		Responses: []commandsservice.CreateInputResponse{{
			Text:              &response,
			Order:             0,
			TwitchCategoryIDs: []string{},
			OnlineOnly:        command.OnlineOnly,
			OfflineOnly:       command.OfflineOnly,
		}},
	}
}

func timerCreateInput(channelID, actorID string, timer Timer) timersservice.CreateInput {
	return timersservice.CreateInput{
		ChannelID:       channelID,
		ActorID:         actorID,
		Name:            timer.Name,
		Enabled:         timer.Enabled,
		OnlineEnabled:   timer.OnlineEnabled,
		OfflineEnabled:  timer.OfflineEnabled,
		TimeInterval:    timer.TimeInterval,
		MessageInterval: timer.MessageInterval,
		Responses: []timersservice.CreateResponse{{
			Text:  timer.Message,
			Count: 1,
		}},
	}
}

func expectedFailureReason(err error) (FailureReason, bool) {
	appErr, ok := apperrors.AsAppError(err)
	if !ok {
		return "", false
	}

	reason, ok := appErr.Details["reason"].(string)
	if !ok {
		return "", false
	}

	switch FailureReason(reason) {
	case FailureDuplicate, FailurePlanLimit:
		return FailureReason(reason), true
	default:
		return "", false
	}
}

func (r *Report) addFailure(name string, reason FailureReason) {
	r.FailedCount++
	r.Failures = append(r.Failures, Failure{Name: name, Reason: reason})
}
