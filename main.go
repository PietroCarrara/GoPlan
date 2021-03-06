package main

import (
	"log"

	"github.com/marcusolsson/tui-go"
)

// Main ui
var ui tui.UI

// Screen sectors and their lists
var boxes [3]*tui.Box
var sectors [3]*tui.List

// Sector names
var names [3]string = [3]string{"Projects", "Tasks", "Entries"}

// Projects (which includes tasks, that includes entries)
var projects []*Project

// Which sector (Projects, Tasks...) is selected
var sectorIndex = 0

// Label before input box
var inputLabel *tui.Label

type Mode int

const (
	Insert Mode = iota
	Normal
)

var currentMode = Normal

var focusedStyle = tui.Style{
	Bold: tui.DecorationOn,
}

func main() {

	Setup()

	// Little hack to gather user input
	// from <key>, not from Ctrl-<key>
	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	input.OnSubmit(Input)

	inputLabel = tui.NewLabel("")
	// inputLabel.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(inputLabel, input)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	sectors[0] = tui.NewList()
	sectors[0].SetFocused(true)
	sectors[0].OnSelectionChanged(ProjectChanged)
	boxes[0] = tui.NewVBox(sectors[0])
	boxes[0].SetTitle("Projects")
	boxes[0].SetBorder(true)
	boxes[0].SetFocused(true)
	for _, val := range projects {
		sectors[0].AddItems(val.Name)
	}

	sectors[1] = tui.NewList()
	sectors[1].OnSelectionChanged(TaskChanged)
	boxes[1] = tui.NewVBox(sectors[1])
	boxes[1].SetTitle("Tasks")
	boxes[1].SetBorder(true)

	sectors[2] = tui.NewList()
	boxes[2] = tui.NewVBox(sectors[2])
	boxes[2].SetTitle("Entries")
	boxes[2].SetBorder(true)

	sectorsBox := tui.NewHBox(boxes[0], boxes[1], boxes[2])

	root := tui.NewVBox(sectorsBox, inputBox)

	var err error

	ui, err = tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	input.OnChanged(func(e *tui.Entry) {
		if command(e.Text()) {
			e.SetText("")
		}
	})

	ui.SetKeybinding("Esc", func() { input.SetText("") })

	if len(projects) > 0 {
		sectors[0].Select(0)
	}

	prevSector()

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func Setup() {
	projects = LoadFile()
}

func ProjectChanged(l *tui.List) {

	if l.Selected() < 0 {
		return
	}

	// Clear tasks
	sectors[1].RemoveItems()

	if len(projects[l.Selected()].Tasks) <= 0 {

		for _, sec := range sectors[2:] {
			sec.RemoveItems()
		}

		return
	}

	// Foreach task in the current project...
	for _, val := range projects[l.Selected()].Tasks {
		sectors[1].AddItems(val.Name)
	}

	sectors[1].Select(0)
}

func TaskChanged(l *tui.List) {

	if l.Selected() < 0 {
		return
	}

	// Clear entries
	sectors[2].RemoveItems()

	if len(projects[sectors[0].Selected()].Tasks[l.Selected()].Entries) <= 0 {

		for _, sec := range sectors[3:] {
			sec.RemoveItems()
		}

		return
	}

	// Foreach entry in the current task, in the current project...
	for _, val := range projects[sectors[0].Selected()].Tasks[l.Selected()].Entries {
		sectors[2].AddItems(val.String())
	}

	sectors[2].Select(0)
}

var inChan = make(chan string)

func Input(e *tui.Entry) {

	go input(e.Text())

	e.SetText("")
}

func input(text string) {
	if currentMode != Insert {
		inChan <- ""
	} else {
		inChan <- text
	}
}

// Runs the command 's', and returns true
// if 's' was a valid command string
func command(s string) bool {

	// If we are not in Normal, don't do anything
	if currentMode != Normal {
		return false
	}

	switch s {
	case "a":
		add()
	case "A":
		nextSector(true)
		add()
	case "q":
		SaveFile(projects)
		ui.Quit()
	case "x":
		complete()
	case "l":
		nextSector(false)
	case "h":
		prevSector()
	case "j":
		fallthrough
	case "k":
		return true
	default:
		return false
	}

	// If we're here, we didn't got into the
	// default case, so a command has been run
	return true
}

func add() {

	currentMode = Insert

	switch sectorIndex {
	case 0:
		inputLabel.SetText("Project Name: ")
		go addProject()
	case 1:
		inputLabel.SetText("Task Name: ")
		go addTask()
	case 2:
		inputLabel.SetText("Entry Name: ")
		go addEntry()
	}
}

func addProject() {

	proj := NewProject(<-inChan, nil)

	inputLabel.SetText("")
	currentMode = Normal

	if proj.Name == "" {
		return
	}

	projects = append(projects, proj)

	sectors[0].AddItems(proj.Name)

	sectors[0].Select(sectors[0].Length() - 1)
}

func addTask() {

	task := NewTask(<-inChan, nil)

	inputLabel.SetText("")
	currentMode = Normal

	if task.Name == "" {
		return
	}

	if sectors[0].Selected() < 0 {
		return
	}

	parent := projects[sectors[0].Selected()]

	parent.Tasks = append(parent.Tasks, task)

	sectors[1].AddItems(task.Name)

	sectors[1].Select(sectors[1].Length() - 1)
}

func addEntry() {

	entry := NewEntry(<-inChan)

	inputLabel.SetText("")
	currentMode = Normal

	if entry.Name == "" {
		return
	}

	parent := projects[sectors[0].Selected()].Tasks[sectors[1].Selected()]

	parent.Entries = append(parent.Entries, entry)

	sectors[2].AddItems(entry.String())

	sectors[2].Select(sectors[2].Length() - 1)
}

func complete() {
	index1, index2, index3 := sectors[0].Selected(), sectors[1].Selected(), sectors[2].Selected()

	if index1 < 0 || index2 < 0 || index3 < 0 {
		return
	}

	entries := projects[index1].Tasks[index2].Entries

	entries[index3].Completed = !entries[index3].Completed

	for i := index3; i < sectors[2].Length(); i++ {
		sectors[2].RemoveItem(index3)

		sectors[2].AddItems(entries[i].String())
	}

	sectors[2].Select(index3)
}

func nextSector(force bool) {
	sectors[sectorIndex].SetFocused(false)
	boxes[sectorIndex].SetTitle(names[sectorIndex])

	newIndex := sectorIndex + 1

	if newIndex >= len(sectors) {
		newIndex = len(sectors) - 1
	}

	// Only advance if there are things there
	// or we have been forced to
	if sectors[newIndex].Length() > 0 || force {
		sectorIndex = newIndex
	}

	sectors[sectorIndex].SetFocused(true)
	boxes[sectorIndex].SetTitle("!!" + names[sectorIndex] + "!!")
}

func prevSector() {
	sectors[sectorIndex].SetFocused(false)
	boxes[sectorIndex].SetTitle(names[sectorIndex])

	sectorIndex--

	if sectorIndex < 0 {
		sectorIndex = 0
	}

	sectors[sectorIndex].SetFocused(true)
	boxes[sectorIndex].SetTitle("!!" + names[sectorIndex] + "!!")
}
