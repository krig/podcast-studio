package main

import (
	"github.com/krig/Go-SDL2/sdl"
	"github.com/krig/Go-SDL2/ttf"
	"log"
)

type Visual interface {
	Draw(rend* sdl.Renderer)
	Destroy()
	GetPos() sdl.Rect
	SetPos(pos sdl.Rect)
}

type Layout interface {
	UpdateLayout(space sdl.Rect)
}

type Widget struct {
	Pos sdl.Rect
}

type Pane struct {
	Widget
	elements []Visual
	layouts []Layout
}

type Button struct {
	Widget
	// State: 0 = off, 1 = hover, 2 = on
	State int
	texture *sdl.Texture
}

type Label struct {
	Widget
	Text string
	Font *ttf.Font
	texture *sdl.Texture
}

type HorizontalLayout struct {
	Widget
	elements_left []Visual
	element_center Visual
	elements_right []Visual
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

type CanvasPane struct {
	Pane
	nodes []*Node
	links []*Link
}

func (widget *Widget) GetPos() sdl.Rect {
	return widget.Pos
}

func (widget *Widget) SetPos(pos sdl.Rect) {
	widget.Pos = pos
}

func (pane *Pane) Draw(rend *sdl.Renderer) {
	for _, w := range pane.elements {
		w.Draw(rend)
	}
}

func (pane *Pane) AddVisual(visual Visual) {
	pane.elements = append(pane.elements, visual)
}

func (pane *Pane) AddLayout(layout Layout) {
	pane.layouts = append(pane.layouts, layout)
}

func (pane *Pane) UpdateLayout(space sdl.Rect) {
	pane.Pos = space
	for _, l := range pane.layouts {
		l.UpdateLayout(space)
	}
}

func (pane *Pane) Destroy() {
	for _, e := range pane.elements {
		e.Destroy()
	}
}

func (h *HorizontalLayout) UpdateLayout(space sdl.Rect) {
	x := space.X
	h.Pos = space
	for _, w := range h.elements_left {
		pos := w.GetPos()
		w.SetPos(sdl.Rect{x, space.Y, pos.W, pos.H})
		x += pos.W
	}
	if h.element_center != nil {
		pos := h.element_center.GetPos()
		h.element_center.SetPos(sdl.Rect{x, space.Y, pos.W, pos.H})
	}

	x = space.X + space.W
	for i := len(h.elements_right) - 1; i >= 0; i-- {
		w := h.elements_right[i]
		pos := w.GetPos()
		w.SetPos(sdl.Rect{x - pos.W, space.Y, pos.W, pos.H})
		x -= pos.W
	}
}

func (h *HorizontalLayout) AddLeft(widget Visual) {
	h.elements_left = append(h.elements_left, widget)
}

func (h *HorizontalLayout) AddRight(widget Visual) {
	h.elements_right = append(h.elements_right, widget)
}

func (h *HorizontalLayout) SetCenter(widget Visual) {
	h.element_center = widget
}

func (h *HorizontalLayout) Destroy() {
	for _, w := range h.elements_left {
		w.Destroy()
	}
	if h.element_center != nil {
		h.element_center.Destroy()
	}
	for _, w := range h.elements_right {
		w.Destroy()
	}
}

func (tb *TopBar) Draw(rend *sdl.Renderer) {
	for _, w := range tb.elements_left {
		w.Draw(rend)
	}
	if tb.element_center != nil {
		tb.element_center.Draw(rend)
	}
	for _, w := range tb.elements_right {
		w.Draw(rend)
	}
}

func (canvas *CanvasPane) Draw(rend *sdl.Renderer) {
}

func (canvas *CanvasPane) UpdateLayout(space sdl.Rect) {
}

func (button *Button) Init(rend *sdl.Renderer, space sdl.Rect, texture *sdl.Texture) {
	button.Pos = space
	button.texture = texture
	w, h := texture.GetSize()
	button.Pos.W = int32(w) / 3
	button.Pos.H = int32(h)
	button.State = 0
}

func (label *Label) Init(rend *sdl.Renderer, space sdl.Rect, text string, font *ttf.Font) {
	label.Pos = space
	label.Text = text
	label.Font = font
	label.texture = nil
	label.Update(rend)
}

func (button *Button) Destroy() {
	if button.texture != nil {
		button.texture.Destroy()
		button.texture = nil
	}
}


func (label *Label) Destroy() {
	if label.texture != nil {
		label.texture.Destroy()
		label.texture = nil
	}
}

func hexcolor(code uint32) sdl.Color {
	return sdl.Color{uint8((code >> 16) & 0xff), uint8((code >> 8) & 0xff), uint8(code & 0xff), 0xff}
}

func (label *Label) Update(rend *sdl.Renderer) {
	textw, texth, err := label.Font.SizeText(label.Text)
	if err != nil {
		log.Fatal(err)
	}
	txt_surface := label.Font.RenderText_Blended(label.Text, hexcolor(0xffffff))
	if label.texture != nil {
		label.texture.Destroy()
	}
	label.texture = rend.CreateTextureFromSurface(txt_surface)
	txt_surface.Free()
	label.Pos.W = int32(textw)
	label.Pos.H = int32(texth)
}

func (button *Button) Draw(rend *sdl.Renderer) {
	w := button.Pos.W
	rend.Copy(button.texture, &sdl.Rect{int32(button.State) * w, 0, button.Pos.W, button.Pos.H}, &button.Pos)
}

func (label *Label) Draw(rend *sdl.Renderer) {
	rend.Copy(label.texture, nil, &label.Pos)
}

type Screen struct {
	Pane
	TopBar *TopBar
	F1 *Button
	F2 *Button
	Title *Label
	Play *Button
	Stop *Button

	Canvas *CanvasPane
	//Tracks *TrackPane
	//Current *Pane
}

func (screen *Screen) Init(space sdl.Rect, rsc *Resources) {
	topbarHeight := int32(32)
	screen.Pos = space
	screen.TopBar = &TopBar{}
	screen.TopBar.Pos = sdl.Rect{space.X, space.Y, space.W, topbarHeight}
	screen.F1 = &Button{}
	screen.F2 = &Button{}
	screen.Title = &Label{}
	screen.Play = &Button{}
	screen.Stop = &Button{}

	screen.TopBar.AddLeft(screen.F1)
	screen.TopBar.AddLeft(screen.F2)
	screen.TopBar.SetCenter(screen.Title)
	screen.TopBar.AddRight(screen.Play)
	screen.TopBar.AddRight(screen.Stop)

	screen.F1.Init(rsc.renderer, sdl.Rect{space.X, space.Y, 1, 1}, rsc.LEDButton)
	screen.F2.Init(rsc.renderer, sdl.Rect{screen.F1.Pos.X + screen.F1.Pos.W, space.Y, 1, 1}, rsc.LEDButton)
	screen.Title.Init(rsc.renderer, sdl.Rect{screen.F2.Pos.X + screen.F2.Pos.W, space.Y, 1, 1}, "canvas mode", rsc.TitleFont)
	screen.Play.Init(rsc.renderer, sdl.Rect{screen.Title.Pos.X + screen.Title.Pos.W, space.Y, 1, 1}, rsc.PlayButton)
	screen.Stop.Init(rsc.renderer, sdl.Rect{screen.Play.Pos.X + screen.Play.Pos.W, space.Y, 1, 1}, rsc.LEDButton)

	screen.AddVisual(screen.TopBar)
	screen.AddLayout(screen.TopBar)

	screen.Canvas = &CanvasPane{}
	screen.Canvas.Pos = sdl.Rect{space.X, space.Y + topbarHeight, space.W, space.H - topbarHeight}
	screen.AddVisual(screen.Canvas)
	screen.AddLayout(screen.Canvas)

	screen.UpdateLayout(space)
}

