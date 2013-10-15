package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/krig/Go-SDL2/sdl"
	"github.com/krig/Go-SDL2/ttf"
	"github.com/krig/go-sox"
	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/glib"
)

func openFileDialog(callback func(filename string)) {
	filechooserdialog := gtk.NewFileChooserDialog(
		"Choose File",
		nil,
		gtk.FILE_CHOOSER_ACTION_OPEN,
		gtk.STOCK_OK,
		gtk.RESPONSE_ACCEPT)
	filter := gtk.NewFileFilter()
	filter.AddPattern("*.wav")
	filter.AddPattern("*.mp3")
	filechooserdialog.AddFilter(filter)
	filechooserdialog.Response(func() {
		callback(filechooserdialog.GetFilename())
		filechooserdialog.Destroy()
	})

	filechooserdialog.Show()
}

func main() {
	runtime.LockOSThread()

	// Parse command line
	flag.Parse()
	tracks := []string{}
	if flag.NArg() > 0 {
		tracks = flag.Args()[0:]
	}
	log.Println("Tracks: ", tracks)

	// Init libSoX and SDL
	if !sox.Init() {
		log.Fatal("Failed to init sox")
	}
	defer sox.Quit()

	einfo := sox.GetEncodingsInfo()
	log.Println("Supported encodings:")
	for _, e := range einfo {
		log.Printf("%s: %s (%x)\n", e.Name, e.Desc, e.Flags)
	}

	effects := sox.GetEffectHandlers()
	log.Println("Supported effects:")
	for _, e := range effects {
		log.Printf("%s: %s (%x)\n", e.Name(), e.Usage(), e.Flags())
	}

	if sdl.Init(sdl.INIT_NOPARACHUTE|sdl.INIT_VIDEO|sdl.INIT_EVENTS) != 0 {
		log.Fatal(sdl.GetError())
	}
	defer sdl.Quit()

	gtk.Init(nil)

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
	//run_studio(window, renderer, tracks)
	screen := studioSetup(window, renderer, tracks)
	defer screen.rsc.Free()
	defer screen.Destroy()

	loop := glib.NewMainLoop(nil, false)

	glib.IdleAdd(func() bool {
		ret := studioUpdate(window, renderer, screen)
		if !ret {
			loop.Quit()
			return false
		}
		return true
	})

	loop.Run()
}