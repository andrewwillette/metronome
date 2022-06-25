package main

import (
	metrlog "github.com/andrewwillette/metronome/log"
	"github.com/andrewwillette/metronome/ui"
)

func main() {
	metrlog.ConfigureLog("metronome.log", false)
	ui.StartUi()
}
