package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	// "fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	// "math"
	"time"
)

type CircularProgress struct {
    widget.BaseWidget
    progress float64
    running  bool
}

func NewCircularProgress() *CircularProgress {
    c := &CircularProgress{}
    c.ExtendBaseWidget(c)
    return c
}

func (c *CircularProgress) CreateRenderer() fyne.WidgetRenderer {
    circle := canvas.NewCircle(theme.PrimaryColor())
    circle.StrokeWidth = 5
    circle.StrokeColor = theme.PrimaryColor()
    return &circularProgressRenderer{circle: circle, progress: c}
}

func (c *CircularProgress) Start() {
    c.running = true
    go func() {
        for c.running {
            c.progress += 0.01
            if c.progress > 1 {
                c.progress = 0
            }
            c.Refresh()
            time.Sleep(50 * time.Millisecond)
        }
    }()
}

func (c *CircularProgress) Stop() {
    c.running = false
}

type circularProgressRenderer struct {
    circle   *canvas.Circle
    progress *CircularProgress
}

func (r *circularProgressRenderer) Layout(size fyne.Size) {
    r.circle.Resize(size)
    r.circle.Move(fyne.NewPos(0, 0))
}

func (r *circularProgressRenderer) MinSize() fyne.Size {
    return fyne.NewSize(50, 50)
}

func (r *circularProgressRenderer) Refresh() {
    _ = 360 * r.progress.progress
    r.circle.StrokeWidth = 5
    r.circle.StrokeColor = theme.PrimaryColor()
    r.circle.FillColor = theme.BackgroundColor()
    r.circle.Refresh()
    r.circle.StrokeWidth = 5
    r.circle.StrokeColor = theme.PrimaryColor()
    r.circle.FillColor = theme.BackgroundColor()
    r.circle.Refresh()
}

func (r *circularProgressRenderer) BackgroundColor() color.Color {
    return theme.BackgroundColor()
}

func (r *circularProgressRenderer) Objects() []fyne.CanvasObject {
    return []fyne.CanvasObject{r.circle}
}

func (r *circularProgressRenderer) Destroy() {}