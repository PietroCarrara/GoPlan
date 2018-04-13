package main

type Entry struct {
	Name      string
	Completed bool
}

func NewEntry(name string) *Entry {

	return &Entry{name, false}
}

func (e Entry) String() string {
	done := " "
	if e.Completed {
		done = "X"
	}

	return "[" + done + "] " + e.Name
}
