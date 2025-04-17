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
	TerminalMapDuration time.Duration
	TerminalMapEnabled  bool
	TrainWaitInStation  time.Duration
}

var DefaultConfig = Config{
	DisplayScreenWidth:  500,
	DisplayScreenHeight: 500,
	LogsDirectory:       "logs/",
	LoopDuration:        20 * time.Millisecond,
	LoopDurationOffset:  -1 * time.Millisecond,
	LoopStartingState:   1,
	ReflexDuration:      200 * time.Millisecond,
	StdLogs:             false,
	TerminalMapDuration: 800 * time.Millisecond,
	TerminalMapEnabled:  false,
	TrainWaitInStation:  3000 * time.Millisecond,
}
