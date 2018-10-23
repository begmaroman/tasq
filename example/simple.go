package main

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/begmaroman/tasq"
)

var counter = int32(0)

type printer struct{}

func (p *printer) Do() error {
	fmt.Println("global counter:", atomic.AddInt32(&counter, 1))
	return nil
}

func (p *printer) Done() {
	fmt.Println("done function")
}

func main() {
	tq := tasq.New()
	tq.Start()
	for i := 0; i < 5; i++ {
		log.Print("added task with id:", tq.Enqueue(&printer{}))
	}
	tq.Close()
	time.Sleep(time.Second * 10)
}
