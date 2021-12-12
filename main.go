package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/nsf/termbox-go"
)

// write log to file
func configureLog() {
	f, err := os.OpenFile("metronome.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

func metronome(BeatsPerMinute int) {
	milliseconds := 60000 / BeatsPerMinute
	for range time.Tick(time.Millisecond * time.Duration(milliseconds)) {
		fmt.Println("Tick")
	}
}

type Label struct {
	name string
	x, y int
	w, h int
	body string
}

func NewLabel(name string, x, y int, body string) *Label {
	lines := strings.Split(body, "\n")

	w := 0
	for _, l := range lines {
		if len(l) > w {
			w = len(l)
		}
	}
	h := len(lines) + 1
	w = w + 1

	return &Label{name: name, x: x, y: y, w: w, h: h, body: body}
}

func (l *Label) Layout(g *gocui.Gui) error {
	v, err := g.SetView(l.name, l.x, l.y, l.x+l.w, l.y+l.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprint(v, l.body)
	}
	return nil
}

type Input struct {
	name, value        string
	x, y, w, maxLength int
	metronome          *Metronome
}

func NewInput(name string, x, y, w, maxLength int, metr *Metronome) *Input {
	return &Input{name: name, x: x, y: y, w: w, maxLength: maxLength, metronome: metr}
}

func (i *Input) Layout(g *gocui.Gui) error {
	v, err := g.SetView(i.name, i.x, i.y, i.x+i.w, i.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editor = i
		v.Editable = true
	}
	return nil
}

// Edit the input, set the metronome bpm to the value
func (i *Input) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	cx, _ := v.Cursor()
	ox, _ := v.Origin()
	limit := ox+cx+1 > i.maxLength
	switch {
	case ch != 0 && mod == 0 && !limit:
		v.EditWrite(ch)
	case key == gocui.KeySpace && !limit:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	}
	//metronomeView := g
	i.value = v.Buffer()
	beatsPerMinutes, err := strconv.Atoi(strings.Replace(v.Buffer(), "\n", "", -1))
	if err != nil {
		log.Println(err)
	}
	log.Printf("setting metr bpm to %d", beatsPerMinutes)
	i.metronome.bpm = beatsPerMinutes
}

type Metronome struct {
	name string
	bpm  int
	x, y int
	w, h int
}

func setMetronomeBpm(metronome *Metronome) error {
	metronome.bpm = 5
	return nil
}

const tickerSymbol = "x"

func tickMetronome(beatsPerMinute int, view *gocui.View, g *gocui.Gui) {
	if beatsPerMinute > 0 {
		milliseconds := 60000 / beatsPerMinute
		for range time.Tick(time.Millisecond * time.Duration(milliseconds)) {
			//view.Clear()
			if view.Buffer() == tickerSymbol {
				log.Println("here5")
				fmt.Fprintf(view, "")
			} else {
				log.Println("here6")
				fmt.Fprintf(view, tickerSymbol)
			}
			termbox.Sync()
			log.Println("tick")
			log.Println(view.Buffer())
		}
	}
}

func (metronome *Metronome) Layout(g *gocui.Gui) error {
	log.Println("metronome.Layout()")
	v, err := g.SetView(metronome.name, metronome.x, metronome.y, metronome.x+metronome.w, metronome.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	//v.Clear()
	//fmt.Fprintf(v, fmt.Sprint(metronome.bpm))
	v.Editable = true
	go tickMetronome(metronome.bpm, v, g)
	return nil
}

func NewMetronome(name string, x, y, w int) *Metronome {
	return &Metronome{name: name, x: x, y: y, w: w}
}

// Set focus for gui
func SetFocus(name string) func(g *gocui.Gui) error {
	return func(g *gocui.Gui) error {
		_, err := g.SetCurrentView(name)
		return err
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func startApp() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true
	label := NewLabel("Beats Per Minute Label", 1, 1, "Beats Per Minute")
	metronome := NewMetronome("metronome", 10, 10, 40)
	input := NewInput("beatsPerMinuteInput", 7, 1, 40, 40, metronome)
	focus := gocui.ManagerFunc(SetFocus("beatsPerMinuteInput"))
	g.SetManager(label, input, metronome, focus)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func main() {
	metronome := &Metronome{bpm: 60}
	configureLog()
	log.Println(metronome)
	startApp()
}
