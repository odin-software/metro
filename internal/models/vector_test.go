package models

import "testing"

func TestNewVector(t *testing.T) {
	vec := NewVector(2.0, 1.3)

	if vec.X != 2.0 {
		t.Fatal("The vector was not created correctly.")
	}
}

func TestAdd(t *testing.T) {
	vec := NewVector(2.0, 1.3)
	vec2 := NewVector(8.0, 1.1)
	vec.Add(vec2)

	if RoundFloat(vec.X, 2) != 10.0 {
		t.Fatal("The vectors were added incorrectly.")
	}
	if RoundFloat(vec.Y, 2) != 2.4 {
		t.Fatal("The vectors were added incorrectly.")
	}
}

func TestSub(t *testing.T) {
	vec := NewVector(2.0, 1.3)
	vec2 := NewVector(8.0, 1.1)
	vec.Sub(vec2)

	if RoundFloat(vec.X, 2) != -6.0 {
		t.Fatal("The vectors were substracted incorrectly.")
	}
	if RoundFloat(vec.Y, 2) != 0.2 {
		t.Fatal("The vectors were substracted incorrectly.")
	}
}

func TestScale(t *testing.T) {
	vec := NewVector(2.0, 1.3)
	vec.Scale(3)

	if RoundFloat(vec.X, 2) != 6.0 {
		t.Fatal("The vector was scaled incorrectly.")
	}
	if RoundFloat(vec.Y, 2) != 3.9 {
		t.Fatal("The vector was scaled incorrectly.")
	}
}

func TestDiv(t *testing.T) {
	vec := NewVector(6.0, 3.3)
	vec.Div(3)

	if RoundFloat(vec.X, 2) != 2.0 {
		t.Fatal("The vector was divided incorrectly.")
	}
	if RoundFloat(vec.Y, 2) != 1.1 {
		t.Fatal("The vector was divided incorrectly.")
	}
}

func TestLimit(t *testing.T) {
	vec := NewVector(1.0, 0.0)
	vec.Limit(3)

	if RoundFloat(vec.X, 2) != 1.0 {
		t.Fatal("Shouldn't change.")
	}
	vec.Scale(8)
	if RoundFloat(vec.X, 2) != 8.0 {
		t.Fatal("Value should have scaled.")
	}
	vec.Limit(3)
	if RoundFloat(vec.X, 2) != 3.0 {
		t.Fatal("The limit did not work.")
	}
}

func TestNormalize(t *testing.T) {
	vec := NewVector(4.0, 0.0)
	vec.Normalize()

	if RoundFloat(vec.X, 2) != 1.0 {
		t.Fatal("Should have turn it into a unit vector.")
	}
}

func TestSetMag(t *testing.T) {
	vec := NewVector(2.9, 9.1)
	vec.SetMag(4)

	if vec.Magnitude() != 4.0 {
		t.Fatal("The magnitude changed didn't change.")
	}
}

func TestSoftAdd(t *testing.T) {
	v1 := NewVector(2.9, 9.1)
	v2 := NewVector(2.9, 9.1)
	v3 := v1.SoftAdd(v2)

	if RoundFloat(v3.X, 2) != 5.8 {
		t.Fatal("The new vector has an incorrect value.")
	}
	if RoundFloat(v3.Y, 2) != 18.2 {
		t.Fatal("The new vector has an incorrect value.")
	}
}

func TestSoftSub(t *testing.T) {
	v1 := NewVector(2.9, 9.1)
	v2 := NewVector(2.8, 9.0)
	v3 := v1.SoftSub(v2)

	if RoundFloat(v3.X, 2) != 0.1 {
		t.Fatal("The new vector has an incorrect value.")
	}
	if RoundFloat(v3.Y, 2) != 0.1 {
		t.Fatal("The new vector has an incorrect value.")
	}
}

func TestSoftScale(t *testing.T) {
	v1 := NewVector(2.0, 1.2)
	v2 := v1.SoftScale(3)

	if RoundFloat(v2.X, 2) != 6.0 {
		t.Fatal("The new vector has an incorrect value.")
	}
	if RoundFloat(v2.Y, 2) != 3.6 {
		t.Fatal("The new vector has an incorrect value.")
	}
}

func TestSoftDiv(t *testing.T) {
	v1 := NewVector(9.0, 3.3)
	v2 := v1.SoftDiv(3)

	if RoundFloat(v2.X, 2) != 3.0 {
		t.Fatal("The new vector has an incorrect value.")
	}
	if RoundFloat(v2.Y, 2) != 1.1 {
		t.Fatal("The new vector has an incorrect value.")
	}
}

func TestMagnitude(t *testing.T) {
	v1 := NewVector(2.0, 5.0)

	if RoundFloat(v1.Magnitude(), 2) != 5.39 {
		t.Fatal("The magnitude value is incorrect.")
	}
}

func TestDistance(t *testing.T) {
	v1 := NewVector(1.0, 1.0)
	v2 := NewVector(4.0, 5.0)
	dist1 := v1.Dist(v2)
	dist2 := v2.Dist(v1)

	if RoundFloat(dist1, 2) != 5.00 {
		t.Fatal("The distance is wrong.")
	}
	if dist1 != dist2 {
		t.Fatal("The distance is wrong.")
	}
}
