package main

import (
	"testing"
)

func Test_getHeader(t *testing.T) {
	if GetApplicationTitle() != "Metronome" {
		t.Error("header incorrect")
	}
}

// func Test_getWindow(t *testing.T) {
// // window.FullScreen
// assert true
// }
