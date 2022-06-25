package ui

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func Test_bpm2bps(t *testing.T) {
	tts := []struct {
		exptd time.Duration
		bpm   int
	}{
		{
			exptd: time.Millisecond * 60,
			bpm:   1000,
		},
		{
			exptd: time.Millisecond * 667,
			bpm:   90,
		},
		{
			exptd: time.Second,
			bpm:   60,
		},
	}
	for _, v := range tts {
		result := bpm2bps(v.bpm)
		assert.Equal(t, v.exptd, result.Round(time.Millisecond))
	}
}

func Test_getFrames(t *testing.T) {
	tts := []struct {
		provided, expected []string
	}{
		{
			provided: []string{"G", "B"},
			expected: []string{"G", " B"},
		},
	}
	for _, v := range tts {
		frames := getFrames(v.provided)
		assert.Equal(t, frames, v.expected)
	}
}

func Test_getBpmFromString(t *testing.T) {
	tts := []struct {
		provided  string
		expected  int
		expectErr bool
	}{
		{provided: "54", expected: 54},
		{provided: "a", expected: 0, expectErr: true},
	}
	for _, v := range tts {
		result, err := getBpmFromString(v.provided)
		if v.expectErr {
			assert.NotNil(t, err)
			continue
		}
		assert.Equal(t, v.expected, result)
	}
}

func Test_newModel(t *testing.T) {
	initMod := newModel()
	assert.Equal(t, lastID, initMod.id)
	newerModel := newModel()
	assert.Equal(t, lastID, newerModel.id)
}

func Test_initModel(t *testing.T) {
	m := newModel()
	res := m.Init()
	assert.IsType(t, TickMsg{}, res())
}

func Test_View(t *testing.T) {
	m := newModel()
	view := m.View()
	assert.Contains(t, view, getFrames(defaultMetronome.Frames)[0])
}

func Test_Update(t *testing.T) {
	t.Run("Updating with KeyMsg", func(t *testing.T) {
		m := newModel()
		tts := []struct {
			keyType tea.KeyType
			runes   []rune
		}{
			{
				keyType: tea.KeyCtrlC,
				runes:   nil,
			},
			{
				keyType: tea.KeyCtrlR,
				runes:   nil,
			},
			{
				runes: []rune{'1'},
			},
		}
		for _, v := range tts {
			key := tea.KeyMsg(tea.Key{
				Type:  v.keyType,
				Runes: v.runes,
			})
			m.cursorMode = textinput.CursorHide + 1
			_, _ = m.Update(key)
		}
	})

	t.Run("Updating with TickMsg", func(t *testing.T) {
		m := newModel()
		// much of this is no good
		tts := []struct {
			frame,
			tickId,
			tag,
			modelId,
			expectedTagVal int
			nilCmdReturned bool
		}{
			{
				frame:          1,
				tag:            1,
				modelId:        1,
				tickId:         1,
				expectedTagVal: 2,
			},
		}
		for _, v := range tts {
			m.frame = len(m.metronome.Frames)
			tm := TickMsg{
				ID:  v.tickId,
				tag: v.tag,
			}
			m.id = v.modelId
			m.tag = v.tag
			_, cmd := m.Update(tm)
			if !v.nilCmdReturned {
				assert.NotNil(t, cmd)
			}
		}
	})
}

func Test_ID(t *testing.T) {
	m := newModel()
	res := m.ID()
	assert.Equal(t, lastID, res)
}

func Test_GetSongsFromXdgConfig(t *testing.T) {
	getSongsFromXdgConfig()
}
