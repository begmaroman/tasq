package tasq

type iTask struct {
	id    int64
	state byte
	task  Task
}

func newITask(id int64, state byte, task Task) *iTask {
	return &iTask{
		id:    id,
		state: state,
		task:  task,
	}
}
