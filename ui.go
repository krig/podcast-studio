package main

import (
	"log"

	"github.com/krig/Go-SDL2/sdl"
	"github.com/krig/Go-SDL2/ttf"
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

type Clickable interface {
	OnClick(handler func())
}

type MouseLover interface {
	GetPos() sdl.Rect
	OnMouseButtonEvent(event *sdl.MouseButtonEvent) bool
	OnMouseMotionEvent(event *sdl.MouseMotionEvent) bool
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
	// State: 0 = off, 1 = hover, 2 = on, 3 = on+hover
	State int
	texture *sdl.Texture
	clickhandler func()
}

type Label struct {
	Widget
	Text string
	Font *ttf.Font
	Color sdl.Color
	texture *sdl.Texture
	texwidth int32
	texheight int32
}

type HorizontalLayout struct {
	Widget
	elements_left []Visual
	element_center Visual
	elements_right []Visual
	HSpacing int32
}

type MenuEntry struct {
	Text string
	label *Label
}

type PopupMenu struct {
	Widget
	entries []MenuEntry
	hover int
	Visible bool
	Spacing, Inset int32
	clickhandler func(entry *MenuEntry)
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

func (h *HorizontalLayout) Init(space sdl.Rect) {
	h.Pos = space
	h.HSpacing = 2
}

func (h *HorizontalLayout) UpdateLayout(space sdl.Rect) {
	x := space.X + h.HSpacing
	h.Pos = sdl.Rect{space.X, space.Y, space.W, h.Pos.H}
	for _, w := range h.elements_left {
		pos := w.GetPos()
		w.SetPos(sdl.Rect{x, space.Y, pos.W, pos.H})
		x += pos.W + h.HSpacing
	}

	rx := space.X + space.W - h.HSpacing
	for i := len(h.elements_right) - 1; i >= 0; i-- {
		w := h.elements_right[i]
		pos := w.GetPos()
		w.SetPos(sdl.Rect{rx - pos.W, space.Y, pos.W, pos.H})
		rx -= (pos.W + h.HSpacing)
	}

	if h.element_center != nil {
		h.element_center.SetPos(sdl.Rect{x, space.Y, rx - x, h.Pos.H})
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

func (button *Button) Init(rend *sdl.Renderer, space sdl.Rect, texture *sdl.Texture) {
	button.Pos = space
	button.texture = texture
	w, h := texture.GetSize()
	button.Pos.W = int32(w) / 3
	button.Pos.H = int32(h)
	button.State = 0
}

func (label *Label) Init(rend *sdl.Renderer, space sdl.Rect, text string, font *ttf.Font, color sdl.Color) {
	label.Pos = space
	label.Text = text
	label.Font = font
	label.texture = nil
	label.Color = color
	label.Update(rend)
}

func (button *Button) Destroy() {
	if button.texture != nil {
		button.texture.Destroy()
		button.texture = nil
	}
}

func (button *Button) OnClick(handler func()) {
	button.clickhandler = handler
}

func (button *Button) OnMouseMotionEvent(event *sdl.MouseMotionEvent) bool {
	if button.Pos.Contains(event.X, event.Y) {
		button.State |= 1
	} else {
		button.State &= ^1
	}
	return true
}

func (button *Button) OnMouseButtonEvent(event *sdl.MouseButtonEvent) bool {
	if event.Button == sdl.BUTTON_LEFT {
		contains := button.Pos.Contains(event.X, event.Y)
		pressed := event.State == sdl.PRESSED
		was_pressed := (button.State & 2) != 0
		if pressed && contains {
			button.State |= 2
		} else if was_pressed && !pressed && contains {
			button.State &= ^2
			if button.clickhandler != nil {
				button.clickhandler()
			}
		} else {
			button.State &= ^2
		}
	}
	return true
}


func (label *Label) Destroy() {
	if label.texture != nil {
		label.texture.Destroy()
		label.texture = nil
	}
}

func (label *Label) Update(rend *sdl.Renderer) {
	textw, texth, err := label.Font.SizeText(label.Text)
	if err != nil {
		log.Fatal(err)
	}
	txt_surface := label.Font.RenderText_Blended(label.Text, label.Color)
	if label.texture != nil {
		label.texture.Destroy()
	}
	label.texture = rend.CreateTextureFromSurface(txt_surface)
	txt_surface.Free()
	label.texwidth = int32(textw)
	label.texheight = int32(texth)
}

func (button *Button) Draw(rend *sdl.Renderer) {
	w := button.Pos.W
	state := button.State
	if state > 2 {
		state = 2
	}
	rend.Copy(button.texture, &sdl.Rect{int32(state) * w, 0, button.Pos.W, button.Pos.H}, &button.Pos)
}

func (label *Label) Draw(rend *sdl.Renderer) {
	pos := sdl.Rect{label.Pos.X + (label.Pos.W - label.texwidth) / 2, label.Pos.Y + (label.Pos.H - label.texheight)/2, label.texwidth, label.texheight}
	rend.Copy(label.texture, nil, &pos)
}

// Init a popup menu
func (menu *PopupMenu) Init(rend *sdl.Renderer, space sdl.Rect, entries []string, font *ttf.Font) {
	menu.Pos = space
	menu.Inset = 4
	menu.Spacing = 16
	entry_pos := sdl.Rect{menu.Pos.X + menu.Inset,
		menu.Pos.Y + menu.Inset,
		menu.Pos.W - menu.Inset*2,
		menu.Spacing}
	max_w := int32(0)
	for i, e := range entries {
		menu.entries = append(menu.entries, MenuEntry{})
		menu.entries[i].Init(rend, entry_pos, e, font, hexcolor(0xeeeeec))
		if menu.entries[i].label.texwidth > max_w {
			max_w = menu.entries[i].label.texwidth
		}
		entry_pos.Y += menu.Spacing
	}
	menu.Pos.W = max_w + menu.Inset*2
	menu.Pos.H = int32(len(entries))*menu.Spacing + menu.Inset*2

	for _, e := range menu.entries {
		e.label.Pos.W = menu.Pos.W - menu.Inset*2
		e.label.Pos.H = menu.Spacing
	}
}

func (menu *PopupMenu) OnClick(handler func(entry *MenuEntry)) {
	menu.clickhandler = handler
}

// Draw a popup menu
func (menu *PopupMenu) Draw(rend *sdl.Renderer) {
	if !menu.Visible {
		return
	}
	rend.SetDrawColor(hexcolor(0x303030))
	rend.FillRect(&menu.Pos)
	rend.SetDrawColor(hexcolor(0x363636))
	rend.DrawRect(&menu.Pos)
	if menu.hover != -1 {
		rend.SetDrawColor(hexcolor(0x406f40))
		rend.FillRect(&menu.entries[menu.hover].label.Pos)
	}
	for _, e := range menu.entries {
		e.Draw(rend)
	}
}

func (menu *PopupMenu) UpdateLayout(space sdl.Rect) {
	// TODO
}

func (menu *PopupMenu) Show(x, y int32) {
	menu.Visible = true
	menu.Pos.X = x
	menu.Pos.Y = y
	new_x := menu.Pos.X + menu.Inset
	new_y := menu.Pos.Y + menu.Inset
	for _, e := range menu.entries {
		e.label.Pos.X = new_x
		e.label.Pos.Y = new_y
		new_y += menu.Spacing
	}
}

func (menu *PopupMenu) Hide() {
	menu.Visible = false
	menu.hover = -1
}

func (entry *MenuEntry) Init(rend *sdl.Renderer, space sdl.Rect, text string, font *ttf.Font, color sdl.Color) {
	entry.Text = text
	entry.label = &Label{}
	entry.label.Init(rend, space, text, font, color)
}

func (entry *MenuEntry) Draw(rend *sdl.Renderer) {
	entry.label.Draw(rend)
}

func (menu *PopupMenu) OnMouseMotionEvent(event *sdl.MouseMotionEvent) bool {
	menu.hover = -1
	if menu.Visible {
		for i, e := range menu.entries {
			if e.label.Pos.Contains(event.X, event.Y) {
				//log.Println(e.Text, e.label.Pos, "contains", event.X, event.Y)
				menu.hover = i
			}
		}
	} else {
	}
	return true
}

func (menu *PopupMenu) OnMouseButtonEvent(event *sdl.MouseButtonEvent) bool {
	if menu.Visible {
		if event.Button == sdl.BUTTON_LEFT {
			contains := menu.Pos.Contains(event.X, event.Y)
			pressed := event.State == sdl.PRESSED
			was_pressed := menu.hover != -1
			if was_pressed && !pressed && contains {
				if menu.clickhandler != nil {
					menu.clickhandler(&menu.entries[menu.hover])
				}
			}
			if was_pressed && !pressed {
				menu.Hide()
			}
		}
	}
	return true
}
