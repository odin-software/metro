package control

import (
	"time"
)

type Config struct {
	DisplayScreenWidth   int
	DisplayScreenHeight  int
	DisplayMonitor       int // Monitor index (0 = primary, 1 = second monitor, etc.)
	LogsDirectory        string
	LoopDuration         time.Duration
	LoopDurationOffset   time.Duration
	LoopStartingState    int
	ReflexDuration       time.Duration
	StdLogs              bool
	TrainWaitInStation   time.Duration
	TenjinEnabled        bool
	TenjinTickRate       time.Duration
	PassengerSpawnRate   time.Duration
	PassengersPerStation int

	// Real-world metrics scaling
	PixelsPerMeter      float64 // Scale factor: 1 pixel = X meters
	SimulationSpeed     float64 // Multiplier for time (1.0 = real-time, 2.0 = 2x speed)

	// Simulation time clock
	SimulationStartHour int     // Starting hour (0-23), e.g., 8 for 8:00 AM
	SimulationStartMin  int     // Starting minute (0-59)
}

var DefaultConfig = Config{
	DisplayScreenWidth:   800,
	DisplayScreenHeight:  600,
	DisplayMonitor:       1, // 0 = primary, 1 = second monitor
	LogsDirectory:        "logs/",
	LoopDuration:         time.Second / 60,
	LoopDurationOffset:   -1 * time.Millisecond,
	LoopStartingState:    1,
	ReflexDuration:       2 * time.Second,
	StdLogs:              true,
	TrainWaitInStation:   5 * time.Second,
	TenjinEnabled:        true,
	TenjinTickRate:       time.Second,
	PassengerSpawnRate:   5 * time.Second,
	PassengersPerStation: 3,

	// Real-world scaling: 1 pixel = 100 meters (map is ~70km x 50km)
	PixelsPerMeter:      0.01,  // 1 pixel = 100 meters
	SimulationSpeed:     1.0,   // 1.0 = real-time

	// Simulation starts at 8:00 AM
	SimulationStartHour: 8,
	SimulationStartMin:  0,
}
