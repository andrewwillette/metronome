package main

import (
	"testing"
	"time"

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
