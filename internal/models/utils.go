package model

import "math"

func Map(value, start1, stop1, start2, stop2 float64) float64 {
	newval := (value-start1)/(stop1-start1)*(stop2-start2) + start2
	return newval
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
