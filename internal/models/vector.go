package model

import "math"

type Vector struct {
	X float64
	Y float64
}

func NewVector(x, y float64) Vector {
	return Vector{x, y}
}

// type Points = []Vector

// // Points methods.
// func (p Points) Reverse() Points {
// }

// modification methods.
func (v *Vector) Add(v2 Vector) {
	v.X += v2.X
	v.Y += v2.Y
}

func (v *Vector) Sub(v2 Vector) {
	v.X -= v2.X
	v.Y -= v2.Y
}

func (v *Vector) Scale(s float64) {
	v.X *= s
	v.Y *= s
}

func (v *Vector) Div(d float64) {
	v.X /= d
	v.Y /= d
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
	return NewVector(v.X+v2.X, v.Y+v2.Y)
}

func (v *Vector) SoftSub(v2 Vector) Vector {
	return NewVector(v.X-v2.X, v.Y-v2.Y)
}

func (v *Vector) SoftScale(s float64) Vector {
	return NewVector(v.X*s, v.Y*s)
}

func (v *Vector) SoftDiv(d float64) Vector {
	return NewVector(v.X/d, v.Y/d)
}

func (v *Vector) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vector) Dist(v2 Vector) float64 {
	X := v.X - v2.X
	Y := v.Y - v2.Y
	return math.Sqrt(X*X + Y*Y)
}

func (v *Vector) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

func (v *Vector) CloseTo(x, y float64, precision int) bool {
	return v.Dist(NewVector(x, y)) <= float64(precision)
}

func (v *Vector) Copy() Vector {
	return NewVector(v.X, v.Y)
}
