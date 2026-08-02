package deploy

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/swarm"
)

type fakeServiceInspector struct {
	services []swarm.Service
	err      error
	index    int
}

func (f *fakeServiceInspector) ServiceInspectWithRaw(
	context.Context,
	string,
	swarm.ServiceInspectOptions,
) (swarm.Service, []byte, error) {
	if f.err != nil {
		return swarm.Service{}, nil, f.err
	}

	index := min(f.index, len(f.services)-1)
	service := f.services[index]
	f.index++
	return service, nil, nil
}

func serviceWithUpdate(version uint64, state swarm.UpdateState, message string) swarm.Service {
	return swarm.Service{
		Meta: swarm.Meta{Version: swarm.Version{Index: version}},
		UpdateStatus: &swarm.UpdateStatus{
			State:   state,
			Message: message,
		},
	}
}

func TestWaitForServiceUpdateCompletes(t *testing.T) {
	inspector := &fakeServiceInspector{services: []swarm.Service{
		serviceWithUpdate(10, swarm.UpdateStateCompleted, "previous update"),
		serviceWithUpdate(11, swarm.UpdateStateUpdating, "update in progress"),
		serviceWithUpdate(11, swarm.UpdateStateCompleted, "update completed"),
	}}

	err := waitForServiceUpdate(t.Context(), inspector, "service-id", 10, time.Millisecond)
	if err != nil {
		t.Fatalf("waitForServiceUpdate() error = %v", err)
	}
}

func TestWaitForServiceUpdateReportsTerminalFailure(t *testing.T) {
	tests := []struct {
		name      string
		state     swarm.UpdateState
		wantError string
	}{
		{name: "paused", state: swarm.UpdateStatePaused, wantError: "update paused"},
		{name: "rollback paused", state: swarm.UpdateStateRollbackPaused, wantError: "rollback paused"},
		{name: "rolled back", state: swarm.UpdateStateRollbackCompleted, wantError: "update rolled back"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inspector := &fakeServiceInspector{services: []swarm.Service{
				serviceWithUpdate(11, tt.state, "task failed"),
			}}

			err := waitForServiceUpdate(t.Context(), inspector, "service-id", 10, time.Millisecond)
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("waitForServiceUpdate() error = %v, want containing %q", err, tt.wantError)
			}
		})
	}
}

func TestWaitForServiceUpdateReportsInspectFailure(t *testing.T) {
	inspector := &fakeServiceInspector{err: errors.New("daemon unavailable")}

	err := waitForServiceUpdate(t.Context(), inspector, "service-id", 10, time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "inspect swarm service update") {
		t.Fatalf("waitForServiceUpdate() error = %v", err)
	}
}

func TestWaitForServiceUpdateHonorsContext(t *testing.T) {
	inspector := &fakeServiceInspector{services: []swarm.Service{
		serviceWithUpdate(11, swarm.UpdateStateUpdating, "update in progress"),
	}}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer cancel()

	err := waitForServiceUpdate(ctx, inspector, "service-id", 10, time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("waitForServiceUpdate() error = %v, want context deadline exceeded", err)
	}
}
