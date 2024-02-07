package model

import "math"

type Vector struct {
	x float64
	y float64
}

func NewVector(x, y float64) *Vector {
	return &Vector{x, y}
}

func (v *Vector) Magnitude() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v *Vector) Add(v2 *Vector) *Vector {
	return &Vector{v.x + v2.x, v.y + v2.y}
}

func (v *Vector) Dist(v2 *Vector) float64 {
	x := v.x - v2.x
	y := v.y - v2.y
	return math.Sqrt(x*x + y*y)
}

func (v *Vector) Scale(s float64) *Vector {
	return &Vector{v.x * s, v.y * s}
}

func (v *Vector) Sub(v2 *Vector) *Vector {
	return &Vector{v.x - v2.x, v.y - v2.y}
}

func (v *Vector) Unit(v2 *Vector) *Vector {
	return v.Scale(1 / v.Dist(v2))
}
