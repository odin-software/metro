package newspaper

import (
	"fmt"
	"time"

	"github.com/odin-software/metro/internal/tenjin/analysis"
)

// StoryType represents different categories of news stories
type StoryType string

const (
	StoryTypePerformance StoryType = "performance"
	StoryTypeRecord      StoryType = "record"
	StoryTypeIncident    StoryType = "incident"
	StoryTypeSentiment   StoryType = "sentiment"
)

// Story represents a generated newspaper article
type Story struct {
	Type      StoryType
	Headline  string
	Article   string
	Timestamp time.Time
}

// StoryData contains information needed to generate stories
type StoryData struct {
	Performance map[string]interface{}
	Records     []map[string]interface{}
	Incidents   []map[string]interface{}
	Sentiment   map[string]interface{}
}

// CollectStoryData gathers interesting data from Tenjin metrics
func CollectStoryData(metrics *analysis.Metrics) StoryData {
	data := StoryData{
		Performance: make(map[string]interface{}),
		Records:     []map[string]interface{}{},
		Incidents:   []map[string]interface{}{},
		Sentiment:   make(map[string]interface{}),
	}

	// Performance story data (always generated)
	data.Performance["score"] = fmt.Sprintf("%.1f", metrics.Score.Overall)
	data.Performance["grade"] = metrics.Score.Grade
	data.Performance["satisfaction"] = metrics.Score.PassengerSatisfaction
	data.Performance["efficiency"] = metrics.Score.ServiceEfficiency
	data.Performance["capacity"] = metrics.Score.SystemCapacity
	data.Performance["reliability"] = metrics.Score.Reliability

	// Sentiment story data (if passengers exist)
	if metrics.TotalPassengers > 0 {
		data.Sentiment["sentiment"] = metrics.AverageSentiment
		data.Sentiment["total"] = metrics.TotalPassengers
		data.Sentiment["waiting"] = metrics.PassengersWaiting
		data.Sentiment["riding"] = metrics.PassengersRiding

		// Find busiest station
		maxWaiting := 0
		busiestStation := ""
		for stationID, count := range metrics.ArrivalsPerStation {
			countInt := int(count)
			if countInt > maxWaiting {
				maxWaiting = countInt
				busiestStation = fmt.Sprintf("Station %d", stationID)
			}
		}

		if busiestStation != "" {
			data.Sentiment["station"] = busiestStation
			data.Sentiment["waiting"] = maxWaiting
		} else {
			data.Sentiment["station"] = "N/A"
			data.Sentiment["waiting"] = 0
		}

		// Determine trend (simple heuristic)
		if metrics.AverageSentiment >= 70 {
			data.Sentiment["trend"] = "stable and positive"
		} else if metrics.AverageSentiment >= 50 {
			data.Sentiment["trend"] = "holding steady"
		} else {
			data.Sentiment["trend"] = "declining"
		}
	}

	// Record detection (busiest station)
	if len(metrics.ArrivalsPerStation) > 0 {
		maxArrivals := int64(0)
		busiestStationID := int64(0)

		for stationID, count := range metrics.ArrivalsPerStation {
			if int64(count) > maxArrivals {
				maxArrivals = int64(count)
				busiestStationID = stationID
			}
		}

		if maxArrivals > 20 { // Threshold for "record"
			data.Records = append(data.Records, map[string]interface{}{
				"event": fmt.Sprintf("Station %d handles record traffic", busiestStationID),
				"value": fmt.Sprintf("%d arrivals today", maxArrivals),
			})
		}
	}

	// Record detection (high passenger throughput)
	if metrics.PassengersArrived > 50 {
		data.Records = append(data.Records, map[string]interface{}{
			"event": "System achieves high passenger throughput",
			"value": fmt.Sprintf("%d passengers reached their destinations", metrics.PassengersArrived),
		})
	}

	// Incident detection (low sentiment)
	if metrics.AverageSentiment < 50 && metrics.TotalPassengers > 10 {
		data.Incidents = append(data.Incidents, map[string]interface{}{
			"incident": "Passenger frustration reaches concerning levels",
			"impact":   metrics.TotalPassengers,
		})
	}

	// Incident detection (errors)
	if metrics.ErrorCount > 0 {
		data.Incidents = append(data.Incidents, map[string]interface{}{
			"incident": "System experiences technical difficulties",
			"impact":   fmt.Sprintf("%d errors reported", metrics.ErrorCount),
		})
	}

	return data
}

// SelectStoriesToGenerate chooses which stories to generate based on available data
func SelectStoriesToGenerate(data StoryData) []StoryType {
	stories := []StoryType{}

	// Always include performance story
	stories = append(stories, StoryTypePerformance)

	// Add sentiment story if passengers exist
	if val, ok := data.Sentiment["total"]; ok && val.(int) > 0 {
		stories = append(stories, StoryTypeSentiment)
	}

	// Add one record story if any exist
	if len(data.Records) > 0 {
		stories = append(stories, StoryTypeRecord)
	}

	// Add one incident story if any exist
	if len(data.Incidents) > 0 {
		stories = append(stories, StoryTypeIncident)
	}

	return stories
}
