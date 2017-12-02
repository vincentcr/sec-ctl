package main

import (
	"sync"
)

type workQueueFunc func(o interface{}) error

type workQueue struct {
	in chan interface{}

	quit chan struct{}

	tasks    []interface{}
	taskLock *sync.Cond
	worker   workQueueFunc
}

func newWorkQueue(worker workQueueFunc) *workQueue {
	return &workQueue{
		in:       make(chan interface{}, 128),
		quit:     make(chan struct{}),
		tasks:    make([]interface{}, 0, 128),
		taskLock: sync.NewCond(&sync.Mutex{}),
		worker:   worker,
	}
}

func (q *workQueue) enqueue(o interface{}) {
	q.taskLock.L.Lock()
	defer q.taskLock.L.Unlock()
	q.tasks = append(q.tasks, o)
	q.taskLock.Signal()
}

func (q *workQueue) start() {
	go q.consumeLoop()
}

func (q *workQueue) consumeLoop() {
	q.drain()
	for {
		q.taskLock.L.Lock()
		for len(q.tasks) == 0 {
			q.taskLock.Wait()
		}
		q.taskLock.L.Unlock()
		q.drain()
	}
}

func (q *workQueue) drain() {

	for len(q.tasks) > 0 {
		task := q.tasks[0]
		if err := q.worker(task); err != nil {
			q.handleTaskError(err)
		} else {
			q.taskLock.L.Lock()
			q.tasks = q.tasks[1:]
			q.taskLock.L.Unlock()
		}
	}
}

func (q *workQueue) handleTaskError(err error) {
	logger.Panicln(err)
}
