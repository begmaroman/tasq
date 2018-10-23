package tasq

import "sync"

type pendingQ struct {
	sync.Mutex
	queue []*taskProcess
}

func newPendingQ(size int) *pendingQ {
	return &pendingQ{
		queue: make([]*taskProcess, 0, size),
	}
}

func (q *pendingQ) enq(it *taskProcess) {
	q.Lock()
	q.queue = append(q.queue, it)
	q.Unlock()
}

func (q *pendingQ) deq() *taskProcess {
	q.Lock()
	defer q.Unlock()

	if len(q.queue) == 0 {
		return nil
	}
	it := q.queue[0]
	q.queue = q.queue[1:]
	return it
}
