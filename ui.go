package main

import (
	"github.com/krig/Go-SDL2/sdl"
	"github.com/krig/Go-SDL2/ttf"
)

type Visual interface {
	Draw(rend* sdl.Renderer)
}

type Layout interface {
	UpdateLayout(space sdl.Rect)
}

type Widget struct {
	Pos sdl.Rect
}

type Pane struct {
	Widget
	Elements []Visual
	Layouts []Layout
}

type Button struct {
	Widget
	// State: 0 = off, 1 = hover, 2 = on
	State int
}

type Label struct {
	Widget
	Text string
	Font *ttf.Font
	texture *sdl.Texture
	dirty bool
}

type HorizontalLayout struct {
	Widget
	elements_left []*Widget
	element_center *Widget
	elements_right []*Widget
}

type TopBar struct {
	HorizontalLayout
}

type Node struct {
	Widget
}

type Link struct {
	Widget
	A *Node
	B *Node
}

struct CanvasPane struct {
	Pane
	nodes []*Node
	links []*Link
}

func (pane *Pane) Draw(rend *sdl.Renderer) {
	for w := range pane.elements {
		w.Draw(rend)
	}
}

func (pane *Pane) AddVisual(visual Visual) {
	pane.elements = append(pane.elements, visual)
}

func (pane *Pane) AddLayout(layout Layout) {
	pane.layouts = append(pane.layouts, visual)
}

func (pane *Pane) UpdateLayout(space sdl.Rect) {
	for l := range pane.layouts {
		layout.UpdateLayout()
	}
}

func (h *HorizontalLayout) AddLeft(widget *Widget) {
	h.elements_left = append(h.elements_left, widget)
}

func (h *HorizontalLayout) AddRight(widget *Widget) {
	h.elements_right = append(h.elements_right, widget)
}

func (h *HorizontalLayout) SetCenter(widget *Widget) {
	h.element_center = widget
}

func (tb *TopBar) Draw(rend *sdl.Renderer) {
	for w := range tb.elements_left {
		w.Draw(rend)
	}
	if tb.element_center != nil {
		tb.element_center.Draw(rend)
	}
	for w := range tb.elements_right {
		w.Draw(rend)
	}
}

func (canvas *CanvasPane) Draw(rend *sdl.Renderer) {
}

func (canvas *CanvasPane) UpdateLayout(space sdl.Rect) {
}

func (button *Button) Init() {
}

func (label *Label) Init() {
}

func (button *Button) Draw(rend *sdl.Renderer) {
	//switch button.State {
	//case 0:
	//	rend.Copy(button.
	//}
}

func (label *Label) Draw(rend *sdl.Renderer) {
	//rend.Copy(label.texture, nil, label.Pos)
}

type Screen struct {
	Pane
	TopBar *TopBar
	F1 *LEDButton
	F2 *LEDButton
	Title *Label
	Play *Button
	Stop *Button

	Canvas *CanvasPane
	//Tracks *TrackPane
	//Current *Pane
}

func (screen *Screen) Init(space sdl.Rect) {
	topbarHeight := 32
	screen.Pos = space
	screen.TopBar := &TopBar{}
	screen.TopBar.Pos = sdl.Rect{space.X, space.Y, space.W, topbarHeight}
	screen.F1 := &Button{}
	screen.F2 := &Button{}
	screen.Title := &Label{}
	screen.Play := &Button{}
	screen.Stop := &Button{}

	screen.TopBar.AddLeft(screen.F1)
	screen.TopBar.AddLeft(screen.F2)
	screen.TopBar.SetCenter(screen.Title)
	screen.TopBar.AddRight(screen.Play)
	screen.TopBar.AddRight(screen.Stop)

	screen.AddVisual(screen.TopBar)
	screen.AddLayout(screen.TopBar)

	screen.Canvas = &CanvasPane{}
	screen.Canvas.Pos = sdl.Rect{space.X, space.Y + topbarHeight, space.W, space.H - topbarHeight}
	screen.AddVisual(screen.Canvas)
	screen.AddLayout(screen.Canvas)

	pane.UpdateLayout(space)
}

