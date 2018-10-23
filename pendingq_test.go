package tasq

import (
	"reflect"
	"testing"
)

func TestPendingQ_CheckQueueCap(t *testing.T) {
	size := 10
	p := newPendingQ(size)

	if cap(p.queue) != size {
		t.Errorf("invalid queue capacity: expected %d given %d", size, cap(p.queue))
	}
}

func TestPendingQ_Enq(t *testing.T) {
	ts := newTaskProcess(1, 2, nil)
	p := newPendingQ(10)
	p.enq(ts)

	if len(p.queue) != 1 {
		t.Errorf("invalid queue length: expected 1 given %d", cap(p.queue))
	}

	if !reflect.DeepEqual(p.queue[0], ts) {
		t.Errorf("invalid task value: expected %v given %v", ts, p.queue[0])
	}
}

func TestPendingQ_Deq(t *testing.T) {
	ts := newTaskProcess(1, 2, nil)
	p := newPendingQ(10)
	p.enq(ts)
	dTs := p.deq()

	if len(p.queue) != 0 {
		t.Errorf("invalid queue length: expected 0 given %d", cap(p.queue))
	}

	if !reflect.DeepEqual(dTs, ts) {
		t.Errorf("invalid task value: expected %v given %v", ts, dTs)
	}
}
