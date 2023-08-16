package utils

import "sync"

type GoroutinesGroup struct {
	wg sync.WaitGroup
}

func NewGoroutinesGroup() *GoroutinesGroup {
	return &GoroutinesGroup{wg: sync.WaitGroup{}}
}

func (c *GoroutinesGroup) Go(f func()) {
	c.wg.Add(1)

	go func() {
		defer c.wg.Done()
		f()
	}()
}

func (c *GoroutinesGroup) Wait() {
	c.wg.Wait()
}
