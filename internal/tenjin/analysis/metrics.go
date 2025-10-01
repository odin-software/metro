package analysis

import (
	"fmt"
	"sync"
	"time"

	"github.com/odin-software/metro/internal/models"
	"github.com/odin-software/metro/internal/tenjin/scoring"
)

// Metrics holds the current state of system metrics
type Metrics struct {
	TotalTrains             int
	ArrivalsPerStation      map[int64]int
	DeparturesPerStation    map[int64]int
	AverageSpeed            float64
	TotalDistanceTraveled   float64
	TrainsPerLine           map[string]int // Future: when we track line info
	ErrorCount              int
	TotalPassengers         int
	PassengersWaiting       int
	PassengersRiding        int
	PassengersArrived       int
	AverageSentiment        float64
	PassengerBoardings      int
	PassengerDisembarkments int
	LastUpdated             time.Time
	Score                   scoring.ScoreComponents // System score
	// Punctuality metrics
	TotalArrivalsChecked int
	OnTimeArrivals       int     // Within ±2 minutes of schedule
	EarlyArrivals        int     // More than 2 minutes early
	LateArrivals         int     // More than 2 minutes late
	AverageDelay         float64 // Average delay in seconds (negative = early)
	OnTimePercentage     float64 // Percentage of on-time arrivals
}

// MetricsEngine calculates and maintains metrics from events
type MetricsEngine struct {
	current                Metrics
	mu                     sync.RWMutex
	trainSpeeds            map[string]float64    // Track individual train speeds for averaging
	trainDistances         map[string]float64    // Track cumulative distance per train
	passengerStates        map[string]string     // Track passenger states (waiting/riding/arrived)
	passengerSentiment     map[string]float64    // Track passenger sentiment
	scoreHistory           *scoring.ScoreHistory // Score tracking
	totalStations          int                   // Total stations in system
	stationsWithPassengers map[int64]bool        // Stations that have served passengers
	currentDay             time.Time             // Track current day for daily resets
	delays                 []float64             // Track all delays for averaging
	scheduleDB             ScheduleDB            // Interface for schedule lookups
}

// ScheduleDB provides schedule lookup functionality
type ScheduleDB interface {
	GetScheduleByTrainAndStation(trainID, stationID int64) (Schedule, error)
}

// Schedule represents a scheduled stop
type Schedule struct {
	TrainID       int64
	StationID     int64
	ScheduledTime int // Seconds since midnight
}

// NewMetricsEngine creates a new metrics engine
func NewMetricsEngine(totalTrains int, scheduleDB ScheduleDB) *MetricsEngine {
	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return &MetricsEngine{
		current: Metrics{
			TotalTrains:          totalTrains,
			ArrivalsPerStation:   make(map[int64]int),
			DeparturesPerStation: make(map[int64]int),
			TrainsPerLine:        make(map[string]int),
			LastUpdated:          now,
			Score:                scoring.ScoreComponents{Overall: 100.0, Grade: "S"},
		},
		trainSpeeds:            make(map[string]float64),
		trainDistances:         make(map[string]float64),
		passengerStates:        make(map[string]string),
		passengerSentiment:     make(map[string]float64),
		scoreHistory:           scoring.NewScoreHistory(),
		stationsWithPassengers: make(map[int64]bool),
		totalStations:          0, // Will be set based on events
		currentDay:             dayStart,
		delays:                 make([]float64, 0),
		scheduleDB:             scheduleDB,
	}
}

// ProcessEvents updates metrics based on a batch of events
// Events are passed as interface{} to avoid import cycles
func (m *MetricsEngine) ProcessEvents(events []interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, event := range events {
		// Use type assertion on the struct itself
		if e, ok := event.(struct {
			Type        string
			TrainID     int64
			Train       string
			StationID   int64
			StationName string
			Time        time.Time
			SimTime     int
			Position    models.Vector
		}); ok && e.Type == "train_arrival" {
			m.current.ArrivalsPerStation[e.StationID]++
			// Track punctuality
			m.trackPunctuality(e.TrainID, e.StationID, e.SimTime)
		} else if e, ok := event.(struct {
			Type        string
			Train       string
			StationID   int64
			StationName string
			NextStation string
			Time        time.Time
			Position    models.Vector
		}); ok && e.Type == "train_departure" {
			m.current.DeparturesPerStation[e.StationID]++
		} else if e, ok := event.(struct {
			Type           string
			Train          string
			Position       models.Vector
			Velocity       models.Vector
			Speed          float64
			CurrentStation int64
			NextStation    int64
			Time           time.Time
		}); ok && e.Type == "train_tick" {
			// Update speed tracking
			m.trainSpeeds[e.Train] = e.Speed
			// Update distance (speed * time since last tick, roughly)
			// Since we tick every second, distance = speed * 1 second
			m.trainDistances[e.Train] += e.Speed
		} else if e, ok := event.(struct {
			Type    string
			Train   string
			Error   string
			Context string
			Time    time.Time
		}); ok && e.Type == "train_error" {
			m.current.ErrorCount++
		} else if e, ok := event.(struct {
			Type            string
			PassengerID     string
			PassengerName   string
			StationID       int64
			StationName     string
			DestinationID   int64
			DestinationName string
			Time            time.Time
		}); ok && e.Type == "passenger_spawn" {
			m.passengerStates[e.PassengerID] = "waiting"
			m.passengerSentiment[e.PassengerID] = 100.0
			m.stationsWithPassengers[e.StationID] = true
		} else if e, ok := event.(struct {
			Type          string
			PassengerID   string
			PassengerName string
			TrainName     string
			StationID     int64
			StationName   string
			Sentiment     float64
			Time          time.Time
		}); ok && e.Type == "passenger_board" {
			m.passengerStates[e.PassengerID] = "riding"
			m.passengerSentiment[e.PassengerID] = e.Sentiment
			m.current.PassengerBoardings++
		} else if e, ok := event.(struct {
			Type          string
			PassengerID   string
			PassengerName string
			StationID     int64
			StationName   string
			Sentiment     float64
			Time          time.Time
		}); ok && e.Type == "passenger_disembark" {
			m.passengerStates[e.PassengerID] = "waiting"
			m.passengerSentiment[e.PassengerID] = e.Sentiment
			m.current.PassengerDisembarkments++
		} else if e, ok := event.(struct {
			Type            string
			PassengerID     string
			PassengerName   string
			DestinationID   int64
			DestinationName string
			JourneyDuration time.Duration
			Sentiment       float64
			Time            time.Time
		}); ok && e.Type == "passenger_arrive" {
			// Passenger has arrived - remove from tracking and increment counter
			delete(m.passengerStates, e.PassengerID)
			delete(m.passengerSentiment, e.PassengerID)
			m.current.PassengersArrived++
		} else if e, ok := event.(struct {
			Type          string
			PassengerID   string
			PassengerName string
			StationID     int64
			StationName   string
			WaitDuration  time.Duration
			Sentiment     float64
			Time          time.Time
		}); ok && e.Type == "passenger_wait" {
			m.passengerSentiment[e.PassengerID] = e.Sentiment
		}
	}

	// Recalculate aggregates
	m.calculateAverages()
	m.current.LastUpdated = time.Now()
}

// calculateAverages recomputes average speed and total distance
// trackPunctuality compares actual arrival time with scheduled time
func (m *MetricsEngine) trackPunctuality(trainID, stationID int64, actualTime int) {
	// Skip if no schedule DB available
	if m.scheduleDB == nil {
		return
	}

	// Look up scheduled time
	schedule, err := m.scheduleDB.GetScheduleByTrainAndStation(trainID, stationID)
	if err != nil {
		// No schedule found for this train/station combo (not an error, just skip)
		return
	}

	// Calculate delay (negative = early, positive = late)
	delay := float64(actualTime - schedule.ScheduledTime)
	m.delays = append(m.delays, delay)

	// Increment counters
	m.current.TotalArrivalsChecked++

	// Categorize: on-time = within ±2 minutes (±120 seconds)
	if delay <= -120 {
		m.current.EarlyArrivals++
	} else if delay >= 120 {
		m.current.LateArrivals++
	} else {
		m.current.OnTimeArrivals++
	}

	// Update average delay
	if len(m.delays) > 0 {
		totalDelay := 0.0
		for _, d := range m.delays {
			totalDelay += d
		}
		m.current.AverageDelay = totalDelay / float64(len(m.delays))
	}

	// Update on-time percentage
	if m.current.TotalArrivalsChecked > 0 {
		m.current.OnTimePercentage = (float64(m.current.OnTimeArrivals) / float64(m.current.TotalArrivalsChecked)) * 100.0
	}
}

func (m *MetricsEngine) calculateAverages() {
	// Average speed across all trains
	if len(m.trainSpeeds) > 0 {
		totalSpeed := 0.0
		for _, speed := range m.trainSpeeds {
			totalSpeed += speed
		}
		m.current.AverageSpeed = totalSpeed / float64(len(m.trainSpeeds))
	}

	// Total distance traveled
	totalDistance := 0.0
	for _, distance := range m.trainDistances {
		totalDistance += distance
	}
	m.current.TotalDistanceTraveled = totalDistance

	// Count passenger states
	m.current.PassengersWaiting = 0
	m.current.PassengersRiding = 0
	m.current.TotalPassengers = len(m.passengerStates)
	for _, state := range m.passengerStates {
		switch state {
		case "waiting":
			m.current.PassengersWaiting++
		case "riding":
			m.current.PassengersRiding++
		}
	}

	// Average sentiment (only for active passengers - waiting or riding)
	activeSentiment := 0.0
	activeCount := 0
	for passengerID, sentiment := range m.passengerSentiment {
		if state := m.passengerStates[passengerID]; state == "waiting" || state == "riding" {
			activeSentiment += sentiment
			activeCount++
		}
	}
	if activeCount > 0 {
		m.current.AverageSentiment = activeSentiment / float64(activeCount)
	} else {
		m.current.AverageSentiment = 0.0
	}

	// Calculate and update system score
	m.calculateScore()
}

// calculateScore computes the overall system score based on current metrics
func (m *MetricsEngine) calculateScore() {
	// Check if we need to reset for a new day
	now := time.Now()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if dayStart.After(m.currentDay) {
		m.currentDay = dayStart
		// Reset cumulative metrics for new day (without locking - already locked by caller)
		m.current.ArrivalsPerStation = make(map[int64]int)
		m.current.DeparturesPerStation = make(map[int64]int)
		m.current.ErrorCount = 0
		m.current.TotalDistanceTraveled = 0
		m.current.PassengersArrived = 0
		m.current.PassengerBoardings = 0
		m.current.PassengerDisembarkments = 0
		m.trainSpeeds = make(map[string]float64)
		m.trainDistances = make(map[string]float64)
		m.stationsWithPassengers = make(map[int64]bool)
		// Reset punctuality metrics
		m.current.TotalArrivalsChecked = 0
		m.current.OnTimeArrivals = 0
		m.current.EarlyArrivals = 0
		m.current.LateArrivals = 0
		m.current.AverageDelay = 0
		m.current.OnTimePercentage = 0
		m.delays = make([]float64, 0)
		// Note: Don't reset passengerStates/passengerSentiment - those track active passengers
	}

	// Count passengers with 0 sentiment
	zeroSentimentCount := 0
	for _, sentiment := range m.passengerSentiment {
		if sentiment <= 0.0 {
			zeroSentimentCount++
		}
	}

	// Calculate train capacity metrics
	totalCapacity := m.current.TotalTrains * 50 // Assuming 50 passenger capacity per train
	averageOccupancy := 0.0
	if totalCapacity > 0 {
		averageOccupancy = (float64(m.current.PassengersRiding) / float64(totalCapacity)) * 100.0
	}

	// Find max station congestion (most waiting at any single station)
	maxCongestion := 0
	for _, count := range m.current.ArrivalsPerStation {
		if count > maxCongestion {
			maxCongestion = count
		}
	}

	// Calculate average station wait (rough estimate based on waiting passengers)
	avgStationWait := 0
	if m.totalStations > 0 {
		avgStationWait = m.current.PassengersWaiting / m.totalStations
	}

	// Set totalStations from arrivals map if not set
	if m.totalStations == 0 {
		m.totalStations = len(m.current.ArrivalsPerStation)
	}

	// Prepare scoring inputs
	inputs := scoring.ScoreInputs{
		TotalPassengers:        m.current.TotalPassengers,
		WaitingPassengers:      m.current.PassengersWaiting,
		RidingPassengers:       m.current.PassengersRiding,
		ArrivedPassengers:      m.current.PassengersArrived,
		AverageSentiment:       m.current.AverageSentiment,
		ZeroSentimentCount:     zeroSentimentCount,
		TotalBoardings:         m.current.PassengerBoardings,
		TotalDisembarkments:    m.current.PassengerDisembarkments,
		TotalTrains:            m.current.TotalTrains,
		TotalTrainCapacity:     totalCapacity,
		PassengersOnTrains:     m.current.PassengersRiding,
		AverageTrainOccupancy:  averageOccupancy,
		TotalStations:          m.totalStations,
		StationsWithPassengers: len(m.stationsWithPassengers),
		MaxStationCongestion:   maxCongestion,
		AverageStationWait:     avgStationWait,
		TrainErrors:            m.current.ErrorCount,
	}

	// Calculate score
	score := scoring.CalculateScore(inputs)

	// Update current score
	m.current.Score = score

	// Record in history
	m.scoreHistory.Update(score)
}

// GetMetrics returns a copy of current metrics (thread-safe)
func (m *MetricsEngine) GetMetrics() Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy to avoid data races
	metrics := m.current

	// Deep copy maps
	metrics.ArrivalsPerStation = make(map[int64]int)
	for k, v := range m.current.ArrivalsPerStation {
		metrics.ArrivalsPerStation[k] = v
	}

	metrics.DeparturesPerStation = make(map[int64]int)
	for k, v := range m.current.DeparturesPerStation {
		metrics.DeparturesPerStation[k] = v
	}

	return metrics
}

// GetFormattedOutput returns a human-readable metrics summary
func (m *MetricsEngine) GetFormattedOutput() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	output := fmt.Sprintf("\n=== TENJIN METRICS (Updated: %s) ===\n",
		m.current.LastUpdated.Format("15:04:05"))

	// System Score (prominent display)
	output += fmt.Sprintf("\n*** SYSTEM SCORE: %.1f (%s) ***\n",
		m.current.Score.Overall, m.current.Score.Grade)
	output += fmt.Sprintf("  Satisfaction: %.1f | Efficiency: %.1f | Capacity: %.1f | Reliability: %.1f\n",
		m.current.Score.PassengerSatisfaction,
		m.current.Score.ServiceEfficiency,
		m.current.Score.SystemCapacity,
		m.current.Score.Reliability)

	output += fmt.Sprintf("\nTotal Trains: %d\n", m.current.TotalTrains)
	output += fmt.Sprintf("Average Speed: %.2f\n", m.current.AverageSpeed)
	output += fmt.Sprintf("Total Distance Traveled: %.2f\n", m.current.TotalDistanceTraveled)
	output += fmt.Sprintf("Total Errors: %d\n", m.current.ErrorCount)

	output += "\n--- PASSENGERS ---\n"
	output += fmt.Sprintf("Total Passengers: %d\n", m.current.TotalPassengers)
	output += fmt.Sprintf("Waiting: %d | Riding: %d | Arrived: %d\n",
		m.current.PassengersWaiting, m.current.PassengersRiding, m.current.PassengersArrived)
	output += fmt.Sprintf("Total Boardings: %d | Disembarkments: %d\n",
		m.current.PassengerBoardings, m.current.PassengerDisembarkments)
	output += fmt.Sprintf("Average Sentiment: %.1f/100\n", m.current.AverageSentiment)

	// Punctuality metrics
	if m.current.TotalArrivalsChecked > 0 {
		output += "\n--- PUNCTUALITY ---\n"
		output += fmt.Sprintf("Total Arrivals Checked: %d\n", m.current.TotalArrivalsChecked)
		output += fmt.Sprintf("On-Time: %d (%.1f%%) | Early: %d | Late: %d\n",
			m.current.OnTimeArrivals, m.current.OnTimePercentage,
			m.current.EarlyArrivals, m.current.LateArrivals)
		// Format average delay nicely
		if m.current.AverageDelay < 0 {
			output += fmt.Sprintf("Average Delay: %.0f seconds (early)\n", -m.current.AverageDelay)
		} else {
			output += fmt.Sprintf("Average Delay: %.0f seconds (late)\n", m.current.AverageDelay)
		}
	}

	output += fmt.Sprintf("\nStation Arrivals (%d stations):\n", len(m.current.ArrivalsPerStation))
	for stationID, count := range m.current.ArrivalsPerStation {
		output += fmt.Sprintf("  Station %d: %d arrivals\n", stationID, count)
	}

	output += fmt.Sprintf("\nStation Departures (%d stations):\n", len(m.current.DeparturesPerStation))
	for stationID, count := range m.current.DeparturesPerStation {
		output += fmt.Sprintf("  Station %d: %d departures\n", stationID, count)
	}

	output += "=====================================\n"

	return output
}

// Reset clears all metrics (useful for testing or daily resets)
func (m *MetricsEngine) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.current.ArrivalsPerStation = make(map[int64]int)
	m.current.DeparturesPerStation = make(map[int64]int)
	m.current.ErrorCount = 0
	m.current.TotalDistanceTraveled = 0
	m.current.PassengersArrived = 0
	m.current.PassengerBoardings = 0
	m.current.PassengerDisembarkments = 0
	m.trainSpeeds = make(map[string]float64)
	m.trainDistances = make(map[string]float64)
	m.stationsWithPassengers = make(map[int64]bool)
	m.current.LastUpdated = time.Now()
	// Note: Don't reset passengerStates/passengerSentiment - those track active passengers
}
