package ui

import (
	"fmt"
	"os"
	"strconv"

	"strings"
	"sync"
	"time"

	metrlog "github.com/andrewwillette/metronome/log"
	"github.com/andrewwillette/metronome/song"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// https://hextobinary.com/unit/frequency/from/bpm/to/fps
const bpmConversion float64 = .016666666666667

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle  = focusedStyle.Copy()
	chordsPerBar = []string{"G", "G", "G", "G", "D", "D", "D", "D"}
	lastID       int
	idMtx        sync.Mutex
)

func StartUi() {
	// songs := musicparse.GetDefaultSongs()
	model := newModel()
	if err := tea.NewProgram(model).Start(); err != nil {
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
	return time.Duration(float64(time.Second) / (float64(bpm) * bpmConversion))
}

type BaseModel struct {
	metronomeDisplay    string
	songs               []song.Song
	activeSong          song.Song
	bpmInputModel       textinput.Model
	bpmInputStyle       lipgloss.Style
	metronomeFrameStyle lipgloss.Style
	songFrames          []string
	fps                 time.Duration
	frame               int
	id                  int
	tag                 int
	cursorMode          textinput.CursorMode
}

func (m BaseModel) ID() int {
	return m.id
}

type TickMsg struct {
	Time time.Time
	tag  int
	ID   int
}

func newModel() *BaseModel {
	var t textinput.Model
	t = textinput.New()
	t.CursorStyle = cursorStyle
	t.CharLimit = 32

	t.Placeholder = "Beats Per Minute"
	t.Focus()
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle

	songs := song.GetSongsXdg()
	activeSong := songs[0]
	songFrames := song.GetSongFrames(activeSong)
	return &BaseModel{
		songs:               songs,
		activeSong:          activeSong,
		bpmInputModel:       t,
		songFrames:          songFrames,
		fps:                 bpm2bps(1),
		id:                  nextID(),
		bpmInputStyle:       lipgloss.NewStyle().BorderStyle(lipgloss.DoubleBorder()),
		metronomeFrameStyle: lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()),
	}
}

func nextID() int {
	idMtx.Lock()
	defer idMtx.Unlock()
	lastID++
	return lastID
}

func (m BaseModel) Init() tea.Cmd {
	return m.Tick
}

// var tickernum = 0

func (m BaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			cmd := m.bpmInputModel.SetCursorMode(m.cursorMode)
			return m, cmd
		default:
			// m.bpmUpdated <- struct{}{}
			cmd := m.updateInputs(msg)
			return m, cmd
		}
	case TickMsg:
		// tickernum++
		if msg.ID > 0 && msg.ID != m.id {
			return m, nil
		}
		if msg.tag > 0 && msg.tag != m.tag {
			return m, nil
		}
		m.frame++
		if m.frame >= len(m.songFrames) {
			m.frame = 0
		}
		m.tag++
		return m, m.tick(m.id, m.tag)
	}
	return m, nil
}

func (m BaseModel) Tick() tea.Msg {
	return TickMsg{
		Time: time.Now(),
		ID:   m.id,
		tag:  m.tag,
	}
}

func (m BaseModel) tick(id, tag int) tea.Cmd {
	return tea.Tick(m.fps, func(t time.Time) tea.Msg {
		return TickMsg{
			Time: t,
			ID:   id,
			tag:  tag,
		}
	})
}

func (m *BaseModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	m.bpmInputModel, cmd = m.bpmInputModel.Update(msg)

	bpmVal, err := getBpmFromString(m.bpmInputModel.Value())
	if err != nil {
		return cmd
	}
	m.fps = bpm2bps(bpmVal)
	cmd = m.tick(m.id, m.tag)

	return cmd
}

func (m BaseModel) View() string {
	var b strings.Builder
	b.WriteString(m.bpmInputStyle.Render(m.bpmInputModel.View()))
	b.WriteString(fmt.Sprintf("\n%s\n", m.activeSong.Title))
	b.WriteString(fmt.Sprintf("\n%s\n", m.metronomeFrameStyle.Render(m.songFrames[m.frame])))
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
