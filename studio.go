package main

import (
	"github.com/krig/Go-SDL2/sdl"
	"github.com/krig/Go-SDL2/ttf"
	"log"
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

type PopupMenu struct {
	Widget
	entries []*sdl.Texture
	Visible bool
}

type Node struct {
	Widget
	Type int
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

func (menu *PopupMenu) Draw(rend *sdl.Renderer) {
}

func (canvas *CanvasPane) Draw(rend *sdl.Renderer) {
}

func (canvas *CanvasPane) UpdateLayout(space sdl.Rect) {
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
	screen.Canvas.Pos = sdl.Rect{space.X, space.Y + topbarHeight, space.W, space.H - topbarHeight}
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
	defer screen.Destroy()

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
			w = neww
			h = newh
			screen.UpdateLayout(sdl.Rect{0, 0, int32(w), int32(h)})
		}

		rend.SetDrawBlendMode(sdl.BLENDMODE_NONE)
		rend.SetDrawColor(rsc.BackgroundColor)
		rend.Clear()
		rend.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
		screen.Draw(rend)
		rend.Present()
	}
}
