package main

import (
	"errors"
	"fmt"
	"github.com/begmaroman/tasq"
	"log"
	"sync"
	"sync/atomic"
)

var (
	ctr         = int32(0)
	handleError = errors.New("task handle error")
)

type taskHandler struct {
	err error
}

func newTaskHandler(err error) *taskHandler {
	return &taskHandler{
		err: err,
	}
}

func (p *taskHandler) Do() error {
	fmt.Println("global counter:", atomic.AddInt32(&ctr, 1))
	return p.err
}

func main() {
	var wg sync.WaitGroup

	tq := tasq.New()
	tq.TaskDone = func(id int64) {
		fmt.Printf("task %d done\n", id)
		wg.Done()
	}
	tq.TaskFailed = func(id int64, err error) {
		log.Printf("task %d failed with err %q\n", id, err)
		wg.Done()
	}
	tq.Start()

	// without error
	for i := 0; i < 100; i++ {
		wg.Add(1)

		log.Print("added task with id:", tq.Enqueue(newTaskHandler(nil)))
	}

	// with custom error
	for i := 0; i < 100; i++ {
		wg.Add(1)

		log.Print("added task (custom error) with id:", tq.Enqueue(newTaskHandler(handleError)))
	}

	// with tasq retry error
	for i := 0; i < 100; i++ {
		wg.Add(1)

		log.Print("added task (retry error) with id:", tq.Enqueue(newTaskHandler(tasq.ErrRetryTask)))
	}

	tq.Close()

	wg.Wait()

	fmt.Println("final global counter:", ctr)
}
