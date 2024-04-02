package control

import (
	"strconv"
	"time"

	"github.com/odin-software/metro/internal/models"
)

type Config struct {
	LoopDuration        time.Duration
	LoopDurationOffset  time.Duration
	PortCity            int
	PortReporter        int
	ReflexDuration      time.Duration
	TerminalMapDuration time.Duration
	TerminalMapEnabled  bool
}

var DefaultConfig = Config{
	LoopDuration:        20 * time.Millisecond,
	LoopDurationOffset:  -1 * time.Millisecond,
	PortCity:            2221,
	PortReporter:        2222,
	ReflexDuration:      3000 * time.Millisecond,
	TerminalMapDuration: 800 * time.Millisecond,
	TerminalMapEnabled:  false,
}

var StationHashFunction = func(station models.Station) string {
	return strconv.FormatInt(station.ID, 10)
}
