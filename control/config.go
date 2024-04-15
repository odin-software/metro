package control

import (
	"time"
)

type Config struct {
	LogsDirectory       string
	LoopDuration        time.Duration
	LoopDurationOffset  time.Duration
	LoopStartingState   int
	PortCity            int
	PortEvents          int
	PortReporter        int
	PortVirtual         int
	ReflexDuration      time.Duration
	StdLogs             bool
	TerminalMapDuration time.Duration
	TerminalMapEnabled  bool
	TrainWaitInStation  time.Duration
	WSLogsDuration      time.Duration
	WSTrainDuration     time.Duration
}

var DefaultConfig = Config{
	LogsDirectory:       "logs/",
	LoopDuration:        20 * time.Millisecond,
	LoopDurationOffset:  -1 * time.Millisecond,
	LoopStartingState:   1,
	PortCity:            2221,
	PortEvents:          2223,
	PortReporter:        2222,
	PortVirtual:         2224,
	ReflexDuration:      200 * time.Millisecond,
	StdLogs:             false,
	TerminalMapDuration: 800 * time.Millisecond,
	TerminalMapEnabled:  false,
	TrainWaitInStation:  3000 * time.Millisecond,
	WSLogsDuration:      4000 * time.Millisecond,
	WSTrainDuration:     200 * time.Millisecond,
}
