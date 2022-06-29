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

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle  = focusedStyle.Copy()
	chordsPerBar = []string{"G", "G", "G", "G", "D", "D", "D", "D"}
)

func StartUi() {
	// songs := musicparse.GetDefaultSongs()
	model := newModel()
	go model.manageMetronomeDisplay()
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
	// https://hextobinary.com/unit/frequency/from/bpm/to/fps
	const bpmConversion float64 = .016666666666667
	return time.Duration(float64(time.Second) / (float64(bpm) * bpmConversion))
}

type BaseModel struct {
	bpmUpdated          chan struct{}
	metronomeDisplay    string
	songs               []song.Song
	bpmInputModel       textinput.Model
	bpmInputStyle       lipgloss.Style
	metronomeFrameStyle lipgloss.Style
	frames              []string
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

	return &BaseModel{
		songs:               song.GetSongsXdg(),
		bpmInputModel:       t,
		bpmUpdated:          make(chan struct{}),
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
		if m.frame >= len(m.frames) {
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

func getSongFrames(song song.Song) []string {
	frames := []string{}
	if len(song.Sections.ASection) > 0 {
		section := song.Sections.ASection
		for _, bar := range section {
			for _, beat := range bar {
				frames = append(frames, beat)
			}
		}
	}
	return frames
}

func (m *BaseModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	m.bpmInputModel, cmd = m.bpmInputModel.Update(msg)

	bpmVal, err := getBpmFromString(m.bpmInputModel.Value())
	if err != nil {
		return cmd
	}
	// <-m.bpmUpdated
	m.fps = bpm2bps(bpmVal)
	cmd = m.tick(m.id, m.tag)

	return cmd
}

func (m *BaseModel) viewMetronomeDisplay() string {
	song := m.songs[0]
	toReturn := song.Sections.ASection[0][0]
	// for {
	// 	// toReturn = song.Sections.ASection[i]
	// 	i++
	// 	if i == m.frame {
	// 		break
	// 	}
	// }
	// toDisplay = song
	return toReturn
}

func getMaxLenOfSection(section [][]string) int {
	i := 0
	for _, v1 := range section {
		for range v1 {
			i++
		}
	}
	return i
}

// manageMetronomeDisplay set UI display to appropriate string
//
func (baseModel *BaseModel) manageMetronomeDisplay() {
	song := baseModel.songs[0]
	maxLen := getMaxLenOfSection(song.Sections.ASection)
	for {
		select {
		case <-baseModel.bpmUpdated:
			metrlog.Lg("bpm updated caught in manageMetronomeDisplay")
			// break
		default:
			toDisplayFrameIndex := baseModel.frame % maxLen
			i := 0
			for {
				for _, bar := range song.Sections.ASection {
					for _, beat := range bar {
						if i == toDisplayFrameIndex {
							metrlog.Lg("i == toDisplay")
							metrlog.Lg(beat)
							baseModel.metronomeDisplay = beat
							i = i % maxLen
						}
						i++
					}
				}
			}
		}
	}
}

func (m BaseModel) View() string {
	var b strings.Builder

	b.WriteString(m.bpmInputStyle.Render(m.bpmInputModel.View()))
	// b.WriteString(fmt.Sprintf("\n\n%s\n\n", m.metronomeFrameStyle.Render(m.frames[m.frame])))
	b.WriteString(fmt.Sprintf("\n\n%s\n\n", m.metronomeDisplay))

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
