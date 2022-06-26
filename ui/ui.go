package ui

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	metrlog "github.com/andrewwillette/metronome/log"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle  = focusedStyle.Copy()
	chordsPerBar = []string{"G", "G", "G", "G", "D", "D", "D", "D"}
)

func StartUi() {
	// songs := musicparse.GetDefaultSongs()
	if err := tea.NewProgram(newModel()).Start(); err != nil {
		metrlog.Lg(fmt.Sprintf("could not start program: %s\n", err))
		os.Exit(1)
	}
}

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

type MetronomeModel struct {
	bpmInput            textinput.Model
	bpmInputStyle       lipgloss.Style
	metronomeFrameStyle lipgloss.Style
	frames              []string
	fps                 time.Duration
	frame               int
	id                  int
	tag                 int
	cursorMode          textinput.CursorMode
}

func (m MetronomeModel) ID() int {
	return m.id
}

type TickMsg struct {
	Time time.Time
	tag  int
	ID   int
}

func newModel() MetronomeModel {
	var t textinput.Model
	t = textinput.New()
	t.CursorStyle = cursorStyle
	t.CharLimit = 32

	t.Placeholder = "Beats Per Minute"
	t.Focus()
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle

	return MetronomeModel{
		bpmInput:            t,
		frames:              getFrames(chordsPerBar),
		fps:                 bpm2bps(1),
		id:                  nextID(),
		bpmInputStyle:       lipgloss.NewStyle().BorderStyle(lipgloss.DoubleBorder()),
		metronomeFrameStyle: lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()),
	}
}

var (
	lastID int
	idMtx  sync.Mutex
)

func nextID() int {
	metrlog.Lg("nextID()")
	idMtx.Lock()
	defer idMtx.Unlock()
	lastID++
	return lastID
}

func (m MetronomeModel) Init() tea.Cmd {
	metrlog.Lg("model.Init()")
	return m.Tick
}

var tickernum = 0

func (m MetronomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	metrlog.Lg(fmt.Sprintf("m.Update().\nmsg: %+v", msg))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		metrlog.Lg("m.Update() tea.KeyMsg")
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
		tickernum++
		metrlog.Lg("tickMsg")
		if msg.ID > 0 && msg.ID != m.id {
			return m, nil
		}
		if msg.tag > 0 && msg.tag != m.tag {
			return m, nil
		}
		m.frame++
		if m.frame >= len(m.frames) {
			m.frame = 0
		}
		m.tag++
		return m, m.tick(m.id, m.tag)
	}
	return m, nil
}

func (m MetronomeModel) Tick() tea.Msg {
	return TickMsg{
		Time: time.Now(),
		ID:   m.id,
		tag:  m.tag,
	}
}

func (m MetronomeModel) tick(id, tag int) tea.Cmd {
	metrlog.Lg("m.tick()")
	return tea.Tick(m.fps, func(t time.Time) tea.Msg {
		return TickMsg{
			Time: t,
			ID:   id,
			tag:  tag,
		}
	})
}

func (m *MetronomeModel) updateInputs(msg tea.Msg) tea.Cmd {
	metrlog.Lg("model.updateInputs()")
	var cmd tea.Cmd

	m.bpmInput, cmd = m.bpmInput.Update(msg)

	bpmVal, err := getBpmFromString(m.bpmInput.Value())
	if err != nil {
		return cmd
	}
	m.fps = bpm2bps(bpmVal)
	cmd = m.tick(m.id, m.tag)

	return cmd
}

func (m MetronomeModel) View() string {
	metrlog.Lg("m.View()")
	var b strings.Builder

	b.WriteString(m.bpmInputStyle.Render(m.bpmInput.View()))
	b.WriteString(fmt.Sprintf("\n\n%s\n\n", m.metronomeFrameStyle.Render(m.frames[m.frame])))
	b.WriteString(fmt.Sprintf("\n\n%d\n\n", tickernum))

	return b.String()
}

func getBpmFromString(bpmInput string) (int, error) {
	intVar, err := strconv.Atoi(bpmInput)
	if err != nil {
		metrlog.Lg("bad number in bpm")
		return 0, err
	}
	return intVar, nil
}
