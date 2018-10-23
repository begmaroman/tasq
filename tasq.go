package tasq

import (
	"sync/atomic"
)

// Process statuses
const (
	Pending byte = iota
	InProgress
	Done
	Failed
)

var (
	// Count of workers (go-routines) which will be run
	WorkersCount = 10
	// Count of retry when task processing is fail
	MaxRetry = 3
)

type TasQ struct {
	lastID       int64
	queue        chan *taskProcess
	pending      *pendingQ
	maxRetry     int
	workersCount int

	TaskDone   func(int64)
	TaskFailed func(int64, error)
}

func New() *TasQ {
	return &TasQ{
		workersCount: WorkersCount,
		maxRetry:     MaxRetry,
		queue:        make(chan *taskProcess, WorkersCount),
		pending:      newPendingQ(WorkersCount),
	}
}

func (t *TasQ) Enqueue(task Task) int64 {
	if task == nil {
		return -1
	}

	it := newTaskProcess(atomic.AddInt64(&t.lastID, 1), Pending, task)
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

func (t *TasQ) process(tp *taskProcess) {
	if tDone, ok := tp.task.(TaskDone); ok {
		defer tDone.Done()
	}

	tp.state = InProgress
	var try int
	for {
		if err := tp.task.Do(); err != nil {
			if err == ErrRetryTask {
				try++
				if try < t.maxRetry {
					continue
				} else {
					err = ErrMaxRetry
				}
			}

			tp.state = Failed
			if t.TaskFailed != nil {
				t.TaskFailed(tp.id, err)
			}

			break
		}

		tp.state = Done
		if t.TaskDone != nil {
			t.TaskDone(tp.id)
		}

		break
	}
}
