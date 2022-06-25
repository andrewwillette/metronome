package main

import (
	metrlog "github.com/andrewwillette/metronome/log"
)

func main() {
	metrlog.ConfigureLog("metronome.log", false)
}
