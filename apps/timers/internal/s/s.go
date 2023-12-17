package s

import (
	"context"

	"github.com/satont/twir/apps/timers/internal/repositories/timers"
	"github.com/satont/twir/apps/timers/internal/workflow"
)

func New(repository timers.Repository, w *workflow.Workflow) {
	timers, err := repository.GetAll()
	if err != nil {
		return
	}

	for _, timer := range timers {
		w.AddTimer(context.TODO(), timer.ID)
	}
}
