package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

// write log to file
func configureLog() {
    f, err := os.OpenFile("metronome.log", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
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
	name      string
	x, y      int
	w         int
	maxLength int
}

func NewInput(name string, x, y, w, maxLength int) *Input {
	return &Input{name: name, x: x, y: y, w: w, maxLength: maxLength}
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
}

func layout(g *gocui.Gui) error {
	//maxX, maxY := g.Size()
	//if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		//if err != gocui.ErrUnknownView {
			//return err
		//}

		//fmt.Fprintln(v, "Hello world!")
	//}
    if v, err := g.SetView("stdin", 0, 0, 80, 35); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("stdin"); err != nil {
			return err
        }
        fmt.Fprintln(v, "this garbage works lmao")
        fmt.Fprintln(v, "alright alright")
		v.Wrap = true
	}
	return nil
}

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
    //var beatsPerMinute int
    //_, err := fmt.Scan(&beatsPerMinute)
    //if err != nil {
        //log.Fatal("failed to scan bpm", err)
    //}
    g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
    g.Cursor = true
    label := NewLabel("bitchlabel\nand tits", 1, 1, "bitchShit")
    input := NewInput("bitchInput", 7, 1, 40, 40)
    focus := gocui.ManagerFunc(SetFocus("bitchInput"))
    g.SetManager(label, input, focus)
	//g.SetManagerFunc(layout)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func main() {
    configureLog()
    log.Println("gets this")
    startApp()
}
