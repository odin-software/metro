package scoring

import (
	"sync"
	"time"
)

// ScoreHistory tracks scores over time with daily reset
type ScoreHistory struct {
	CurrentDay       time.Time       // Start of current day
	CurrentScore     ScoreComponents // Latest score
	ScoreHistory     []ScoreSnapshot // Historical snapshots
	DailyStats       DailyStatistics // Stats for current day
	mu               sync.RWMutex    // Thread safety
	maxHistoryLength int             // Max snapshots to keep
}

// ScoreSnapshot represents a score at a specific time
type ScoreSnapshot struct {
	Timestamp  time.Time
	Components ScoreComponents
}

// DailyStatistics tracks aggregate stats for the day
type DailyStatistics struct {
	DayStart           time.Time
	MinScore           float64
	MaxScore           float64
	AverageScore       float64
	TotalUpdates       int
	TimeAtGradeS       time.Duration
	TimeAtGradeA       time.Duration
	TimeAtGradeB       time.Duration
	TimeAtGradeC       time.Duration
	TimeAtGradeD       time.Duration
	TimeAtGradeF       time.Duration
	LastGrade          string
	LastGradeStartTime time.Time
}

// NewScoreHistory creates a new score history tracker
func NewScoreHistory() *ScoreHistory {
	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return &ScoreHistory{
		CurrentDay:       dayStart,
		CurrentScore:     ScoreComponents{Overall: 100.0, Grade: "S"},
		ScoreHistory:     make([]ScoreSnapshot, 0),
		maxHistoryLength: 3600, // Keep 1 hour of history at 1 second intervals
		DailyStats: DailyStatistics{
			DayStart:           dayStart,
			MinScore:           100.0,
			MaxScore:           100.0,
			AverageScore:       100.0,
			LastGrade:          "S",
			LastGradeStartTime: now,
		},
	}
}

// Update records a new score, checking for daily reset
func (sh *ScoreHistory) Update(components ScoreComponents) {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Check if we need to reset for a new day
	if dayStart.After(sh.CurrentDay) {
		sh.resetForNewDay(dayStart)
	}

	// Update current score
	sh.CurrentScore = components

	// Add to history
	snapshot := ScoreSnapshot{
		Timestamp:  now,
		Components: components,
	}
	sh.ScoreHistory = append(sh.ScoreHistory, snapshot)

	// Trim history if too long
	if len(sh.ScoreHistory) > sh.maxHistoryLength {
		sh.ScoreHistory = sh.ScoreHistory[len(sh.ScoreHistory)-sh.maxHistoryLength:]
	}

	// Update daily statistics
	sh.updateDailyStats(components, now)
}

// resetForNewDay resets tracking for a new day
func (sh *ScoreHistory) resetForNewDay(dayStart time.Time) {
	// Archive previous day's stats if needed (future: save to DB)

	// Reset for new day
	sh.CurrentDay = dayStart
	sh.ScoreHistory = make([]ScoreSnapshot, 0)
	sh.DailyStats = DailyStatistics{
		DayStart:           dayStart,
		MinScore:           100.0,
		MaxScore:           100.0,
		AverageScore:       100.0,
		LastGrade:          "S",
		LastGradeStartTime: dayStart,
	}
}

// updateDailyStats updates the daily statistics with the new score
func (sh *ScoreHistory) updateDailyStats(components ScoreComponents, now time.Time) {
	stats := &sh.DailyStats

	// Update min/max
	if components.Overall < stats.MinScore {
		stats.MinScore = components.Overall
	}
	if components.Overall > stats.MaxScore {
		stats.MaxScore = components.Overall
	}

	// Update average (rolling)
	stats.TotalUpdates++
	stats.AverageScore = ((stats.AverageScore * float64(stats.TotalUpdates-1)) + components.Overall) / float64(stats.TotalUpdates)

	// Track time at each grade
	if components.Grade != stats.LastGrade {
		// Grade changed, record time spent at previous grade
		timeDiff := now.Sub(stats.LastGradeStartTime)

		switch stats.LastGrade {
		case "S":
			stats.TimeAtGradeS += timeDiff
		case "A":
			stats.TimeAtGradeA += timeDiff
		case "B":
			stats.TimeAtGradeB += timeDiff
		case "C":
			stats.TimeAtGradeC += timeDiff
		case "D":
			stats.TimeAtGradeD += timeDiff
		case "F":
			stats.TimeAtGradeF += timeDiff
		}

		stats.LastGrade = components.Grade
		stats.LastGradeStartTime = now
	}
}

// GetCurrentScore returns the current score (thread-safe)
func (sh *ScoreHistory) GetCurrentScore() ScoreComponents {
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	return sh.CurrentScore
}

// GetDailyStats returns the current day's statistics (thread-safe)
func (sh *ScoreHistory) GetDailyStats() DailyStatistics {
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	return sh.DailyStats
}

// GetRecentHistory returns the last N score snapshots
func (sh *ScoreHistory) GetRecentHistory(count int) []ScoreSnapshot {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	if count > len(sh.ScoreHistory) {
		count = len(sh.ScoreHistory)
	}

	if count == 0 {
		return []ScoreSnapshot{}
	}

	start := len(sh.ScoreHistory) - count
	result := make([]ScoreSnapshot, count)
	copy(result, sh.ScoreHistory[start:])
	return result
}

// GetHistoryInRange returns score snapshots within a time range
func (sh *ScoreHistory) GetHistoryInRange(start, end time.Time) []ScoreSnapshot {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	result := make([]ScoreSnapshot, 0)
	for _, snapshot := range sh.ScoreHistory {
		if snapshot.Timestamp.After(start) && snapshot.Timestamp.Before(end) {
			result = append(result, snapshot)
		}
	}
	return result
}

// GetFormattedStats returns a human-readable summary of daily stats
func (sh *ScoreHistory) GetFormattedStats() string {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	stats := sh.DailyStats

	// Calculate current time at grade (if still in same grade)
	now := time.Now()
	currentGradeTime := now.Sub(stats.LastGradeStartTime)

	totalTime := stats.TimeAtGradeS + stats.TimeAtGradeA + stats.TimeAtGradeB +
		stats.TimeAtGradeC + stats.TimeAtGradeD + stats.TimeAtGradeF + currentGradeTime

	percentS := 0.0
	percentA := 0.0
	percentB := 0.0
	percentC := 0.0
	percentD := 0.0
	percentF := 0.0

	if totalTime > 0 {
		percentS = (stats.TimeAtGradeS.Seconds() / totalTime.Seconds()) * 100
		percentA = (stats.TimeAtGradeA.Seconds() / totalTime.Seconds()) * 100
		percentB = (stats.TimeAtGradeB.Seconds() / totalTime.Seconds()) * 100
		percentC = (stats.TimeAtGradeC.Seconds() / totalTime.Seconds()) * 100
		percentD = (stats.TimeAtGradeD.Seconds() / totalTime.Seconds()) * 100
		percentF = (stats.TimeAtGradeF.Seconds() / totalTime.Seconds()) * 100

		// Add current grade time to its percentage
		currentPercent := (currentGradeTime.Seconds() / totalTime.Seconds()) * 100
		switch stats.LastGrade {
		case "S":
			percentS += currentPercent
		case "A":
			percentA += currentPercent
		case "B":
			percentB += currentPercent
		case "C":
			percentC += currentPercent
		case "D":
			percentD += currentPercent
		case "F":
			percentF += currentPercent
		}
	}

	return "--- DAILY SCORE STATISTICS ---\n" +
		"Current Score: " + formatFloat(sh.CurrentScore.Overall) + " (" + sh.CurrentScore.Grade + ")\n" +
		"Day Average: " + formatFloat(stats.AverageScore) + "\n" +
		"Min/Max: " + formatFloat(stats.MinScore) + " / " + formatFloat(stats.MaxScore) + "\n" +
		"Time at Grade: S=" + formatPercent(percentS) + " | A=" + formatPercent(percentA) +
		" | B=" + formatPercent(percentB) + " | C=" + formatPercent(percentC) +
		" | D=" + formatPercent(percentD) + " | F=" + formatPercent(percentF) + "\n"
}

// Helper function to format floats
func formatFloat(f float64) string {
	return time.Duration(f * float64(time.Second)).String()[:4] // Trim to 2 decimal places hack
}

// Helper function to format percentages
func formatPercent(f float64) string {
	return time.Duration(f * float64(time.Second)).String()[:4] + "%"
}
