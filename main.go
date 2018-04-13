package main

import (
	"fmt"
	"log"
	"time"

	"github.com/marcusolsson/tui-go"
)

type post struct {
	username string
	message  string
	time     string
}

var posts = []post{
	{username: "john", message: "hi, what's up?", time: "14:41"},
	{username: "jane", message: "not much", time: "14:43"},
}

var ui tui.UI

var sectors [3]*tui.List

var projects []*Project

var sectorIndex = 0

func main() {

	Setup()

	// Little hack to gather user input
	// from <key>, not from Ctrl-<key>
	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	sectors[0] = tui.NewList()
	sectors[0].SetFocused(true)
	sectors[0].OnSelectionChanged(ProjectChanged)
	projectsBox := tui.NewVBox(sectors[0])
	projectsBox.SetTitle("Projects")
	projectsBox.SetBorder(true)
	projectsBox.SetFocused(true)
	for _, val := range projects {
		sectors[0].AddItems(val.Name)
	}

	sectors[1] = tui.NewList()
	sectors[1].OnSelectionChanged(TaskChanged)
	tasksBox := tui.NewVBox(sectors[1])
	tasksBox.SetTitle("Tasks")
	tasksBox.SetBorder(true)

	sectors[2] = tui.NewList()
	entryBox := tui.NewVBox(sectors[2])
	entryBox.SetTitle("Entries")
	entryBox.SetBorder(true)

	sectorsBox := tui.NewHBox(projectsBox, tasksBox, entryBox)

	root := tui.NewVBox(sectorsBox, input)

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

	sectors[0].Select(0)

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func Setup() {
	entries := []*Entry{{"Adicionar referências"}, {"Ler artigos"}}
	entries2 := []*Entry{{"Corrigir posição do mapa"}, {"Implementar heróis"}}

	tasks := []*Task{{"Escrever anexo II", entries}, {"Implementar fases", entries2}}

	projects = []*Project{{"Code Overlord", tasks}}
}

func ProjectChanged(l *tui.List) {

	if l.Selected() < 0 {
		return
	}

	// Clear tasks
	sectors[1].RemoveItems()

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

	// Foreach entry in the current task, in the current project...
	for _, val := range projects[sectors[0].Selected()].Tasks[l.Selected()].Entries {
		sectors[2].AddItems(val.Name)
	}

	sectors[2].Select(0)
}

func command(s string) bool {
	switch s {
	case "q":
		ui.Quit()
	case "l":
		nextSector()
	case "h":
		prevSector()
	case "j":
		fallthrough
	case "k":
		return true
	default:
		return false
	}
	return true
}

func nextSector() {
	sectors[sectorIndex].SetFocused(false)

	sectorIndex++

	if sectorIndex >= len(sectors) {
		sectorIndex = len(sectors) - 1
	}

	sectors[sectorIndex].SetFocused(true)
}

func prevSector() {
	sectors[sectorIndex].SetFocused(false)

	sectorIndex--

	if sectorIndex < 0 {
		sectorIndex = 0
	}

	sectors[sectorIndex].SetFocused(true)
}

func sample() {
	sidebar := tui.NewVBox(
		tui.NewLabel("CHANNELS"),
		tui.NewLabel("general"),
		tui.NewLabel("random"),
		tui.NewLabel(""),
		tui.NewLabel("DIRECT MESSAGES"),
		tui.NewLabel("slackbot"),
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)

	history := tui.NewVBox()

	for _, m := range posts {
		history.Append(tui.NewHBox(
			tui.NewLabel(m.time),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", m.username))),
			tui.NewLabel(m.message),
			tui.NewSpacer(),
		))
	}

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", "john"))),
			tui.NewLabel(e.Text()),
			tui.NewSpacer(),
		))
		input.SetText("")
	})

	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
