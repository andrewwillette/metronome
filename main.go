package main

import (
	"time"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

type metronomeLayout struct {
    beatBlinker *canvas.Circle
    canvas fyne.CanvasObject
    stop bool
}

func (metronomeLayout *metronomeLayout) Layout(_ []fyne.CanvasObject, size fyne.Size) {
    diameter := fyne.Min(size.Width, size.Height)
    radius := diameter / 2
    size = fyne.NewSize(diameter, diameter)
    middle := fyne.NewPos(size.Width/2, size.Height/2)
    topleft := fyne.NewPos(middle.X-radius, middle.Y-radius)
    metronomeLayout.beatBlinker.Resize(size)
    metronomeLayout.beatBlinker.Move(topleft)
}

func (metronomeLayout *metronomeLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
    return fyne.NewSize(200, 200)
}

func (metronomeLayout *metronomeLayout) animate(canvasObject fyne.CanvasObject) {
	tick := time.NewTicker(time.Millisecond)
	go func() {
		for !metronomeLayout.stop {
			metronomeLayout.Layout(nil, canvasObject.Size())
			canvas.Refresh(metronomeLayout.canvas)
			<-tick.C
		}
	}()
}

func (metronome *metronomeLayout) render() *fyne.Container {
    metronome.beatBlinker = &canvas.Circle{StrokeColor: theme.TextColor(), StrokeWidth: 5}
    container := fyne.NewContainer(metronome.beatBlinker)
    container.Layout = metronome
    metronome.canvas = container
    return container
}

func (metronomeLayout *metronomeLayout) applyTheme (_ fyne.Settings) {
    metronomeLayout.beatBlinker.StrokeColor = theme.PrimaryColor()
}

func ShowMetronome(win fyne.Window) fyne.CanvasObject {
    metronome := &metronomeLayout{}
	content := metronome.render()
    go metronome.animate(content)
    listener := make(chan fyne.Settings)
    fyne.CurrentApp().Settings().AddChangeListener(listener)
    go func() {
        for {
            settings := <-listener
            metronome.applyTheme(settings)
        }
    }()
    return content
}

func GetApplicationTitle() string {
    return "Metronome"
}

func GetSize() fyne.Size {
    return fyne.NewSize(float32(60), float32(60))
}

func main() {
	a := app.New()
	w := a.NewWindow(GetApplicationTitle())
    w.SetContent(ShowMetronome(w))
    w.Resize(fyne.NewSize(480, 360))
    w.ShowAndRun()
}
