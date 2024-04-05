package control

import (
	"strconv"
	"time"

	"github.com/odin-software/metro/internal/models"
)

type Config struct {
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
	WSTrainDuration     time.Duration
}

var DefaultConfig = Config{
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
	WSTrainDuration:     200 * time.Millisecond,
}

var StationHashFunction = func(station models.Station) string {
	return strconv.FormatInt(station.ID, 10)
}
