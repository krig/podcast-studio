package main

import (
	"time"
	"github.com/krig/Go-SDL2/sdl"
	"github.com/krig/Go-SDL2/ttf"
	"log"
)

type Track interface {
}

type Sample interface {
}

type Node interface {
}

type Link interface {
}

type Studio struct {
	mode int
}

type Model struct {
	tracks *[]Track
	samples *[]Sample
	nodes *[]Node
	links *[]Link
	playback_state uint64
	position uint64
}

type CanvasView struct {
}

type TrackView struct {
}

type StudioController struct {
}

func update_playback() {
}

func update_animations() {
}

func player_loop() {
	for {
		update_playback()
		update_animations()
		time.Sleep(10 * 1e6) // 10 ms
	}
}

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

func main() {
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		log.Fatal(sdl.GetError())
	}
	defer sdl.Quit()

	window, rend := sdl.CreateWindowAndRenderer(640, 480, sdl.WINDOW_SHOWN |
		sdl.RENDERER_ACCELERATED |
		sdl.RENDERER_PRESENTVSYNC)
	if (window == nil) || (rend == nil) {
		log.Fatal(sdl.GetError())
	}
	defer window.Destroy()

	if ttf.Init() != 0 {
		log.Fatal(sdl.GetError())
	}
	defer ttf.Quit()

	window.SetTitle("Podcast Studio")


	garoa := ttf.OpenFont("data/GaroaHackerClubeBold.otf", 10)
	defer garoa.Close()

	txt_tex, txt_w, txt_h := RenderTextToTexture(rend, garoa, "PODCAST STUDIO", sdl.Color{0xFF, 0xFF, 0xFF, 0xFF})
	defer txt_tex.Destroy()

	running := true
	event := &sdl.Event{}
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

		rend.SetDrawColor(sdl.Color{0x30, 0x30, 0x30, 0xFF})
		rend.Clear()
		rend.SetDrawColor(sdl.Color{0xFF, 0x1F, 0x69, 0xFF})
		w, h := window.GetSize()
		rend.DrawLine(w/2, h/2 - 100, w/2 + 100, h/2 + 100)
		rend.DrawLine(w/2 - 100, h/2 + 100, w/2 + 100, h/2 + 100)
		rend.DrawLine(w/2, h/2 - 100, w/2 - 100, h/2 + 100)

		rend.Copy(txt_tex, nil, &sdl.Rect{int32(w/2 - txt_w/2), int32(h/2 + 100 + txt_h), int32(txt_w), int32(txt_h)})

		rend.Present()
		sdl.Delay(10);
	}
}