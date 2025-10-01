# Tenjin Implementation - Executive Summary

**Last Updated**: September 30, 2025
**Status**: Core Features Complete ✅

---

## What is Tenjin?

Tenjin is the "brain" of the Metro simulation - a layered service architecture that observes, analyzes, and scores the entire transit system in real-time.

---

## Completed Features

### ✅ Phase 1: Foundation (Observation & Analysis)

**Event System** → Trains emit events (arrivals, departures, ticks, errors)
**Observation Layer** → Collects 500-buffered events from all actors
**Analysis Layer** → Calculates metrics (speed, distance, arrivals/departures)
**Logging** → Auto-rotating file logs (5MB limit)
**Config** → Toggle on/off, adjustable tick rate (default: 1 second)

### ✅ Phase 2: Passengers

**Passenger Model** → ID, sentiment (0-100), state machine (waiting→boarding→riding→disembarking→arrived)
**Station Queues** → Thread-safe passenger waiting lists
**Train Capacity** → 50 passengers/train, boarding/disembarking logic
**Spawning** → Initial + random (5s intervals), reachable destinations only
**Sentiment System** → Time-based decay (-2pts/5sec waiting, -0.5pts/15sec riding)
**Metrics Tracking** → Total/waiting/riding/arrived, boardings, disembarkments, average sentiment

### ✅ Phase 2b: Database Persistence

**Passenger Table** → Active passengers synced every 2s (reflexTick)
**Event Table** → Historical log of all passenger events
**Cleanup** → Arrived passengers removed (prevents memory leak)
**Batch Sync** → Delete-all-then-insert strategy for simplicity

### ✅ Phase 3: Visualization

**Passenger Dots** → Color-coded by sentiment (green/yellow/red) around stations
**Click Interactions** → Trains show data panel, stations open interior scene
**Scene Management** → Map view ↔ Station interior with back button
**Counts Display** → Show waiting/riding passenger numbers

### ✅ Phase 4: Scoring System

**4 Components** (weighted):

- **Passenger Satisfaction** (40%) - Average sentiment, -5pts per 0-sentiment passenger
- **Service Efficiency** (30%) - Completion rate + throughput
- **System Capacity** (20%) - Train utilization (ideal: 60-80%) + station congestion
- **Reliability** (10%) - Train errors, abandoned passengers, coverage

**Grading**: S (95-100), A (85-94), B (75-84), C (65-74), D (50-64), F (<50)
**Daily Reset**: Automatic at midnight, historical stats preserved
**UI Overlay**: Top-left panel, clickable to expand/collapse components
**Color Coding**: Border matches grade (Gold/Green/Yellow/Orange/Red)

---

## Key Technical Decisions

1. **Event Channel**: `interface{}` events to avoid import cycles
2. **Thread Safety**: RWMutex for all shared state (passengers, metrics, scores)
3. **Non-Blocking**: Channel sends use `select/default` to prevent train goroutine blocking
4. **Pointer Architecture**: Stations/trains use pointers to avoid mutex copying
5. **Daily Reset**: Both scoring and metrics reset at midnight for accurate daily performance
6. **Arrived Cleanup**: Passengers removed from tracking when they arrive (prevents inflation)

---

## Architecture Overview

```
Events → Observation → Analysis → Metrics/Score → Logging/UI
  ↑                                                    ↓
Trains/Passengers ←──────────────────────── Display Overlay
```

**Goroutines**:

- Train goroutines (1 per train) - emit events
- Tenjin main loop (1 second tick) - collect & process events
- Passenger spawner - create new passengers every 5s
- Database sync (2 second tick) - persist state

---

## File Structure

```
/internal/tenjin/
  ├── tenjin.go                 # Main coordinator
  ├── observation/collector.go  # Event collection
  ├── analysis/metrics.go       # Metrics calculation
  ├── analysis/logger.go        # File logging
  └── scoring/
      ├── calculator.go         # Score computation
      └── history.go            # Daily tracking

/internal/models/
  ├── passenger.go              # Passenger with sentiment
  ├── train.go                  # With capacity & events
  └── station.go                # With passenger queues

/data/
  ├── passengers.go             # Spawning logic
  └── dump.go                   # DB persistence

/display/
  └── game.go                   # UI with score overlay
```

---

## Configuration

```go
// In control/config.go
TenjinEnabled:        true
TenjinTickRate:       1 * time.Second
PassengerSpawnRate:   5 * time.Second
PassengersPerStation: 3
DisplayMonitor:       1  // Second monitor
```

---

## Metrics Output Example

```
=== TENJIN METRICS (Updated: 17:04:05) ===

*** SYSTEM SCORE: 87.3 (A) ***
  Satisfaction: 78.5 | Efficiency: 92.1 | Capacity: 95.2 | Reliability: 90.0

Total Trains: 5
Average Speed: 0.11
Total Distance Traveled: 59.70
Total Errors: 0

--- PASSENGERS ---
Total Passengers: 68
Waiting: 28 | Riding: 40 | Arrived: 12
Total Boardings: 44 | Disembarkments: 15
Average Sentiment: 75.3/100

Station Arrivals (9 stations):
  Station 4: 5 arrivals
  Station 9: 2 arrivals
  ...
```

---

## What's Next (Planned)

### Intelligence Layer

- Congestion detection algorithms
- Automated intervention recommendations
- Performance optimization suggestions

### Action Layer

- Direct train control (speed, route changes)
- Dynamic scheduling adjustments
- Actor spawning/removal

### Reporting Layer

- LLM-generated daily newspaper summaries
- Key event highlighting
- Narrative storytelling from metrics

### Enhancements

- Historical score graphs
- Real-time alerts for score drops
- Database persistence for daily stats
- Score trend indicators (↑↓)

---

## Bug Fixes Applied

1. ✅ **Memory Leak**: Arrived passengers now removed from tracking maps
2. ✅ **Score Inflation**: Daily reset ensures today's score = today's performance
3. ✅ **Unreachable Destinations**: Passengers only spawn with reachable targets
4. ✅ **Mutex Copying**: Changed to pointers for Station/Train in Line
5. ✅ **Sentiment Timing**: Rate-limited to prevent frame-by-frame drops
6. ✅ **Event Matching**: Fixed type assertions for passenger events

---

## Performance Notes

- **Score Calculation**: O(n) where n = passengers, ~few milliseconds
- **Event Buffer**: 500 capacity, non-blocking sends
- **Database Sync**: Every 2 seconds, batch operations
- **Tenjin Tick**: 1 second (adjustable via config)
- **Passenger Spawn**: 5 seconds (3 per station initially)

---

**For detailed implementation notes, see `tenjin-implementation.md` (1093 lines)**
