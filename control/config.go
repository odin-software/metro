package control

import (
	"time"
)

type Config struct {
	DisplayScreenWidth  int
	DisplayScreenHeight int
	LogsDirectory       string
	LoopDuration        time.Duration
	LoopDurationOffset  time.Duration
	LoopStartingState   int
	ReflexDuration      time.Duration
	StdLogs             bool
	TrainWaitInStation  time.Duration
}

var DefaultConfig = Config{
	DisplayScreenWidth:  800,
	DisplayScreenHeight: 600,
	LogsDirectory:       "logs/",
	LoopDuration:        time.Second / 60,
	LoopDurationOffset:  -1 * time.Millisecond,
	LoopStartingState:   1,
	ReflexDuration:      200 * time.Millisecond,
	StdLogs:             false,
	TrainWaitInStation:  3 * time.Second,
}
