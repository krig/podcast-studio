// Podcast Studio, a podcast mixing/mastering app under construction.
package main

import (
	"log"

	"github.com/krig/Go-SDL2/sdl"
	"github.com/krig/Go-SDL2/ttf"
	"github.com/krig/Go-SDL2/gfx"
	"github.com/krig/go-sox"
)

type Resources struct {
	renderer *sdl.Renderer
	TitleFont *ttf.Font
	PlayButton *sdl.Texture
	StopButton *sdl.Texture
	LEDButton *sdl.Texture
	BackgroundColor sdl.Color
	TitleBarColor sdl.Color
	TitleColor sdl.Color
}

type Node struct {
	Widget
	label string
	color sdl.Color
	tex *sdl.Texture
}

type InputNode struct {
	Node
}

type OutputNode struct {
	Node
}

type EffectNode struct {
	Node
}

type Link struct {
	Widget
	A *Node
	B *Node
	Color sdl.Color
}

type SoundChain struct {
	in, out *sox.Format
	chain *sox.EffectsChain
	effects []*sox.Effect

	nodes []*Node
	links []*Link
}

type TopBar struct {
	HorizontalLayout
	BackgroundColor sdl.Color
}

type CanvasPane struct {
	Pane
	soundchain SoundChain
	menu PopupMenu
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

type InputStack struct {
	lovers []MouseLover
}

func (r *Resources) Load(rend *sdl.Renderer) {

	r.renderer = rend

	cfg := NewConfig("data/config.json")
	r.TitleFont = ttf.OpenFont(cfg.String("TitleFont"), cfg.Int("TitleFontSize"))
	r.PlayButton = cfg.Texture(rend, "PlayButton")
	r.StopButton = cfg.Texture(rend, "StopButton")
	r.LEDButton = cfg.Texture(rend, "LEDButton")
	r.BackgroundColor = cfg.Color("BackgroundColor")
	r.TitleBarColor = cfg.Color("TitleBarColor")
	r.TitleColor = cfg.Color("TitleColor")
}

func (r *Resources) Free() {
	r.TitleFont.Close()
	r.PlayButton.Destroy()
}

func (stack *InputStack) Add(lover MouseLover) {
	stack.lovers = append(stack.lovers, lover)
}

func (stack *InputStack) OnMouseMotionEvent(event *sdl.MouseMotionEvent) bool {
	for _, l := range stack.lovers {
		l.OnMouseMotionEvent(event)
	}
	return true
}

func (stack *InputStack) OnMouseButtonEvent(event *sdl.MouseButtonEvent) bool {
	for _, l := range stack.lovers {
		l.OnMouseButtonEvent(event)
	}
	return true
}

func (tb *TopBar) Draw(rend *sdl.Renderer) {
	rend.SetDrawColor(tb.BackgroundColor)
	rend.FillRect(&tb.Pos)
	rend.SetDrawColor(lighten(tb.BackgroundColor, 9))
	rend.DrawLine(tb.Pos.X, tb.Pos.Y + tb.Pos.H - 1, tb.Pos.X + tb.Pos.W, tb.Pos.Y + tb.Pos.H - 1)
	rend.SetDrawColor(darken(tb.BackgroundColor, 9))
	rend.DrawLine(tb.Pos.X, tb.Pos.Y + tb.Pos.H, tb.Pos.X + tb.Pos.W, tb.Pos.Y + tb.Pos.H)

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

func (canvas *CanvasPane) Init(rsc *Resources, space sdl.Rect) {
	canvas.Pos = space
	canvas.menu.Init(rsc.renderer, space, []string{"New Input", "New Output", "New Effect"}, rsc.TitleFont)

	canvas.menu.OnClick(func(entry *MenuEntry) {
		log.Println("Clicked: " + entry.Text)
		if entry.Text == "New Input" {
			canvas.NewInput()
		} else if entry.Text == "New Output" {
			canvas.NewOutput()
		} else if entry.Text == "New Effect" {
			canvas.NewEffect()
		}
	})
}

func (canvas *CanvasPane) NewInput() {
	
}

func (canvas *CanvasPane) NewOutput() {
}

func (canvas *CanvasPane) NewEffect() {
}

func (canvas *CanvasPane) Draw(rend *sdl.Renderer) {
	canvas.menu.Draw(rend)
}

func (canvas *CanvasPane) UpdateLayout(space sdl.Rect) {
	canvas.menu.UpdateLayout(space)
}

func (canvas *CanvasPane) OnMouseMotionEvent(event *sdl.MouseMotionEvent) bool {
	if canvas.menu.Visible {
		canvas.menu.OnMouseMotionEvent(event)
	}
	return true
}

func (canvas *CanvasPane) OnMouseButtonEvent(event *sdl.MouseButtonEvent) bool {
	if event.Button == sdl.BUTTON_RIGHT && event.State == sdl.PRESSED {
		if !canvas.menu.Visible {
			canvas.menu.Show(event.X, event.Y)
		} else {
			canvas.menu.Hide()
		}
	} else if canvas.menu.Visible {
		canvas.menu.OnMouseButtonEvent(event)
	}
	return true
}


func (screen *Screen) Init(space sdl.Rect, rsc *Resources) {
	topbarHeight := int32(32)
	screen.Pos = space
	screen.TopBar = &TopBar{}
	screen.TopBar.Init(sdl.Rect{space.X, space.Y, space.W, topbarHeight})
	screen.TopBar.BackgroundColor = rsc.TitleBarColor
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
	screen.Title.Init(rsc.renderer, sdl.Rect{screen.F2.Pos.X + screen.F2.Pos.W, space.Y, 1, topbarHeight}, "canvas mode", rsc.TitleFont, rsc.TitleColor)
	screen.Play.Init(rsc.renderer, sdl.Rect{screen.Title.Pos.X + screen.Title.Pos.W, space.Y, 1, 1}, rsc.PlayButton)
	screen.Stop.Init(rsc.renderer, sdl.Rect{screen.Play.Pos.X + screen.Play.Pos.W, space.Y, 1, 1}, rsc.StopButton)

	screen.F1.OnClick(func() {
		log.Println("Canvas Mode clicked!")
	})

	screen.F2.OnClick(func() {
		log.Println("Track Mode clicked!")
	})

	screen.Play.OnClick(func() {
		log.Println("Play clicked!")
	})

	screen.Stop.OnClick(func() {
		log.Println("Stop clicked!")
	})

	screen.AddVisual(screen.TopBar)
	screen.AddLayout(screen.TopBar)

	screen.Canvas = &CanvasPane{}
	screen.Canvas.Init(rsc, sdl.Rect{space.X, space.Y + topbarHeight, space.W, space.H - topbarHeight})
	screen.AddVisual(screen.Canvas)
	screen.AddLayout(screen.Canvas)

	screen.UpdateLayout(space)
}


func run_studio(window *sdl.Window, rend *sdl.Renderer) {
	rsc := &Resources{}
	rsc.Load(rend)
	defer rsc.Free()

	w, h := window.GetSize()
	neww, newh := w, h
	screen := &Screen{}
	screen.Init(sdl.Rect{0, 0, int32(w), int32(h)}, rsc)
	stack := &InputStack{}
	stack.Add(screen.F1)
	stack.Add(screen.F2)
	stack.Add(screen.Play)
	stack.Add(screen.Stop)
	stack.Add(screen.Canvas)
	defer screen.Destroy()

	framerate := gfx.NewFramerate()
	framerate.SetFramerate(20)
	event := &sdl.Event{}
	running := true
	for running {
		for event.Poll() {
			switch e := event.Get().(type) {
			case sdl.QuitEvent:
				running = false

			case sdl.KeyboardEvent:
				if e.Keysym.Keycode == sdl.K_ESCAPE {
					running = false
				}

			case sdl.MouseMotionEvent:
				stack.OnMouseMotionEvent(&e)

			case sdl.MouseButtonEvent:
				stack.OnMouseButtonEvent(&e)
			}
		}

		neww, newh = window.GetSize()
		if neww != w || newh != h {
			w, h = neww, newh
			screen.UpdateLayout(sdl.Rect{0, 0, int32(w), int32(h)})
		}

		rend.SetDrawBlendMode(sdl.BLENDMODE_NONE)
		rend.SetDrawColor(rsc.BackgroundColor)
		rend.Clear()
		rend.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
		screen.Draw(rend)
		rend.Present()
		framerate.FramerateDelay()
	}
}
