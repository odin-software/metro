package analysis

import (
	"fmt"
	"sync"
	"time"

	"github.com/odin-software/metro/internal/models"
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
}

// MetricsEngine calculates and maintains metrics from events
type MetricsEngine struct {
	current            Metrics
	mu                 sync.RWMutex
	trainSpeeds        map[string]float64 // Track individual train speeds for averaging
	trainDistances     map[string]float64 // Track cumulative distance per train
	passengerStates    map[string]string  // Track passenger states (waiting/riding/arrived)
	passengerSentiment map[string]float64 // Track passenger sentiment
}

// NewMetricsEngine creates a new metrics engine
func NewMetricsEngine(totalTrains int) *MetricsEngine {
	return &MetricsEngine{
		current: Metrics{
			TotalTrains:          totalTrains,
			ArrivalsPerStation:   make(map[int64]int),
			DeparturesPerStation: make(map[int64]int),
			TrainsPerLine:        make(map[string]int),
			LastUpdated:          time.Now(),
		},
		trainSpeeds:        make(map[string]float64),
		trainDistances:     make(map[string]float64),
		passengerStates:    make(map[string]string),
		passengerSentiment: make(map[string]float64),
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
			Train       string
			StationID   int64
			StationName string
			Time        time.Time
			Position    models.Vector
		}); ok && e.Type == "train_arrival" {
			m.current.ArrivalsPerStation[e.StationID]++
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
			m.passengerStates[e.PassengerID] = "arrived"
			m.passengerSentiment[e.PassengerID] = e.Sentiment
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
	output += fmt.Sprintf("Total Trains: %d\n", m.current.TotalTrains)
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
	m.trainSpeeds = make(map[string]float64)
	m.trainDistances = make(map[string]float64)
	m.current.LastUpdated = time.Now()
}
