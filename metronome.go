package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	metronomeDisplay    = ""
	bpm                 = 0
	done                = make(chan struct{})
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	currentBarIterator  = 0
	chordsPerBar        = []string{"G", "G", "G", "G", "D", "D", "D", "D"}
)

func main() {
	configureLog()
	go tickMetronome()
	if err := tea.NewProgram(initialModel()).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}

// write log to file
func configureLog() {
	f, err := os.OpenFile("metronome.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

type model struct {
	focusIndex int
	bpmInput   textinput.Model
	cursorMode textinput.CursorMode
}

func initialModel() model {
	m := model{
		bpmInput: textinput.Model{},
	}

	var t textinput.Model
	t = textinput.New()
	t.CursorStyle = cursorStyle
	t.CharLimit = 32

	t.Placeholder = "Beats Per Minute"
	t.Focus()
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle

	m.bpmInput = t

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > textinput.CursorHide {
				m.cursorMode = textinput.CursorBlink
			}
			cmds := m.bpmInput.SetCursorMode(m.cursorMode)
			return m, cmds
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds tea.Cmd

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	// for i := range m.bpmInput {
	m.bpmInput, cmds = m.bpmInput.Update(msg)
	// close(done)
	setMetronomeBpm(m.bpmInput.Value())

	return cmds
}

// Updates metronome view each bpm
func tickMetronome() {
	var ti = 1
	for {
		log.Printf("iteratting and bpm is %d\n", bpm)
		if !(bpm > 0) {
			continue
		}
		milliseconds := 60000 / bpm
		select {
		case <-done:
		case <-time.After(time.Duration(milliseconds) * time.Millisecond):
			log.Println(ti)
			ti += 1
			spaceToPrepend := strings.Repeat(" ", currentBarIterator)
			currentBarIterator = (currentBarIterator + 1) % len(chordsPerBar)
			metronomeDisplay = spaceToPrepend + chordsPerBar[currentBarIterator]
		}
	}
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(m.bpmInput.View())
	// for i := range m.bpmInput {
	// 	b.WriteString(m.bpmInput[i].View())
	// 	if i < len(m.bpmInput)-1 {
	// 		b.WriteRune('\n')
	// 	}
	// }

	// button := &blurredButton
	// if m.focusIndex == 1 {
	// 	button = &focusedButton
	// }
	// fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	fmt.Fprintf(&b, "\n\n%s\n\n", metronomeDisplay)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

func setMetronomeBpm(bpmInput string) {
	intVar, err := strconv.Atoi(bpmInput)
	if err != nil {
		// println(
		// log.Fatal("cant convert")
		log.Println("bad number in bpm")
		return
	}
	bpm = intVar
	// tickMetronome(intVar, m)
	// return strconv.Itoa(intVar)
}
