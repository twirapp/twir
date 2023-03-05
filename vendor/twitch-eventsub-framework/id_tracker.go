package eventsub_framework

import "context"

type IDTracker interface {
	// AddAndCheckIfDuplicate returns if the ID is a duplicate and an error.
	AddAndCheckIfDuplicate(ctx context.Context, id string) (bool, error)
}

// MapTracker uses an in-memory map to check if a notification ID is
// a duplicate.
type MapTracker struct {
	seen map[string]struct{}
}

func (m *MapTracker) AddAndCheckIfDuplicate(_ context.Context, id string) (bool, error) {
	_, ok := m.seen[id]
	if ok {
		return true, nil
	}

	m.seen[id] = struct{}{}
	return false, nil
}

// NewMapTracker creates a new MapTracker instance which uses an in-memory map
// to check if a notification ID is a duplicate.
func NewMapTracker() *MapTracker {
	return &MapTracker{
		seen: make(map[string]struct{}),
	}
}
