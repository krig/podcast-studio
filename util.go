// Utility functions for Podcast Studio
package main

import (
	"github.com/krig/Go-SDL2/sdl"
)

func hexcolor(code uint32) sdl.Color {
	return sdl.Color{uint8((code >> 16) & 0xff), uint8((code >> 8) & 0xff), uint8(code & 0xff), 0xff}
}
