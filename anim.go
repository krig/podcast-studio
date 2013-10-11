package main

import (
	"math"

	"github.com/krig/Go-SDL2/sdl"
)

// Stuff related to UI animation.
// Position animation, color animation...?

func interp(current, target, dt float64) float64 {
	if target - current < 0.00001 {
		target = current
	}
	return current + (target - current) * dt
}

// FloatColor holds colors as floats for smoother animation
type FloatColor struct {
	R, G, B, A float64
}

func makeFloatColor(clr sdl.Color) FloatColor {
	return FloatColor{float64(clr.R) / 255.0, float64(clr.G) / 255.0, float64(clr.B) / 255.0, float64(clr.A) / 255.0}
}

func makeSDLColor(clr FloatColor) sdl.Color {
	return sdl.Color{uint8(math.Floor(clr.R + 0.5)), uint8(math.Floor(clr.G + 0.5)), uint8(math.Floor(clr.B + 0.5)), uint8(math.Floor(clr.A + 0.5))}
}

// ColorAnimation animates a color
type ColorAnimation struct {
	target FloatColor
	current FloatColor
	speed float64
}

func (anim *ColorAnimation) Init(start, target sdl.Color, speed float64) {
	anim.current = makeFloatColor(start)
	anim.target = makeFloatColor(target)
	anim.speed = speed
}

func (anim *ColorAnimation) SetTarget(target sdl.Color) {
	anim.target = makeFloatColor(target)
}

func (anim *ColorAnimation) Finish() {
	anim.current = anim.target
}

func (anim *ColorAnimation) Update(dt float64) {
	sp := dt * anim.speed
	anim.current.R = interp(anim.current.R, anim.target.R, sp)
	anim.current.G = interp(anim.current.G, anim.target.G, sp)
	anim.current.B = interp(anim.current.B, anim.target.B, sp)
	anim.current.A = interp(anim.current.A, anim.target.A, sp)
}

func (anim *ColorAnimation) Get() sdl.Color {
	return makeSDLColor(anim.current)
}

