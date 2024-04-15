package control

import (
	"io"
	"log"
	"os"
	"time"
)

// Logs created by day.
type Loporter struct {
	currentDay time.Time
	log        *log.Logger
}

const LogFormat = "2006-01-02"
const LogFileFlags = os.O_RDWR | os.O_CREATE | os.O_APPEND
const LogFilePerms = 0666

var LPT Loporter

func InitLogger() {
	// Checks if the logs directory exists.
	logDirectory := DefaultConfig.LogsDirectory
	err := os.MkdirAll(logDirectory, 0755)
	if err != nil {
		log.Fatalf("Cannot create logs directory")
	}

	logFile := logDirectory + time.Now().Format(LogFormat) + ".log"

	f, err := os.OpenFile(logFile, LogFileFlags, LogFilePerms)
	if err != nil {
		log.Panicf("error opening file: %v", err)
	}

	var wr io.Writer
	if DefaultConfig.StdLogs {
		wr = io.MultiWriter(f, os.Stdout)
	} else {
		wr = io.MultiWriter(f)
	}
	_log := log.New(wr, "INFO: ", log.LstdFlags)

	LPT = Loporter{
		currentDay: time.Now(),
		log:        _log,
	}
}

func dateSameDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func Log(message string) {
	if !dateSameDay(time.Now(), LPT.currentDay) {
		LPT.currentDay = time.Now()

		dailyLogFile := DefaultConfig.LogsDirectory + time.Now().Format(LogFormat) + ".log"
		newFile, err := os.OpenFile(dailyLogFile, LogFileFlags, LogFilePerms)
		if err != nil {
			log.Panicf("error opening file: %v", err)
		}

		var wr io.Writer
		if DefaultConfig.StdLogs {
			wr = io.MultiWriter(newFile, os.Stdout)
		} else {
			wr = io.MultiWriter(newFile)
		}
		LPT.log.SetOutput(wr)
	}

	LPT.log.Println(message)
}
