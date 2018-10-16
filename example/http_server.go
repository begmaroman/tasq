package main

import (
	"log"
	"net/http"
	"time"

	"github.com/begmaroman/tasq"
)

var (
	tq *tasq.TasQ
)

type CustomTask struct {
	logData interface{}
}

func NewCustomTask(logData interface{}) *CustomTask {
	return &CustomTask{
		logData: logData,
	}
}

func (c *CustomTask) Do() error {
	// imitation some huge logic
	time.Sleep(time.Second)

	log.Println("request query parameters: ", c.logData)
	return nil
}

func init() {
	tq = tasq.New()
}

func main() {
	// start tasq background workers
	tq.Start()
	defer tq.Close()

	http.HandleFunc("/tasq", handler)
	log.Fatal(http.ListenAndServe(":7575", nil))
}

func handler(writer http.ResponseWriter, request *http.Request) {
	// add some task to background task queue
	// If a task requires some logic which not needs for client You can move this logic to background.
	tq.Enqueue(NewCustomTask(request.URL.Query().Encode()))
}
