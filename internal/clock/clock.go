package clock

import (
	"fmt"
	"sync"
	"time"

	"github.com/odin-software/metro/control"
)

// SimulationClock tracks the current time in the simulation
type SimulationClock struct {
	startTime       time.Time     // When simulation started (real time)
	simulationStart int           // Seconds since midnight when sim starts (e.g., 8:00 AM = 28800)
	elapsedSeconds  float64       // Elapsed simulation time in seconds
	mutex           sync.RWMutex
}

// NewSimulationClock creates a new simulation clock
func NewSimulationClock() *SimulationClock {
	startHour := control.DefaultConfig.SimulationStartHour
	startMin := control.DefaultConfig.SimulationStartMin

	// Convert start time to seconds since midnight
	simulationStart := startHour*3600 + startMin*60

	return &SimulationClock{
		startTime:       time.Now(),
		simulationStart: simulationStart,
		elapsedSeconds:  0,
	}
}

// Update increments the simulation time by the loop duration
// Should be called every tick (60 times per second)
func (c *SimulationClock) Update() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Increment by loop duration, scaled by simulation speed
	increment := control.DefaultConfig.LoopDuration.Seconds() * control.DefaultConfig.SimulationSpeed
	c.elapsedSeconds += increment
}

// GetCurrentTimeOfDay returns the current time of day in seconds since midnight
// Wraps around after 24 hours (86400 seconds)
func (c *SimulationClock) GetCurrentTimeOfDay() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	currentSeconds := c.simulationStart + int(c.elapsedSeconds)

	// Wrap around at 24 hours (86400 seconds in a day)
	return currentSeconds % 86400
}

// GetCurrentTime returns the current simulation time as a formatted string (HH:MM:SS)
func (c *SimulationClock) GetCurrentTime() string {
	seconds := c.GetCurrentTimeOfDay()
	return FormatSecondsAsTime(seconds)
}

// GetElapsedSeconds returns total elapsed simulation time in seconds
func (c *SimulationClock) GetElapsedSeconds() float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.elapsedSeconds
}

// FormatSecondsAsTime converts seconds since midnight to HH:MM:SS format
func FormatSecondsAsTime(seconds int) string {
	hours := (seconds / 3600) % 24
	minutes := (seconds % 3600) / 60
	secs := seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}

// TimeToSeconds converts HH:MM:SS to seconds since midnight
func TimeToSeconds(hours, minutes, seconds int) int {
	return hours*3600 + minutes*60 + seconds
}

// ParseTimeString parses "HH:MM:SS" string to seconds since midnight
func ParseTimeString(timeStr string) (int, error) {
	var hours, minutes, seconds int
	_, err := fmt.Sscanf(timeStr, "%d:%d:%d", &hours, &minutes, &seconds)
	if err != nil {
		return 0, err
	}
	return TimeToSeconds(hours, minutes, seconds), nil
}
