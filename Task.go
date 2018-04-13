package main

type Task struct {
	Name    string
	Entries []*Entry
}

func NewTask(name string, entries []*Entry) *Task {

	if entries == nil {
		entries = []*Entry{}
	}

	return &Task{name, entries}
}
