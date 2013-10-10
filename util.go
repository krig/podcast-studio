package main

import (
	"github.com/krig/Go-SDL2/sdl"
)

func hexcolor(code uint32) sdl.Color {
	return sdl.Color{uint8((code >> 16) & 0xff), uint8((code >> 8) & 0xff), uint8(code & 0xff), 0xff}
}

// uint8_clamp returns the given number, clamped to the 0-255 range
func uint8_clamp(num int) uint8 {
	if num > 0xff {
		return 0xff
	} else if num < 0 {
		return 0
	} else {
		return uint8(num)
	}
}

// lighten the given color by ye much
func lighten(clr sdl.Color, ye int) sdl.Color {
	return sdl.Color{uint8_clamp(int(clr.R) + ye),
		uint8_clamp(int(clr.G) + ye),
		uint8_clamp(int(clr.B) + ye),
		uint8_clamp(int(clr.A) + ye)}
}

// darken the given color by ye much
func darken(clr sdl.Color, ye int) sdl.Color {
	return sdl.Color{uint8_clamp(int(clr.R) - ye),
		uint8_clamp(int(clr.G) - ye),
		uint8_clamp(int(clr.B) - ye),
		uint8_clamp(int(clr.A) - ye)}
}