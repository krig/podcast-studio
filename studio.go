package main

import (
	"time"
	"github.com/krig/Go-SDL2/sdl"
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

func main() {
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		log.Fatal(sdl.GetError())
	}

	window, rend := sdl.CreateWindowAndRenderer(640, 480, sdl.WINDOW_SHOWN |
		sdl.RENDERER_ACCELERATED |
		sdl.RENDERER_PRESENTVSYNC)

	if (window == nil) || (rend == nil) {
		log.Fatal(sdl.GetError())
	}

	window.SetTitle("Podcast Studio")

	running := true
	event := &sdl.Event{}
	for running {
		for event.Poll() {
			switch e := event.Get().(type) {
			case sdl.QuitEvent:
				running = false

			case sdl.KeyboardEvent:
				if e.Keysym.Sym == sdl.K_ESCAPE {
					running = false
				}
			}
		}

		rend.SetDrawColor(sdl.Color{0x30, 0x30, 0x30, 0xFF, 0x00})
		rend.Clear()
		rend.Present()
		sdl.Delay(10);
	}

	window.Destroy()
	sdl.Quit()
}