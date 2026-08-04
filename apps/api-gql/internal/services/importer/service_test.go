package importer

import (
	"context"
	stderrors "errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	commandsservice "github.com/twirapp/twir/apps/api-gql/internal/services/commands"
	timersservice "github.com/twirapp/twir/apps/api-gql/internal/services/timers"
	commandentity "github.com/twirapp/twir/libs/entities/command_with_relations"
	timersentity "github.com/twirapp/twir/libs/entities/timers"
	"github.com/twirapp/twir/libs/errors"
)

func TestImportCommandsCreatesNormalizedCommand(t *testing.T) {
	t.Parallel()

	broadcasterID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	moderatorID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	subscriberID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	roles := &fakeRoleLookup{roles: []entity.ChannelRole{
		{ID: subscriberID, Type: entity.ChannelRoleTypeSubscriber},
		{ID: moderatorID, Type: entity.ChannelRoleTypeModerator},
		{ID: broadcasterID, Type: entity.ChannelRoleTypeBroadcaster},
	}}
	commands := &fakeCommandCreator{}
	service := newService(commands, &fakeTimerCreator{}, roles)

	report, err := service.ImportCommands(context.Background(), "channel", "actor", []Command{
		{
			Name:        "hello",
			Response:    "world",
			Enabled:     true,
			Visible:     true,
			IsReply:     true,
			Aliases:     []string{"hi"},
			Cooldown:    10,
			Role:        RoleSubscriber,
			OnlineOnly:  true,
			OfflineOnly: false,
		},
		{Name: "owner", Response: "owner response", Role: RoleBroadcaster},
	})
	if err != nil {
		t.Fatalf("ImportCommands() error = %v", err)
	}

	if want := (Report{ImportedCount: 2, Failures: []Failure{}}); !reflect.DeepEqual(report, want) {
		t.Fatalf("ImportCommands() report = %#v, want %#v", report, want)
	}
	if roles.calls != 1 {
		t.Fatalf("role lookup calls = %d, want 1", roles.calls)
	}
	if len(commands.inputs) != 2 {
		t.Fatalf("created commands = %d, want 2", len(commands.inputs))
	}

	response := "world"
	want := commandsservice.CreateInput{
		ChannelID:                 "channel",
		ActorID:                   "actor",
		Name:                      "hello",
		Cooldown:                  10,
		CooldownType:              "GLOBAL",
		Enabled:                   true,
		Aliases:                   []string{"hi"},
		Description:               "",
		Visible:                   true,
		IsReply:                   true,
		KeepResponsesOrder:        true,
		DeniedUsersIDS:            []string{},
		AllowedUsersIDS:           []string{},
		RolesIDS:                  []string{broadcasterID.String(), moderatorID.String(), subscriberID.String()},
		OnlineOnly:                true,
		EnabledCategories:         []string{},
		RequiredWatchTime:         0,
		RequiredMessages:          0,
		RequiredUsedChannelPoints: 0,
		Responses: []commandsservice.CreateInputResponse{{
			Text:              &response,
			Order:             0,
			TwitchCategoryIDs: []string{},
			OnlineOnly:        true,
			OfflineOnly:       false,
		}},
	}
	if got := commands.inputs[0]; !reflect.DeepEqual(got, want) {
		t.Fatalf("created command = %#v, want %#v", got, want)
	}
}

func TestImportCommandsReportsExpectedCreateFailures(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name   string
		err    error
		reason FailureReason
	}{
		{
			name:   "duplicate",
			err:    errors.NewConflictError("A command with this name or alias already exists").WithDetails(map[string]any{"reason": "DUPLICATE"}),
			reason: FailureDuplicate,
		},
		{
			name:   "plan limit",
			err:    errors.NewBadRequestError("You have reached the maximum limit of 10 commands").WithDetails(map[string]any{"reason": "PLAN_LIMIT"}),
			reason: FailurePlanLimit,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			commands := &fakeCommandCreator{errForName: map[string]error{"first": tt.err}}
			service := newService(commands, &fakeTimerCreator{}, &fakeRoleLookup{})

			report, err := service.ImportCommands(context.Background(), "channel", "actor", []Command{{
				Name: "first", Response: "response", Role: RoleEveryone,
			}})
			if err != nil {
				t.Fatalf("ImportCommands() error = %v", err)
			}

			want := Report{FailedCount: 1, Failures: []Failure{{Name: "first", Reason: tt.reason}}}
			if !reflect.DeepEqual(report, want) {
				t.Fatalf("ImportCommands() report = %#v, want %#v", report, want)
			}
		})
	}
}

func TestImportCommandsKeepsSuccessfulSiblings(t *testing.T) {
	t.Parallel()

	commands := &fakeCommandCreator{errForName: map[string]error{
		"duplicate": errors.NewConflictError("A command with this name or alias already exists").WithDetails(map[string]any{"reason": "DUPLICATE"}),
	}}
	service := newService(commands, &fakeTimerCreator{}, &fakeRoleLookup{})

	report, err := service.ImportCommands(context.Background(), "channel", "actor", []Command{
		{Name: "first", Response: "one", Role: RoleEveryone},
		{Name: "duplicate", Response: "two", Role: RoleEveryone},
		{Name: "third", Response: "three", Role: RoleEveryone},
	})
	if err != nil {
		t.Fatalf("ImportCommands() error = %v", err)
	}

	want := Report{ImportedCount: 2, FailedCount: 1, Failures: []Failure{{Name: "duplicate", Reason: FailureDuplicate}}}
	if !reflect.DeepEqual(report, want) {
		t.Fatalf("ImportCommands() report = %#v, want %#v", report, want)
	}
	if got := len(commands.inputs); got != 3 {
		t.Fatalf("creator calls = %d, want 3", got)
	}
}

func TestImportCommandsReportsUnsupportedAndInvalidRecords(t *testing.T) {
	t.Parallel()

	commands := &fakeCommandCreator{}
	service := newService(commands, &fakeTimerCreator{}, &fakeRoleLookup{})

	report, err := service.ImportCommands(context.Background(), "channel", "actor", []Command{
		{Name: "unknown-role", Response: "response", Role: RoleRequirement("UNKNOWN")},
		{Name: "", Response: "response", Role: RoleEveryone},
		{Name: "both", Response: "response", Role: RoleEveryone, OnlineOnly: true, OfflineOnly: true},
	})
	if err != nil {
		t.Fatalf("ImportCommands() error = %v", err)
	}

	want := Report{FailedCount: 3, Failures: []Failure{
		{Name: "unknown-role", Reason: FailureUnsupportedRole},
		{Name: "", Reason: FailureInvalidRecord},
		{Name: "both", Reason: FailureInvalidRecord},
	}}
	if !reflect.DeepEqual(report, want) {
		t.Fatalf("ImportCommands() report = %#v, want %#v", report, want)
	}
	if got := len(commands.inputs); got != 0 {
		t.Fatalf("creator calls = %d, want 0", got)
	}
}

func TestImportCommandsReturnsInfrastructureErrorsImmediately(t *testing.T) {
	t.Parallel()

	commands := &fakeCommandCreator{errForName: map[string]error{"first": stderrors.New("database unavailable")}}
	service := newService(commands, &fakeTimerCreator{}, &fakeRoleLookup{})

	_, err := service.ImportCommands(context.Background(), "channel", "actor", []Command{
		{Name: "first", Response: "one", Role: RoleEveryone},
		{Name: "second", Response: "two", Role: RoleEveryone},
	})
	if !stderrors.Is(err, commands.errForName["first"]) {
		t.Fatalf("ImportCommands() error = %v, want infrastructure error", err)
	}
	if got := len(commands.inputs); got != 1 {
		t.Fatalf("creator calls = %d, want 1", got)
	}
}

func TestImportCommandsReturnsMissingRequiredRoleAsAnOperationError(t *testing.T) {
	t.Parallel()

	commands := &fakeCommandCreator{}
	service := newService(commands, &fakeTimerCreator{}, &fakeRoleLookup{roles: []entity.ChannelRole{
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Type: entity.ChannelRoleTypeBroadcaster},
	}})

	_, err := service.ImportCommands(context.Background(), "channel", "actor", []Command{{
		Name: "subscriber", Response: "response", Role: RoleSubscriber,
	}})
	if err == nil {
		t.Fatal("ImportCommands() error = nil, want missing required role error")
	}
	if got := len(commands.inputs); got != 0 {
		t.Fatalf("creator calls = %d, want 0", got)
	}
}

func TestImportTimersCreatesNormalizedTimer(t *testing.T) {
	t.Parallel()

	timers := &fakeTimerCreator{}
	service := newService(&fakeCommandCreator{}, timers, &fakeRoleLookup{})

	report, err := service.ImportTimers(context.Background(), "channel", "actor", []Timer{{
		Name: "reminder", Message: "hydrate", Enabled: true, OnlineEnabled: true,
		TimeInterval: 15, MessageInterval: 3,
	}})
	if err != nil {
		t.Fatalf("ImportTimers() error = %v", err)
	}
	if want := (Report{ImportedCount: 1, Failures: []Failure{}}); !reflect.DeepEqual(report, want) {
		t.Fatalf("ImportTimers() report = %#v, want %#v", report, want)
	}
	want := timersservice.CreateInput{
		ChannelID: "channel", ActorID: "actor", Name: "reminder", Enabled: true, OnlineEnabled: true,
		TimeInterval: 15, MessageInterval: 3,
		Responses: []timersservice.CreateResponse{{Text: "hydrate", Count: 1}},
	}
	if got := timers.inputs[0]; !reflect.DeepEqual(got, want) {
		t.Fatalf("created timer = %#v, want %#v", got, want)
	}
}

func TestImportTimersReportsInvalidRecordsAndPlanLimit(t *testing.T) {
	t.Parallel()

	timers := &fakeTimerCreator{errForName: map[string]error{
		"limited": errors.NewBadRequestError("You have reached the maximum limit of 3 timers").WithDetails(map[string]any{"reason": "PLAN_LIMIT"}),
	}}
	service := newService(&fakeCommandCreator{}, timers, &fakeRoleLookup{})

	report, err := service.ImportTimers(context.Background(), "channel", "actor", []Timer{
		{Name: "limited", Message: "one", OnlineEnabled: true, TimeInterval: 1, MessageInterval: 1},
		{Name: "invalid", Message: "two", TimeInterval: 0, MessageInterval: 1},
	})
	if err != nil {
		t.Fatalf("ImportTimers() error = %v", err)
	}
	want := Report{FailedCount: 2, Failures: []Failure{
		{Name: "limited", Reason: FailurePlanLimit},
		{Name: "invalid", Reason: FailureInvalidRecord},
	}}
	if !reflect.DeepEqual(report, want) {
		t.Fatalf("ImportTimers() report = %#v, want %#v", report, want)
	}
}

func TestImportTimersReportsDuplicateCreateFailure(t *testing.T) {
	t.Parallel()

	timers := &fakeTimerCreator{errForName: map[string]error{
		"duplicate": errors.NewConflictError("A timer with this name already exists").WithDetails(map[string]any{"reason": "DUPLICATE"}),
	}}
	service := newService(&fakeCommandCreator{}, timers, &fakeRoleLookup{})

	report, err := service.ImportTimers(context.Background(), "channel", "actor", []Timer{{
		Name: "duplicate", Message: "message", OnlineEnabled: true, TimeInterval: 1, MessageInterval: 1,
	}})
	if err != nil {
		t.Fatalf("ImportTimers() error = %v", err)
	}
	want := Report{FailedCount: 1, Failures: []Failure{{Name: "duplicate", Reason: FailureDuplicate}}}
	if !reflect.DeepEqual(report, want) {
		t.Fatalf("ImportTimers() report = %#v, want %#v", report, want)
	}
}

func TestImportTimersReturnsInfrastructureErrorsImmediately(t *testing.T) {
	t.Parallel()

	timers := &fakeTimerCreator{errForName: map[string]error{"first": stderrors.New("timer repository unavailable")}}
	service := newService(&fakeCommandCreator{}, timers, &fakeRoleLookup{})

	_, err := service.ImportTimers(context.Background(), "channel", "actor", []Timer{
		{Name: "first", Message: "one", OnlineEnabled: true, TimeInterval: 1, MessageInterval: 1},
		{Name: "second", Message: "two", OnlineEnabled: true, TimeInterval: 1, MessageInterval: 1},
	})
	if !stderrors.Is(err, timers.errForName["first"]) {
		t.Fatalf("ImportTimers() error = %v, want infrastructure error", err)
	}
	if got := len(timers.inputs); got != 1 {
		t.Fatalf("creator calls = %d, want 1", got)
	}
}

type fakeCommandCreator struct {
	inputs     []commandsservice.CreateInput
	errForName map[string]error
}

func (f *fakeCommandCreator) Create(_ context.Context, input commandsservice.CreateInput) (commandentity.Command, error) {
	f.inputs = append(f.inputs, input)
	if err := f.errForName[input.Name]; err != nil {
		return commandentity.CommandNil, err
	}
	return commandentity.Command{}, nil
}

type fakeTimerCreator struct {
	inputs     []timersservice.CreateInput
	errForName map[string]error
}

func (f *fakeTimerCreator) Create(_ context.Context, input timersservice.CreateInput) (timersentity.Timer, error) {
	f.inputs = append(f.inputs, input)
	if err := f.errForName[input.Name]; err != nil {
		return timersentity.Nil, err
	}
	return timersentity.Timer{}, nil
}

type fakeRoleLookup struct {
	roles []entity.ChannelRole
	err   error
	calls int
}

func (f *fakeRoleLookup) GetManyByChannelID(_ context.Context, _ string) ([]entity.ChannelRole, error) {
	f.calls++
	return f.roles, f.err
}
