package main

import (
	"github.com/krig/Go-SDL2/sdl"
	"github.com/krig/Go-SDL2/ttf"
	"github.com/krig/go-sox"
	"log"
)

func main() {
	if !sox.Init() {
		log.Fatal("Failed to init sox")
	}
	defer sox.Quit()

	if sdl.Init(sdl.INIT_NOPARACHUTE|sdl.INIT_VIDEO|sdl.INIT_EVENTS) != 0 {
		log.Fatal(sdl.GetError())
	}
	defer sdl.Quit()
	window, renderer := sdl.CreateWindowAndRenderer(640, 480, sdl.WINDOW_SHOWN |
		sdl.WINDOW_OPENGL |
		sdl.RENDERER_ACCELERATED |
		sdl.RENDERER_PRESENTVSYNC)
	if (window == nil) || (renderer == nil) {
		log.Fatal(sdl.GetError())
	}
	defer window.Destroy()
	defer renderer.Destroy()
	if ttf.Init() != 0 {
		log.Fatal(sdl.GetError())
	}
	defer ttf.Quit()
	window.SetTitle("Podcast Studio")
	log.Println("Video Driver:", sdl.GetCurrentVideoDriver())

	run_studio(window, renderer)
}