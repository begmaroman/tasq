package tasq

import (
	"sync/atomic"
)

const (
	Pending byte = iota
	InProgress
	Done
	Failed
)

var (
	// Count of workers (go-routines) which will be run
	WorkersPollSize = 10
	// Count of retry when task processing is fail
	TaskMaxRetry = 3
)

type Task interface {
	Do() error
}

type TasQ struct {
	lastInc       int64
	queue         chan *iTask
	pending       *pendingQ
	tasksMaxRetry int
	workersCount  int

	TaskDone   func(int64)
	TaskFailed func(int64, error)
}

func New() *TasQ {
	return &TasQ{
		workersCount:  WorkersPollSize,
		tasksMaxRetry: TaskMaxRetry,
		queue:         make(chan *iTask, WorkersPollSize),
		pending:       newPendingQ(WorkersPollSize),
	}
}

func (t *TasQ) Enqueue(task Task) int64 {
	if task == nil {
		return -1
	}

	it := newITask(atomic.AddInt64(&t.lastInc, 1), Pending, task)
	select {
	case t.queue <- it:
		// successfully sent to the workers
	default:
		// if we can't send task to directly to the workers
		// we add it to pending queue
		t.pending.enq(it)
	}

	return it.id
}

func (t *TasQ) Start() error {
	// run process workers
	for i := 0; i < t.workersCount; i++ {
		// each worker will make task.Do
		go func(workerID int) {
			for task := range t.queue {
				for task != nil {
					t.process(task)
					task = t.pending.deq()
				}
			}
		}(i)
	}
	return nil
}

func (t *TasQ) Close() {
	close(t.queue)
}

func (t *TasQ) process(it *iTask) {
	it.state = InProgress
	var try int
	for {
		if err := it.task.Do(); err != nil {
			if err == ErrRetryTask {
				try++
				if try < t.tasksMaxRetry {
					continue
				} else {
					err = ErrMaxRetry
				}
			}

			it.state = Failed
			if t.TaskFailed != nil {
				t.TaskFailed(it.id, err)
			}

			break
		}

		it.state = Done
		if t.TaskDone != nil {
			t.TaskDone(it.id)
		}

		break
	}
}
