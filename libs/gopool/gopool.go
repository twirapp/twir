package gopool

import (
	"sync"
)

type worker struct {
	tasksChannel chan func()
	readyWg      *sync.WaitGroup
}

func (w *worker) run() {
	w.readyWg.Done()
	for task := range w.tasksChannel {
		task()
	}
}

type Pool struct {
	workers          []*worker
	latestUsedWorker int
	submitLock       *sync.Mutex
}

func (p *Pool) Submit(f func()) {
	p.submitLock.Lock()
	defer p.submitLock.Unlock()

	neededWorker := 0

	if p.latestUsedWorker != len(p.workers)-1 {
		neededWorker = p.latestUsedWorker + 1
	}

	w := p.workers[neededWorker]
	p.latestUsedWorker = neededWorker

	w.tasksChannel <- f
}

func NewPool(size int) *Pool {
	readyWg := &sync.WaitGroup{}

	pool := &Pool{
		workers:    make([]*worker, size),
		submitLock: &sync.Mutex{},
	}

	for i := 0; i < size; i++ {
		readyWg.Add(1)

		w := &worker{
			readyWg:      readyWg,
			tasksChannel: make(chan func()),
		}
		go w.run()

		pool.workers[i] = w
	}

	readyWg.Wait()

	return pool
}
