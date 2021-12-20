package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

var (
	done               = make(chan struct{})
	globalGui          *gocui.Gui
)

func main() {
	configureLog()
	startMetronome()
}

// write log to file
func configureLog() {
	f, err := os.OpenFile("metronome.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	log.SetOutput(f)
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

type MetronomeInput struct {
	name        string
	x, y, w, maxLength int
}

func NewMetronomeInput(name string, x, y, w, maxLength int) *MetronomeInput {
	return &MetronomeInput{name: name, x: x, y: y, w: w, maxLength: maxLength}
}

func (i *MetronomeInput) Layout(g *gocui.Gui) error {
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

// Edit metronome input
// Trigger goroutine updating the metronome view at bpm
func (i *MetronomeInput) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
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
	beatsPerMinutes, err := strconv.Atoi(strings.Replace(v.Buffer(), "\n", "", -1))
	if err != nil {
		log.Println(err)
	}
    // causes old metronomeCounter goroutine to return
	close(done)
	done = make(chan struct{})
	if beatsPerMinutes > 0 {
        go metronomeCounter(beatsPerMinutes)
	}
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

func (metronome *Metronome) Layout(g *gocui.Gui) error {
	v, err := g.SetView(metronome.name, metronome.x, metronome.y, metronome.x+metronome.w, metronome.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.Editable = true
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

// quit gui app, close display ticker goroutine
func quit(g *gocui.Gui, v *gocui.View) error {
	close(done)
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func startMetronome() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	globalGui = g
	if err != nil {
		log.Panicln(err)
	}
	label := NewLabel("Beats Per Minute Label", 1, 1, "Beats Per Minute")
	bpmInput := NewMetronomeInput("metronomeInput", 1, 3, 40, 40)
	metronome := NewMetronome("metronome", 1, 6, 40)
	focus := gocui.ManagerFunc(SetFocus("metronomeInput"))
	globalGui.SetManager(label, bpmInput, metronome, focus)
	if err := keybindings(globalGui); err != nil {
		log.Panicln(err)
	}
	if err := globalGui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

var tickFlag bool = true

func metronomeCounter(beatsPerMinute int) {
	milliseconds := 60000 / beatsPerMinute
	for {
		select {
		case <-done:
			return
		case <-time.After(time.Duration(milliseconds) * time.Millisecond):
			globalGui.Update(func(g *gocui.Gui) error {
				v, err := g.View("metronome")
				if err != nil {
					log.Panicln("error getting metronome view")
					return err
				}
				v.Clear()
				if tickFlag {
					fmt.Fprintln(v, "❌")
					tickFlag = false
				} else {
					fmt.Fprintln(v, "   ❎")
					tickFlag = true
				}
				return nil
			})
		}
	}
}
