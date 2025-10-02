package scoring

import (
	"math"
)

// ScoreComponents holds the breakdown of the overall score
type ScoreComponents struct {
	PassengerSatisfaction float64 // 0-100, weighted 40%
	ServiceEfficiency     float64 // 0-100, weighted 30%
	SystemCapacity        float64 // 0-100, weighted 20%
	Reliability           float64 // 0-100, weighted 10%
	Overall               float64 // Weighted sum (0-100)
	Grade                 string  // S, A, B, C, D, F
}

// ScoreInputs contains the metrics needed to calculate the score
type ScoreInputs struct {
	// Passenger metrics
	TotalPassengers     int
	WaitingPassengers   int
	RidingPassengers    int
	ArrivedPassengers   int
	AverageSentiment    float64
	ZeroSentimentCount  int // Count of passengers with 0 sentiment
	TotalBoardings      int
	TotalDisembarkments int

	// Journey metrics
	TotalJourneyTime  float64 // Sum of all journey times (seconds)
	CompletedJourneys int     // Number of completed journeys

	// Train metrics
	TotalTrains           int
	TotalTrainCapacity    int
	PassengersOnTrains    int
	AverageTrainOccupancy float64 // Percentage

	// Station metrics
	TotalStations          int
	StationsWithPassengers int
	MaxStationCongestion   int // Max passengers at any station
	AverageStationWait     int // Average passengers per station

	// Reliability metrics
	TrainErrors int
}

// CalculateScore computes the overall score and its components
func CalculateScore(inputs ScoreInputs) ScoreComponents {
	components := ScoreComponents{}

	// 1. Passenger Satisfaction (40%)
	components.PassengerSatisfaction = calculatePassengerSatisfaction(inputs)

	// 2. Service Efficiency (30%)
	components.ServiceEfficiency = calculateServiceEfficiency(inputs)

	// 3. System Capacity (20%)
	components.SystemCapacity = calculateSystemCapacity(inputs)

	// 4. Reliability (10%)
	components.Reliability = calculateReliability(inputs)

	// Weighted Overall Score
	components.Overall = (components.PassengerSatisfaction * 0.40) +
		(components.ServiceEfficiency * 0.30) +
		(components.SystemCapacity * 0.20) +
		(components.Reliability * 0.10)

	// Assign Grade
	components.Grade = assignGrade(components.Overall)

	return components
}

// calculatePassengerSatisfaction computes the passenger satisfaction component
func calculatePassengerSatisfaction(inputs ScoreInputs) float64 {
	if inputs.TotalPassengers == 0 {
		return 100.0 // No passengers = perfect (no complaints)
	}

	// Base score from average sentiment
	baseScore := inputs.AverageSentiment

	// Penalty: -5 points per passenger with 0 sentiment
	zeroSentimentPenalty := float64(inputs.ZeroSentimentCount) * 5.0

	score := baseScore - zeroSentimentPenalty

	// Clamp to 0-100
	return math.Max(0, math.Min(100, score))
}

// calculateServiceEfficiency computes the service efficiency component
func calculateServiceEfficiency(inputs ScoreInputs) float64 {
	if inputs.TotalPassengers == 0 {
		return 100.0 // No passengers = no service needed
	}

	score := 0.0

	// Factor 1: Completion rate (50% of efficiency score)
	// Percentage of passengers who arrived vs total spawned
	completionRate := 0.0
	totalSpawned := inputs.ArrivedPassengers + inputs.WaitingPassengers + inputs.RidingPassengers
	if totalSpawned > 0 {
		completionRate = (float64(inputs.ArrivedPassengers) / float64(totalSpawned)) * 100.0
	}
	score += completionRate * 0.5

	// Factor 2: Passenger throughput (50% of efficiency score)
	// Ratio of boardings to total passengers (want high turnover)
	if inputs.TotalPassengers > 0 {
		throughputRatio := float64(inputs.TotalBoardings) / float64(inputs.TotalPassengers)
		// Normalize: 1.0 throughput = 50 points, 2.0+ = 100 points
		throughputScore := math.Min(100, throughputRatio*50)
		score += throughputScore * 0.5
	}

	return math.Max(0, math.Min(100, score))
}

// calculateSystemCapacity computes the system capacity utilization component
func calculateSystemCapacity(inputs ScoreInputs) float64 {
	if inputs.TotalTrains == 0 {
		return 0.0 // No trains = system failure
	}

	score := 0.0

	// Factor 1: Train capacity utilization (70% of capacity score)
	// Ideal range: 60-80% = 100 points
	// Below 60% or above 80% = penalty
	trainScore := 0.0
	occupancy := inputs.AverageTrainOccupancy

	if occupancy >= 60 && occupancy <= 80 {
		// Ideal range
		trainScore = 100.0
	} else if occupancy < 60 {
		// Underutilized: linear from 0% = 0 points to 60% = 100 points
		trainScore = (occupancy / 60.0) * 100.0
	} else {
		// Overcrowded: linear from 80% = 100 points to 100% = 50 points
		trainScore = 100.0 - ((occupancy - 80.0) / 20.0 * 50.0)
	}
	score += trainScore * 0.7

	// Factor 2: Station congestion (30% of capacity score)
	// Penalize if any station has too many waiting passengers
	stationScore := 100.0
	if inputs.MaxStationCongestion > 20 {
		// Penalty: -2 points per passenger over 20
		penalty := float64(inputs.MaxStationCongestion-20) * 2.0
		stationScore = math.Max(0, 100.0-penalty)
	}
	score += stationScore * 0.3

	return math.Max(0, math.Min(100, score))
}

// calculateReliability computes the reliability component
func calculateReliability(inputs ScoreInputs) float64 {
	score := 100.0

	// Penalty: -10 points per train error
	errorPenalty := float64(inputs.TrainErrors) * 10.0
	score -= errorPenalty

	// Penalty: -1 point per passenger with 0 sentiment (abandoned)
	abandonedPenalty := float64(inputs.ZeroSentimentCount) * 1.0
	score -= abandonedPenalty

	// Bonus: Service coverage (all stations served)
	if inputs.TotalStations > 0 {
		coverageRatio := float64(inputs.StationsWithPassengers) / float64(inputs.TotalStations)
		// Want at least 50% coverage for full points
		if coverageRatio < 0.5 {
			score -= (0.5 - coverageRatio) * 20.0
		}
	}

	return math.Max(0, math.Min(100, score))
}

// assignGrade returns a letter grade based on the overall score
func assignGrade(score float64) string {
	switch {
	case score >= 95:
		return "S"
	case score >= 85:
		return "A"
	case score >= 75:
		return "B"
	case score >= 65:
		return "C"
	case score >= 50:
		return "D"
	default:
		return "F"
	}
}

// GetGradeColor returns a color hex code for a grade (for UI display)
func GetGradeColor(grade string) string {
	switch grade {
	case "S":
		return "#FFD700" // Gold
	case "A":
		return "#00FF00" // Green
	case "B":
		return "#90EE90" // Light Green
	case "C":
		return "#FFFF00" // Yellow
	case "D":
		return "#FFA500" // Orange
	case "F":
		return "#FF0000" // Red
	default:
		return "#FFFFFF" // White
	}
}
