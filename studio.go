package main

import (
	"time"
	"github.com/krig/Go-SDL2/sdl"
	"github.com/krig/Go-SDL2/ttf"
	"log"
	"math"
	"math/rand"
)

func RenderTextToTexture(r *sdl.Renderer, f *ttf.Font, text string, color sdl.Color) (*sdl.Texture, int, int) {
	textw, texth, err := f.SizeText(text)
	if err != nil {
		log.Fatal(err)
	}
	txt_surface := f.RenderText_Blended(text, color)
	txt_tex := r.CreateTextureFromSurface(txt_surface)
	txt_surface.Free()
	return txt_tex, textw, texth
}

func run_studio(window *sdl.Window, rend *sdl.Renderer) {
	garoa := ttf.OpenFont("data/GaroaHackerClubeBold.otf", 10)
	defer garoa.Close()

	txt_tex, txt_w, txt_h := RenderTextToTexture(rend, garoa, "PODCAST STUDIO", sdl.Color{0xFF, 0xFF, 0xFF, 0xFF})
	defer txt_tex.Destroy()

	running := true
	event := &sdl.Event{}
	wobble := 1.0
	dim := 100.0
	state := 0.0
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	rend.SetDrawColor(sdl.Color{0x30, 0x30, 0x30, 0xFF})
	rend.Clear()

	var t uint64 = 0

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

		w, h := window.GetSize()
		wobble = rnd.Float64()
		state += 0.02
		dim = 100.0 + math.Sin(state) * 15.0 + wobble

		rend.SetDrawColor(sdl.Color{0x30, 0x30, 0x30, 0xFF})
		rend.Clear()
		rend.SetDrawColor(sdl.Color{0xFF, 0x1F, 0x69, 0xFF})
		rend.DrawLine(w/2, h/2 - int(dim), w/2 + int(dim), h/2 + int(dim))
		rend.DrawLine(w/2 - int(dim), h/2 + int(dim), w/2 + int(dim), h/2 + int(dim))
		rend.DrawLine(w/2, h/2 - int(dim), w/2 - int(dim), h/2 + int(dim))

		rend.Copy(txt_tex, nil, &sdl.Rect{int32(w/2 - txt_w/2), int32(h/2 + int(dim) + txt_h), int32(txt_w), int32(txt_h)})
		rend.Present()
		sdl.Delay(5);

		t += 1
	}
}