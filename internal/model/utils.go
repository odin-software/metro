package model

func Map(value, start1, stop1, start2, stop2 float64) float64 {
	newval := (value-start1)/(stop1-start1)*(stop2-start2) + start2
	return newval
}
