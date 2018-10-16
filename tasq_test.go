package tasq

import (
	"errors"
	"sync"
	"testing"
)

type testTask struct {
	isDone    bool
	err       error
	execCount int
}

func newTestTask(err error) *testTask {
	return &testTask{
		err: err,
	}
}

func (t *testTask) Do() error {
	t.execCount++
	t.isDone = true
	return t.err
}

func TestTasQ_Enqueue(t *testing.T) {
	tsk := newTestTask(nil)

	var wg sync.WaitGroup
	tq := New()
	tq.TaskFailed = func(i int64, e error) {
		wg.Done()
	}
	tq.TaskDone = func(i int64) {
		wg.Done()
	}
	tq.Start()
	wg.Add(1)
	tq.Enqueue(tsk)
	tq.Close()

	wg.Wait()

	if !tsk.isDone {
		t.Error("task not completed")
	}
}

func TestTasQ_Close(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("no panic when send on closed channel")
		}
	}()

	tsk := newTestTask(nil)

	var wg sync.WaitGroup
	tq := New()
	tq.TaskFailed = func(i int64, e error) {
		wg.Done()
	}
	tq.TaskDone = func(i int64) {
		wg.Done()
	}
	tq.Start()
	tq.Close()

	wg.Add(1)
	tq.Enqueue(tsk)
	wg.Wait()
}

func TestTasQ_CheckMaxRetry(t *testing.T) {
	tsk := newTestTask(ErrRetryTask)

	var wg sync.WaitGroup
	tq := New()
	tq.TaskFailed = func(i int64, e error) {
		wg.Done()
	}
	tq.TaskDone = func(i int64) {
		wg.Done()
	}
	tq.Start()
	wg.Add(1)
	tq.Enqueue(tsk)
	tq.Close()

	wg.Wait()

	if tsk.execCount != TaskMaxRetry {
		t.Errorf("wrong retry count: expected %d given %d", TaskMaxRetry, tsk.execCount)
	}
}

func TestTasQ_CheckError(t *testing.T) {
	tsk := newTestTask(errors.New("test error"))
	var isFailedCall bool
	var wg sync.WaitGroup
	tq := New()
	tq.TaskFailed = func(i int64, e error) {
		isFailedCall = true
		wg.Done()
	}
	tq.TaskDone = func(i int64) {
		wg.Done()
	}
	tq.Start()
	wg.Add(1)
	tq.Enqueue(tsk)
	tq.Close()

	wg.Wait()

	if tsk.execCount != 1 {
		t.Errorf("wrong retry count: expected %d given %d", 1, tsk.execCount)
	}

	if !isFailedCall {
		t.Error("expected call TaskFailed function")
	}
}

func BenchmarkTasQ_Enqueue_SuccessDone(b *testing.B) {
	var wg sync.WaitGroup
	wg.Add(b.N)

	tq := New()
	tq.TaskFailed = func(i int64, e error) {
		wg.Done()
	}
	tq.TaskDone = func(i int64) {
		wg.Done()
	}
	tq.Start()

	for i := 0; i < b.N; i++ {
		tq.Enqueue(newTestTask(nil))
	}

	tq.Close()
	wg.Wait()
}

func BenchmarkTasQ_Enqueue_CustomErrorDone(b *testing.B) {
	var wg sync.WaitGroup
	wg.Add(b.N)

	tq := New()
	tq.TaskFailed = func(i int64, e error) {
		wg.Done()
	}
	tq.TaskDone = func(i int64) {
		wg.Done()
	}
	tq.Start()

	for i := 0; i < b.N; i++ {
		tq.Enqueue(newTestTask(errors.New("test error")))
	}

	tq.Close()
	wg.Wait()
}

func BenchmarkTasQ_Enqueue_RetryError3(b *testing.B) {
	var wg sync.WaitGroup
	wg.Add(b.N)

	tq := New()
	tq.TaskFailed = func(i int64, e error) {
		wg.Done()
	}
	tq.TaskDone = func(i int64) {
		wg.Done()
	}
	tq.Start()

	for i := 0; i < b.N; i++ {
		tq.Enqueue(newTestTask(ErrRetryTask))
	}

	tq.Close()
	wg.Wait()
}
