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
}
