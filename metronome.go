package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	chordsPerBar        = []string{"G", "G", "G", "G", "D", "D", "D", "D"}
	defaultMetronome    = Metronome{
		Frames: getFrames(chordsPerBar),
		FPS:    bpm2bps(1),
	}
)

func getFrames(bars []string) []string {
	toReturn := []string{}
	space := ""
	for _, v := range bars {
		toReturn = append(toReturn, space+v)
		space = " " + space
	}
	return toReturn
}

// bpm2bps get time.Duration of metronome tick for given BPM
func bpm2bps(bpm int) time.Duration {
	// https://hextobinary.com/unit/frequency/from/bpm/to/fps
	const bpmConversion float64 = .016666666666667
	return time.Duration(float64(time.Second) / (float64(bpm) * bpmConversion))
}

type Metronome struct {
	Frames []string
	FPS    time.Duration
}

func main() {
	configureLog()
	// go tickMetronome()
	if err := tea.NewProgram(newModel()).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}

type Model struct {
	bpmInput   textinput.Model
	Style      lipgloss.Style
	metronome  Metronome
	frame      int
	id         int
	tag        int
	cursorMode textinput.CursorMode
}

func (m Model) ID() int {
	return m.id
}

type TickMsg struct {
	Time time.Time
	tag  int
	ID   int
}

func newModel() Model {
	var t textinput.Model
	t = textinput.New()
	t.CursorStyle = cursorStyle
	t.CharLimit = 32

	t.Placeholder = "Beats Per Minute"
	t.Focus()
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle

	return Model{
		bpmInput:  t,
		metronome: defaultMetronome,
		id:        nextID(),
	}
}

var (
	lastID int
	idMtx  sync.Mutex
)

func nextID() int {
	lg("nextID()")
	idMtx.Lock()
	defer idMtx.Unlock()
	lastID++
	return lastID
}

func (m Model) Init() tea.Cmd {
	lg("model.Init()")
	return m.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	lg(fmt.Sprintf("m.Update().\nmsg: %+v", msg))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		lg("m.Update() tea.KeyMsg")
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > textinput.CursorHide {
				m.cursorMode = textinput.CursorBlink
			}
			cmd := m.bpmInput.SetCursorMode(m.cursorMode)
			return m, cmd
		default:
			cmd := m.updateInputs(msg)
			return m, cmd
		}
	case TickMsg:
		lg("tickMsg")
		if msg.ID > 0 && msg.ID != m.id {
			return m, nil
		}
		if msg.tag > 0 && msg.tag != m.tag {
			return m, nil
		}
		m.frame++
		if m.frame >= len(m.metronome.Frames) {
			m.frame = 0
		}
		m.tag++
		return m, m.tick(m.id, m.tag)
	}
	return m, nil
}

func (m Model) Tick() tea.Msg {
	return TickMsg{
		Time: time.Now(),
		ID:   m.id,
		tag:  m.tag,
	}
}

func (m Model) tick(id, tag int) tea.Cmd {
	lg("m.tick()")
	return tea.Tick(m.metronome.FPS, func(t time.Time) tea.Msg {
		return TickMsg{
			Time: t,
			ID:   id,
			tag:  tag,
		}
	})
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	lg("model.updateInputs()")
	var cmd tea.Cmd

	m.bpmInput, cmd = m.bpmInput.Update(msg)

	bpmVal, err := getBpmFromString(m.bpmInput.Value())
	if err == nil {
		m.metronome.FPS = bpm2bps(bpmVal)
		cmd = m.tick(m.id, m.tag)
		return cmd
	}

	return cmd
}

func (m Model) View() string {
	lg("m.View()")
	var b strings.Builder

	b.WriteString(m.bpmInput.View())
	fmt.Fprintf(&b, "\n\n%s\n\n", m.Style.Render(m.metronome.Frames[m.frame]))

	return b.String()
}

func getBpmFromString(bpmInput string) (int, error) {
	intVar, err := strconv.Atoi(bpmInput)
	if err != nil {
		lg("bad number in bpm")
		return 0, err
	}
	return intVar, nil
}

func configureLog() {
	f, err := os.OpenFile("metronome.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

const debug = false

func lg(output string) {
	if debug {
		log.Println(output)
	}
}
