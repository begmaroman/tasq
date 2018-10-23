package tasq

type Task interface {
	Do() error
}

type TaskDone interface {
	Done()
}
