// Podcast Studio, a podcast mixing/mastering app under construction.
package main

import (
	"log"
	"math"
	"runtime"

	"github.com/krig/Go-SDL2/sdl"
	"github.com/krig/Go-SDL2/ttf"
	"github.com/krig/Go-SDL2/gfx"
	"github.com/krig/go-sox"
)

const (
	TOPBAR_HEIGHT = int32(32)
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

type FloatPos struct {
	X float64
	Y float64
}

type Node struct {
	Widget
	label Label
	color sdl.Color
	dragging bool
	curr FloatPos
	goal FloatPos
	menu PopupMenu

	name string
	args []string

	next *Node
}

// Color scheme:
// teal: 15f0e1
// red: ff3015
// green: 5be33b
// yellow: ffe018
// blue: 694ae9


func (node *Node) DrawLink(rend *sdl.Renderer) {
	if node.next != nil {
		rend.SetDrawColor(hexcolor(0x15f0e1))
		p0 := node.Pos
		p1 := node.next.Pos
		rend.DrawLine(p0.X + p0.W/2, p0.Y + p0.H/2, p1.X + p1.W/2, p1.Y + p1.H/2)
	}
}

func (node *Node) Draw(rend *sdl.Renderer) {
	if node.dragging {
		node.curr.X += (node.goal.X - node.curr.X) * (15.0 / 30.0)
		node.curr.Y += (node.goal.Y - node.curr.Y) * (15.0 / 30.0)
		node.Pos.X = int32(node.curr.X)
		node.Pos.Y = int32(node.curr.Y)
		node.label.Pos = node.Pos
	}

	clr := node.color
	if !node.dragging {
		rend.SetDrawColor(clr)
		rend.FillRect(&node.Pos)
	}
	rend.SetDrawColor(lighten(clr, 19))
	rend.DrawRect(&node.Pos)
	node.label.Draw(rend)

	if node.menu.Visible {
		node.menu.Draw(rend)
	}
}

func (node *Node) OnMouseMotionEvent(event *sdl.MouseMotionEvent) bool {
	if node.dragging {
		node.goal.X += float64(event.XRel)
		node.goal.Y += float64(event.YRel)

		node.goal.X = math.Max(node.goal.X, 0.0)
		node.goal.Y = math.Max(node.goal.Y, float64(TOPBAR_HEIGHT))
	}
	if node.menu.Visible {
		node.menu.OnMouseMotionEvent(event)
	}
	return true
}

func (node *Node) OnMouseButtonEvent(event *sdl.MouseButtonEvent) bool {
	lpress := event.State == sdl.PRESSED && event.Button == sdl.BUTTON_LEFT
	rpress := event.State == sdl.PRESSED && event.Button == sdl.BUTTON_RIGHT

	if node.menu.Visible {
		node.menu.OnMouseButtonEvent(event)
	}
	if node.Pos.Contains(event.X, event.Y) {
		if lpress {
			node.dragging = true
			node.goal.X = float64(node.Pos.X)
			node.goal.Y = float64(node.Pos.Y)
			node.curr = node.goal
		}
		if rpress && !node.menu.Visible {
			node.menu.Show(event.X, event.Y)
		} else if rpress {
			node.menu.Hide()
		}
	}
	if node.dragging && event.State == sdl.RELEASED {
		node.dragging = false
	}
	return true
}



// TODO
type SoxChain struct {
	chain *sox.EffectsChain
	in, out *sox.Format
	interrupt bool
	finished bool
}

func (chain *SoxChain) Release() {
	chain.chain.Release()
	chain.in.Release()
	chain.out.Release()
}

func (chain *SoxChain) Flow() {
	chain.chain.FlowCallback(func(all_done bool) int {
		if chain.interrupt {
			return 1
		}
		return 0
	})
	chain.finished = true
}

type TopBar struct {
	HorizontalLayout
	BackgroundColor sdl.Color
}

type CanvasPane struct {
	Pane
	menu PopupMenu

	rsc *Resources

	tracks []string

	nodes []*Node
	new_link *int

	playing *SoxChain
}

type ListWindow struct {
	Widget
	title_text Label
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
		if !l.OnMouseButtonEvent(event) {
			return false
		}
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

func (canvas *CanvasPane) Init(rsc *Resources, space sdl.Rect, tracks []string) {
	canvas.rsc = rsc
	canvas.Pos = space
	canvas.tracks = tracks
	canvas.menu.Init(rsc.renderer, space, []string{"+input", "+output", "+effect"}, rsc.TitleFont)

	canvas.menu.OnClick(func(entry *MenuEntry) {
		log.Println("Clicked: " + entry.Text)
		if entry.Text == "+input" {
			canvas.NewInput()
		} else if entry.Text == "+output" {
			canvas.NewOutput()
		} else if entry.Text == "+effect" {
			canvas.NewEffect()
		}
	})
}

func (canvas *CanvasPane) NewInput() {
	n := &Node{}
	n.Pos = canvas.menu.Pos
	n.Pos.W = 64
	n.Pos.H = 48
	n.color = hexcolor(0x5be33b)
	n.name = "input"
	n.label.Init(canvas.rsc.renderer, n.Pos, n.name, canvas.rsc.TitleFont, hexcolor(0x303030))

	n.menu.Init(canvas.rsc.renderer,
		n.Pos,
		canvas.tracks,
		canvas.rsc.TitleFont)
	canvas.nodes = append(canvas.nodes, n)

	n.menu.OnClick(func(entry *MenuEntry) {
		log.Println("Input clicked: " + entry.Text)
		n.args = make([]string, 1, 1)
		n.args[0] = entry.Text
		n.label.Text = entry.Text
		n.label.Update(canvas.rsc.renderer)
		if n.Pos.W < n.label.texwidth + 8 {
			n.Pos.W = n.label.texwidth + 8
		}
		if n.Pos.H < n.label.texheight + 8 {
			n.Pos.H = n.label.texheight + 8
		}
		n.label.Pos = n.Pos
	})
}

func (canvas *CanvasPane) NewOutput() {
	n := &Node{}
	n.Pos = canvas.menu.Pos
	n.Pos.W = 64
	n.Pos.H = 48
	n.color = hexcolor(0xff3015)
	n.name = "output"
	n.label.Init(canvas.rsc.renderer, n.Pos, n.name, canvas.rsc.TitleFont, hexcolor(0x303030))
	canvas.nodes = append(canvas.nodes, n)
}

func (canvas *CanvasPane) NewEffect() {
	n := &Node{}
	n.Pos = canvas.menu.Pos
	n.Pos.W = 64
	n.Pos.H = 48
	n.color = hexcolor(0xffe018)
	n.name = "(null-fx)"
	n.label.Init(canvas.rsc.renderer, n.Pos, n.name, canvas.rsc.TitleFont, hexcolor(0x303030))

	effects := make([]string, 0, 10)
	for i, h := range sox.GetEffectHandlers() {
		if i < 10 {
			effects = append(effects, h.Name())
		}
	}
	n.menu.Init(canvas.rsc.renderer, n.Pos, effects, canvas.rsc.TitleFont)
	canvas.nodes = append(canvas.nodes, n)

	n.menu.OnClick(func(entry *MenuEntry) {
		n.label.Text = entry.Text
		n.label.Update(canvas.rsc.renderer)
		if n.Pos.W < n.label.texwidth + 8 {
			n.Pos.W = n.label.texwidth + 8
		}
		if n.Pos.H < n.label.texheight + 8 {
			n.Pos.H = n.label.texheight + 8
		}
		n.label.Pos = n.Pos
	})
}

func (canvas *CanvasPane) Draw(rend *sdl.Renderer) {

	for _, n := range canvas.nodes {
		n.DrawLink(rend)
	}

	for _, n := range canvas.nodes {
		n.Draw(rend)
	}

	canvas.menu.Draw(rend)

	if canvas.new_link != nil {
		_, x, y := sdl.GetMouseState()
		n0 := canvas.nodes[*canvas.new_link]
		rend.SetDrawColor(hexcolor(0xffffff))
		pos := n0.GetPos()
		rend.DrawRect(&pos)
		rend.SetDrawColor(hexcolor(0x694ae9))
		rend.DrawLine(n0.GetPos().X + n0.GetPos().W/2, n0.GetPos().Y + n0.GetPos().H/2, int32(x), int32(y))
	}
}

func (canvas *CanvasPane) UpdateLayout(space sdl.Rect) {
	canvas.Pos.W = space.W
	canvas.Pos.H = space.H - canvas.Pos.Y
	canvas.menu.UpdateLayout(space)
}

func (canvas *CanvasPane) OnMouseMotionEvent(event *sdl.MouseMotionEvent) bool {
	if canvas.menu.Visible {
		canvas.menu.OnMouseMotionEvent(event)
	} else {
		for _, n := range canvas.nodes {
			n.OnMouseMotionEvent(event)
		}
	}
	return true
}

func (canvas *CanvasPane) OnMouseButtonEvent(event *sdl.MouseButtonEvent) bool {
	if event.State == sdl.RELEASED && canvas.new_link != nil {
		to := -1
		for i, n := range canvas.nodes {
			pos := n.GetPos()
			if pos.Contains(event.X, event.Y) {
				to = i
				break
			}
		}
		if to != -1 && to != *canvas.new_link && canvas.nodes[to].name != "input" {
			canvas.nodes[*canvas.new_link].next = canvas.nodes[to]
		}
		canvas.new_link = nil
	}

	if event.Button == sdl.BUTTON_RIGHT && event.State == sdl.PRESSED {
		if !canvas.menu.Visible {
			hitsbox := false
			for _, n := range canvas.nodes {
				p := n.GetPos()
				if p.Contains(event.X, event.Y) {
					hitsbox = true
					break
				}
			}
			if !hitsbox && canvas.Pos.Contains(event.X, event.Y) {
				canvas.menu.Show(event.X, event.Y)
				return false
			}
		} else {
			canvas.menu.Hide()
		}
	}
	if canvas.menu.Visible {
		canvas.menu.OnMouseButtonEvent(event)
	} else {
		lpress := event.State == sdl.PRESSED && event.Button == sdl.BUTTON_LEFT
		if ((sdl.GetModState() & sdl.KMOD_SHIFT) != 0) && lpress {
			from := -1
			for i, n := range canvas.nodes {
				pos := n.GetPos()
				if pos.Contains(event.X, event.Y) {
					from = i
					break
				}
			}
			if from >= 0 && canvas.nodes[from].name != "output" {
				log.Println("Start linking from node", from)
				canvas.new_link = &from
			}
		} else {
			for _, n := range canvas.nodes {
				if !n.OnMouseButtonEvent(event) {
					break
				}
			}
		}
	}
	return true
}

func (canvas *CanvasPane) BuildPlayChain() func() {
	if canvas.playing != nil && canvas.playing.finished {
		canvas.playing.Release()
		canvas.playing = nil
	}

	// find an input
	var start *Node
	for _, n := range canvas.nodes {
		if n.name == "input" {
			start = n
			break
		}
	}
	if start == nil || len(start.args) != 1 {
		log.Println("Nothing to play.")
		return func() {}
	}

	var stop *Node
	for n2 := start.next; n2 != nil; n2 = n2.next {
		if n2 == start {
			log.Println("Loop detected!")
			return func() {}
		} else if n2.name == "output" {
			stop = n2
		}
	}

	if stop == nil {
		log.Println("Nothing to play to.")
		return func() {}
	}

	in := sox.OpenRead(start.args[0])
	if in == nil {
		log.Println("Failed to open input file.")
		return func() {}
	}
	out := sox.OpenWrite("default", in.Signal(), nil, "alsa")
	if out == nil {
		log.Println("Failed to open output device.")
		return func() {}
	}
	chain := sox.CreateEffectsChain(in.Encoding(), out.Encoding())

	e := sox.CreateEffect(sox.FindEffect("input"))
	e.Options(in)
	chain.Add(e, in.Signal(), in.Signal())
	e.Release()

	e = sox.CreateEffect(sox.FindEffect("output"))
	e.Options(out)
	chain.Add(e, in.Signal(), in.Signal())
	e.Release()

	canvas.playing = &SoxChain{chain, in, out, false, false}

	return func() {
		canvas.playing.Flow()
	}
}

func (canvas *CanvasPane) Play() {

	fn := canvas.BuildPlayChain()
	go fn()
	//go canvas.playing.Flow()
}

func (canvas *CanvasPane) Stop() {
	if canvas.playing != nil {
		canvas.playing.interrupt = true
	}
}

func (screen *Screen) Init(space sdl.Rect, rsc *Resources, tracks []string) {
	screen.Pos = space
	screen.TopBar = &TopBar{}
	screen.TopBar.Init(sdl.Rect{space.X, space.Y, space.W, TOPBAR_HEIGHT})
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
	screen.Title.Init(rsc.renderer, sdl.Rect{screen.F2.Pos.X + screen.F2.Pos.W, space.Y, 1, TOPBAR_HEIGHT}, "canvas mode", rsc.TitleFont, rsc.TitleColor)
	screen.Play.Init(rsc.renderer, sdl.Rect{screen.Title.Pos.X + screen.Title.Pos.W, space.Y, 1, 1}, rsc.PlayButton)
	screen.Stop.Init(rsc.renderer, sdl.Rect{screen.Play.Pos.X + screen.Play.Pos.W, space.Y, 1, 1}, rsc.StopButton)

	screen.AddVisual(screen.TopBar)
	screen.AddLayout(screen.TopBar)

	screen.Canvas = &CanvasPane{}
	screen.Canvas.Init(rsc, sdl.Rect{space.X, space.Y + TOPBAR_HEIGHT, space.W, space.H - TOPBAR_HEIGHT}, tracks)
	screen.AddVisual(screen.Canvas)
	screen.AddLayout(screen.Canvas)

	screen.UpdateLayout(space)

	screen.F1.OnClick(func() {
		log.Println("Canvas Mode clicked!")
	})

	screen.F2.OnClick(func() {
		log.Println("Track Mode clicked!")
	})

	screen.Play.OnClick(func() {
		log.Println("Play clicked!")
		screen.Canvas.Play()
	})

	screen.Stop.OnClick(func() {
		log.Println("Stop clicked!")
		screen.Canvas.Stop()
	})

}

func (screen *Screen) UpdateAnimations(delta float64) {
	// TODO
}

func run_studio(window *sdl.Window, rend *sdl.Renderer, tracks []string) {
	runtime.LockOSThread()

	rsc := &Resources{}
	rsc.Load(rend)
	defer rsc.Free()

	w, h := window.GetSize()
	neww, newh := w, h
	screen := &Screen{}
	screen.Init(sdl.Rect{0, 0, int32(w), int32(h)}, rsc, tracks)
	stack := &InputStack{}
	stack.Add(screen.F1)
	stack.Add(screen.F2)
	stack.Add(screen.Play)
	stack.Add(screen.Stop)
	stack.Add(screen.Canvas)
	defer screen.Destroy()

	framerate := gfx.NewFramerate()
	framerate.SetFramerate(30)
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

		screen.UpdateAnimations(framerate.Delta())

		rend.SetDrawBlendMode(sdl.BLENDMODE_NONE)
		rend.SetDrawColor(rsc.BackgroundColor)
		rend.Clear()
		rend.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
		screen.Draw(rend)
		rend.Present()
		framerate.FramerateDelay()
	}
}
