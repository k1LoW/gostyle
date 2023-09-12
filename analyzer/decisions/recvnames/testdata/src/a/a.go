package a

type TaskList struct {
	tasks []string
}

func (tasklist *TaskList) Do(task string) error { // want "gostyle.recvnames"
	return nil
}

func Do(task string) error {
	return nil
}

func (h *TaskList) Hello(task string) error { // want "gostyle.recvnames"
	return nil
}

func (l TaskList) World(task string) error {
	return nil
}
