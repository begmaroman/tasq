package tasq

type taskProcess struct {
	id    int64
	state byte
	task  Task
}

func newTaskProcess(id int64, state byte, task Task) *taskProcess {
	return &taskProcess{
		id:    id,
		state: state,
		task:  task,
	}
}
