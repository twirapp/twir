package timers

import (
	"context"
	stderrors "errors"
	"io"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/audit"
	timersbusservice "github.com/twirapp/twir/libs/bus-core/timers"
	planentity "github.com/twirapp/twir/libs/entities/plan"
	timersentity "github.com/twirapp/twir/libs/entities/timers"
	plansrepository "github.com/twirapp/twir/libs/repositories/plans"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
)

func TestCreateCompensatesWhenTimerLifecyclePublishFails(t *testing.T) {
	t.Parallel()

	publishErr := stderrors.New("timer bus unavailable")
	timerID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	repository := &timerCreateRepository{createdTimer: timersentity.Timer{ID: timerID}}
	auditRecorder := &timerCreateAuditRecorder{}
	service := &Service{
		logger:           slog.New(slog.NewTextHandler(io.Discard, nil)),
		auditRecorder:    auditRecorder,
		timersRepository: repository,
		plansRepository:  timerCreatePlansRepository{},
		timerLifecycle:   timerCreateLifecycle{err: publishErr},
	}

	_, err := service.Create(context.Background(), CreateInput{
		ChannelID: "channel",
		ActorID:   "actor",
		Name:      "reminder",
		Enabled:   true,
		Responses: []CreateResponse{{Text: "hydrate"}},
	})
	if !stderrors.Is(err, publishErr) {
		t.Fatalf("Create() error = %v, want lifecycle publish error", err)
	}
	if got := repository.deletedIDs; len(got) != 1 || got[0] != timerID {
		t.Fatalf("deleted timer IDs = %v, want [%s]", got, timerID)
	}
	if auditRecorder.createCalls != 0 {
		t.Fatalf("create audit calls = %d, want 0", auditRecorder.createCalls)
	}
}

func TestCreateReturnsPublishAndCleanupErrorsWhenCompensationFails(t *testing.T) {
	t.Parallel()

	publishErr := stderrors.New("timer bus unavailable")
	cleanupErr := stderrors.New("timer cleanup unavailable")
	repository := &timerCreateRepository{
		createdTimer: timersentity.Timer{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
		deleteErr:    cleanupErr,
	}
	service := &Service{
		logger:           slog.New(slog.NewTextHandler(io.Discard, nil)),
		auditRecorder:    &timerCreateAuditRecorder{},
		timersRepository: repository,
		plansRepository:  timerCreatePlansRepository{},
		timerLifecycle:   timerCreateLifecycle{err: publishErr},
	}

	_, err := service.Create(context.Background(), CreateInput{
		ChannelID: "channel", ActorID: "actor", Name: "reminder", Enabled: true,
		Responses: []CreateResponse{{Text: "hydrate"}},
	})
	if !stderrors.Is(err, publishErr) || !stderrors.Is(err, cleanupErr) {
		t.Fatalf("Create() error = %v, want publish and cleanup errors", err)
	}
}

type timerCreateRepository struct {
	createdTimer timersentity.Timer
	deletedIDs   []uuid.UUID
	deleteErr    error
}

func (r *timerCreateRepository) GetByID(context.Context, uuid.UUID) (timersentity.Timer, error) {
	return timersentity.Nil, nil
}

func (r *timerCreateRepository) GetAllByChannelID(context.Context, string) ([]timersentity.Timer, error) {
	return []timersentity.Timer{}, nil
}

func (r *timerCreateRepository) CountByChannelID(context.Context, string) (int, error) {
	return 0, nil
}

func (r *timerCreateRepository) Create(context.Context, timersrepository.CreateInput) (timersentity.Timer, error) {
	return r.createdTimer, nil
}

func (r *timerCreateRepository) UpdateByID(context.Context, uuid.UUID, timersrepository.UpdateInput) (timersentity.Timer, error) {
	return timersentity.Nil, nil
}

func (r *timerCreateRepository) Delete(_ context.Context, id uuid.UUID) error {
	r.deletedIDs = append(r.deletedIDs, id)
	return r.deleteErr
}

func (r *timerCreateRepository) Count(context.Context, timersrepository.CountInput) (int64, error) {
	return 0, nil
}

func (r *timerCreateRepository) GetMany(context.Context, timersrepository.GetManyInput) ([]timersentity.Timer, error) {
	return nil, nil
}

type timerCreatePlansRepository struct{}

func (timerCreatePlansRepository) GetByID(context.Context, string) (planentity.Plan, error) {
	return planentity.Nil, nil
}

func (timerCreatePlansRepository) GetByName(context.Context, string) (planentity.Plan, error) {
	return planentity.Nil, nil
}

func (timerCreatePlansRepository) GetByChannelID(context.Context, string) (planentity.Plan, error) {
	return planentity.Plan{MaxTimers: 1}, nil
}

func (timerCreatePlansRepository) GetManyByIDs(context.Context, []string) ([]planentity.Plan, error) {
	return nil, nil
}

var _ plansrepository.Repository = timerCreatePlansRepository{}

type timerCreateLifecycle struct {
	err error
}

func (l timerCreateLifecycle) Publish(context.Context, bool, timersbusservice.AddOrRemoveTimerRequest) error {
	return l.err
}

type timerCreateAuditRecorder struct {
	createCalls int
}

func (r *timerCreateAuditRecorder) RecordCreateOperation(context.Context, audit.CreateOperation) error {
	r.createCalls++
	return nil
}

func (*timerCreateAuditRecorder) RecordDeleteOperation(context.Context, audit.DeleteOperation) error {
	return nil
}

func (*timerCreateAuditRecorder) RecordUpdateOperation(context.Context, audit.UpdateOperation) error {
	return nil
}

var _ audit.Recorder = (*timerCreateAuditRecorder)(nil)
var _ timersrepository.Repository = (*timerCreateRepository)(nil)
