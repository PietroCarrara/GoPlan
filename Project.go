package main

import (
	"encoding/json"
	"io"
)

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

func SaveBatch(p []*Project, w io.Writer) error {

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	return enc.Encode(p)
}

func LoadBatch(r io.Reader) ([]*Project, error) {

	dec := json.NewDecoder(r)

	projs := []*Project{}

	err := dec.Decode(&projs)

	return projs, err
}
