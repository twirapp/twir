package manager

import (
	"sync"
	"testing"
	"time"
)

func TestTimerWithTickLockSerializesTicks(t *testing.T) {
	timer := Timer{}

	started := make(chan struct{})
	release := make(chan struct{})
	secondEntered := make(chan struct{})

	go func() {
		timer.withTickLock(func() {
			close(started)
			<-release
		})
	}()

	<-started

	var secondTick sync.WaitGroup
	secondTick.Add(1)
	go func() {
		defer secondTick.Done()
		timer.withTickLock(func() {
			close(secondEntered)
		})
	}()

	select {
	case <-secondEntered:
		t.Fatal("second tick must wait for the first tick")
	case <-time.After(20 * time.Millisecond):
	}

	close(release)
	secondTick.Wait()

	select {
	case <-secondEntered:
	default:
		t.Fatal("second tick must run after the first tick finishes")
	}
}
