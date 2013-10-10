package main

import (
	"flag"
	"log"

	"github.com/krig/Go-SDL2/sdl"
	"github.com/krig/Go-SDL2/ttf"
	"github.com/krig/go-sox"
)

func main() {
	// Parse command line
	flag.Parse()
	session := "default"
	tracks := []string{}
	if flag.NArg() > 0 {
		session = flag.Arg(0)
	}
	if flag.NArg() > 1 {
		tracks = flag.Args()[1:]
	}
	log.Println("Session: ", session)
	log.Println("Tracks: ", tracks)

	// Init libSoX and SDL
	if !sox.Init() {
		log.Fatal("Failed to init sox")
	}
	defer sox.Quit()

	if sdl.Init(sdl.INIT_NOPARACHUTE|sdl.INIT_VIDEO|sdl.INIT_EVENTS) != 0 {
		log.Fatal(sdl.GetError())
	}
	defer sdl.Quit()

	window, renderer := sdl.CreateWindowAndRenderer(640, 480,
		sdl.WINDOW_SHOWN | sdl.WINDOW_OPENGL | sdl.RENDERER_ACCELERATED | sdl.RENDERER_PRESENTVSYNC)
	if (window == nil) || (renderer == nil) {
		log.Fatal(sdl.GetError())
	}
	if ttf.Init() != 0 {
		log.Fatal(sdl.GetError())
	}
	defer window.Destroy()
	defer renderer.Destroy()
	defer ttf.Quit()
	window.SetTitle("Podcast Studio")
	log.Println("Video Driver:", sdl.GetCurrentVideoDriver())
	// Jump to studio.go
	run_studio(window, renderer)
}