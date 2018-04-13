package main

type Entry struct {
	Name string
}

func NewEntry(name string) *Entry {

	return &Entry{name}
}
