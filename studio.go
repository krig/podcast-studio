package main

import (
	"github.com/krig/Go-SDL2/sdl"
	"github.com/krig/Go-SDL2/ttf"
)

type Resources struct {
	renderer *sdl.Renderer
	TitleFont *ttf.Font
	PlayButton *sdl.Texture
	LEDButton *sdl.Texture
}

func (r *Resources) Load(rend *sdl.Renderer) {
	r.renderer = rend

	r.TitleFont = ttf.OpenFont("data/GaroaHackerClubeBold.otf", 10)

	playbutton := sdl.Load("data/testbutton.png")
	r.PlayButton = rend.CreateTextureFromSurface(playbutton)
	playbutton.Free()

	ledbutton := sdl.Load("data/testbutton_blank.png")
	r.LEDButton = rend.CreateTextureFromSurface(ledbutton)
	ledbutton.Free()
}

func (r *Resources) Free() {
	r.TitleFont.Close()
	r.PlayButton.Destroy()
}

func run_studio(window *sdl.Window, rend *sdl.Renderer) {
	rsc := &Resources{}
	rsc.Load(rend)
	defer rsc.Free()

	w, h := window.GetSize()
	neww, newh := w, h
	screen := &Screen{}
	screen.Init(sdl.Rect{0, 0, int32(w), int32(h)}, rsc)
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
			}
		}

		neww, newh = window.GetSize()
		if neww != w || newh != h {
			w = neww
			h = newh
			screen.UpdateLayout(sdl.Rect{0, 0, int32(w), int32(h)})
		}

		rend.SetDrawColor(hexcolor(0x303030))
		rend.Clear()
		rend.SetDrawColor(hexcolor(0xffffff))
		screen.Draw(rend)
		rend.Present()
	}
}
