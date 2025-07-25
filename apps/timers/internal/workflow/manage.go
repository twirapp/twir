package workflow

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/timers/internal/shared"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

func (c *Workflow) AddTimer(ctx context.Context, timerId string) error {
	parsedUuid, err := uuid.Parse(timerId)
	if err != nil {
		return err
	}

	timer, err := c.timersRepository.GetByID(ctx, parsedUuid)
	if err != nil {
		return err
	}

	scheduleID := timer.ID.String()
	workflowID := "timers-" + timer.ID.String()

	var every time.Duration
	if c.config.AppEnv == "development" {
		every = time.Duration(timer.TimeInterval) * time.Second
	} else {
		every = time.Duration(timer.TimeInterval) * time.Minute
	}

	_, err = c.cl.ScheduleClient().Create(
		ctx,
		client.ScheduleOptions{
			ID: scheduleID,
			Spec: client.ScheduleSpec{
				Intervals: []client.ScheduleIntervalSpec{
					{
						Every: every,
					},
				},
			},
			Action: &client.ScheduleWorkflowAction{
				ID:        workflowID,
				Workflow:  c.Flow,
				TaskQueue: shared.TimersWorkerTaskQueueName,
				Args:      []interface{}{timer},
				Memo:      map[string]interface{}{"lastResponse": 0},
				RetryPolicy: &temporal.RetryPolicy{
					InitialInterval:        time.Second,
					BackoffCoefficient:     2.0,
					MaximumInterval:        time.Second * 100,
					MaximumAttempts:        3,
					NonRetryableErrorTypes: []string{},
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Workflow) RemoveTimer(ctx context.Context, timerId string) error {
	handle := c.cl.ScheduleClient().GetHandle(ctx, timerId)
	return handle.Delete(ctx)
}
