package model

import "math"

type Vector struct {
	x float64
	y float64
}

func NewVector(x, y float64) Vector {
	return Vector{x, y}
}

// modification methods.
func (v *Vector) Add(v2 Vector) {
	v.x += v2.x
	v.y += v2.y
}

func (v *Vector) Sub(v2 Vector) {
	v.x -= v2.x
	v.y -= v2.y
}

func (v *Vector) Scale(s float64) {
	v.x *= s
	v.y *= s
}

func (v *Vector) Div(d float64) {
	v.x /= d
	v.y /= d
}

func (v *Vector) Limit(max float64) {
	if v.Magnitude() > max {
		v.Normalize()
		v.Scale(max)
	}
}

func (v *Vector) Normalize() {
	mag := v.Magnitude()
	if mag > 0 {
		v.Div(mag)
	}
}

func (v *Vector) SetMag(mag float64) {
	v.Normalize()
	v.Scale(mag)
}

// non-modifier methods.
func (v *Vector) SoftAdd(v2 Vector) Vector {
	return NewVector(v.x+v2.x, v.y+v2.y)
}

func (v *Vector) SoftSub(v2 Vector) Vector {
	return NewVector(v.x-v2.x, v.y-v2.y)
}

func (v *Vector) SoftScale(s float64) Vector {
	return NewVector(v.x*s, v.y*s)
}

func (v *Vector) SoftDiv(d float64) Vector {
	return NewVector(v.x/d, v.y/d)
}

func (v *Vector) Magnitude() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v *Vector) Dist(v2 Vector) float64 {
	x := v.x - v2.x
	y := v.y - v2.y
	return math.Sqrt(x*x + y*y)
}
