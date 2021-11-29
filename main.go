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
	x, y, w, maxLength int
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
    //log.Printf("editing swag\n, key: %+v", v)
    //dumper := hex.Dumper(v)
    //if _, err := io.Copy(dumper, os.Stdin); err != nil {
        //log.Println("failed to copy dumper lmao")
    //}
    //log.Printf(dumper.Close())
    //log.Printf(dumper)
    //io.Copy(dumper, log.Reader)

    log.Printf("buffer lines %s", v.BufferLines())
    log.Printf("ch: %+v", ch)
    log.Printf("mod: %+v", mod)
    log.Printf("key: %+v", key)
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

// for early implementation using Gocui.SetManagerFunc(layout)
func layout(g *gocui.Gui) error {
    if v, err := g.SetView("stdin", 0, 0, 80, 35); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("stdin"); err != nil {
			return err
        }
        log.Println(v, "this garbage works lmao")
        log.Println(v, "alright alright")
		v.Wrap = true
	}
	return nil
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
    label := NewLabel("test label", 1, 1, "testLabel")
    input := NewInput("testInput", 7, 1, 40, 40)
    focus := gocui.ManagerFunc(SetFocus("testInput"))
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
    startApp()
}
