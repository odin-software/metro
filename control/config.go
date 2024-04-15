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
	TerminalMapDuration time.Duration
	TerminalMapEnabled  bool
	TrainLogs           bool
	TrainWaitInStation  time.Duration
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
	TerminalMapDuration: 800 * time.Millisecond,
	TerminalMapEnabled:  false,
	TrainLogs:           false,
	TrainWaitInStation:  3000 * time.Millisecond,
	WSTrainDuration:     200 * time.Millisecond,
}
