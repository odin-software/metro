package models

import (
	"fmt"

	"github.com/odin-software/metro/control"
)

// RealWorldMetrics provides conversion functions for realistic distance/speed/time
type RealWorldMetrics struct{}

// PixelsToMeters converts pixel distance to meters
func PixelsToMeters(pixels float64) float64 {
	return pixels / control.DefaultConfig.PixelsPerMeter
}

// MetersToPixels converts meters to pixel distance
func MetersToPixels(meters float64) float64 {
	return meters * control.DefaultConfig.PixelsPerMeter
}

// PixelsToKilometers converts pixel distance to kilometers
func PixelsToKilometers(pixels float64) float64 {
	return PixelsToMeters(pixels) / 1000.0
}

// FormatDistance formats a pixel distance as a human-readable string
// Returns meters for < 1km, kilometers otherwise
func FormatDistance(pixels float64) string {
	meters := PixelsToMeters(pixels)
	if meters < 1000 {
		return fmt.Sprintf("%.0f m", meters)
	}
	return fmt.Sprintf("%.2f km", meters/1000.0)
}

// PixelSpeedToKmPerHour converts pixel/tick speed to km/h
// Takes into account the loop duration (ticks per second)
func PixelSpeedToKmPerHour(pixelsPerTick float64) float64 {
	// Convert to pixels per second
	ticksPerSecond := 1.0 / control.DefaultConfig.LoopDuration.Seconds()
	pixelsPerSecond := pixelsPerTick * ticksPerSecond

	// Convert to meters per second
	metersPerSecond := PixelsToMeters(pixelsPerSecond)

	// Convert to km/h
	kmPerHour := metersPerSecond * 3.6

	return kmPerHour
}

// KmPerHourToPixelSpeed converts km/h to pixel/tick speed
// Takes into account the loop duration (ticks per second)
func KmPerHourToPixelSpeed(kmPerHour float64) float64 {
	// Convert km/h to m/s
	metersPerSecond := kmPerHour / 3.6

	// Convert to pixels per second
	pixelsPerSecond := MetersToPixels(metersPerSecond)

	// Convert to pixels per tick
	ticksPerSecond := 1.0 / control.DefaultConfig.LoopDuration.Seconds()
	pixelsPerTick := pixelsPerSecond / ticksPerSecond

	return pixelsPerTick
}

// FormatSpeed formats a pixel/tick speed as km/h
func FormatSpeed(pixelsPerTick float64) string {
	kmh := PixelSpeedToKmPerHour(pixelsPerTick)
	return fmt.Sprintf("%.1f km/h", kmh)
}

// CalculateJourneyTime estimates journey time between two points at given speed
// Returns time in seconds
func CalculateJourneyTime(distance float64, speedKmH float64) float64 {
	if speedKmH == 0 {
		return 0
	}
	distanceKm := PixelsToKilometers(distance)
	timeHours := distanceKm / speedKmH
	timeSeconds := timeHours * 3600
	return timeSeconds
}
