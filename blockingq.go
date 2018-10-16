package tasq

import "sync"

type blockingQ struct {
	sync.Mutex
	queue []*iTask
}

func newBlockingQ(size int) *blockingQ {
	return &blockingQ{
		queue: make([]*iTask, 0, size),
	}
}

func (q *blockingQ) enq(it *iTask) {
	q.Lock()
	q.queue = append(q.queue, it)
	q.Unlock()
}

func (q *blockingQ) deq() *iTask {
	q.Lock()
	defer q.Unlock()

	if len(q.queue) == 0 {
		return nil
	}
	it := q.queue[0]
	q.queue = q.queue[1:]
	return it
}
