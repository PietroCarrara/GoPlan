package main

type Project struct {
	Name  string
	Tasks []*Task
}

func NewProject(name string, tasks []*Task) *Project {

	if tasks == nil {
		tasks = []*Task{}
	}

	return &Project{name, tasks}
}
