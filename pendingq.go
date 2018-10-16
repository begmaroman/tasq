package tasq

import "sync"

type pendingQ struct {
	sync.Mutex
	queue []*iTask
}

func newPendingQ(size int) *pendingQ {
	return &pendingQ{
		queue: make([]*iTask, 0, size),
	}
}

func (q *pendingQ) enq(it *iTask) {
	q.Lock()
	q.queue = append(q.queue, it)
	q.Unlock()
}

func (q *pendingQ) deq() *iTask {
	q.Lock()
	defer q.Unlock()

	if len(q.queue) == 0 {
		return nil
	}
	it := q.queue[0]
	q.queue = q.queue[1:]
	return it
}
