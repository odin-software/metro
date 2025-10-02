# Tenjin Implementation - Executive Summary

**Last Updated**: October 1, 2025
**Status**: Core Features Complete ✅ | Time Schedules Added ⏰

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

### Real City Data Integration

- Santo Domingo metro system data
- Additional cities (NYC, Tokyo, etc.)
- City selector at startup
- Coordinate scaling for different city sizes

### Time Schedules & Punctuality

- Scheduled arrival times per station
- Actual vs. scheduled tracking
- Punctuality metrics in Tenjin
- Intelligence to optimize for on-time performance

### Enhancements

- Historical score graphs
- Real-time alerts for score drops
- Database persistence for daily stats
- Score trend indicators (↑↓)

---

## Phase 5: Newspaper System ✅

### Backend

**Ollama Integration**: Local LLM (llama3.2:1b) for story generation
**Story Types**: Performance, Sentiment, Records, Incidents
**Generation**: Automatic daily at midnight + on-demand
**Caching**: Edition persists until next day
**Tone**: Playful journalism style

### UI

**Newspaper Scene**: Full-screen view with cream/beige background
**Access**: "NEWS" button in top-right corner of map view
**Layout**:

- Dark title bar with "METRO DAILY NEWS"
- Date display
- Story headlines (larger, highlighted background)
- Article text (wrapped at ~90 characters)
- Dividers between stories
- Back button to return to map

**States**:

- Generating: Shows "Generating today's edition..."
- Not Available: Shows "No edition available yet."
- Ready: Displays all stories with formatting

### Performance

- **First Generation**: ~20 seconds (model loads into RAM)
- **Subsequent**: <1 second (model cached for 5 minutes)
- **Model Size**: 1.2GB on disk
- **RAM Usage**: ~1-2GB during generation
- **Raspberry Pi 5**: Compatible (ARM64 support, 8GB RAM sufficient)

---

## Phase 6: Real-World Metrics & Time Schedules ✅

### Real-World Scaling

**Scale Configuration**: 1 pixel = 100 meters (0.01 pixels/meter)
**Map Coverage**: ~70km × 50km (realistic city size)
**Train Speeds**: Fixed from 3000+ km/h to realistic 60-70 km/h
**Conversion Utilities**: `internal/models/metrics.go`

**Functions**:

- `PixelsToMeters`, `MetersToPixels`
- `PixelSpeedToKmPerHour`, `KmPerHourToPixelSpeed`
- `FormatDistance` (e.g., "1.2 km", "500 m")
- `FormatSpeed` (e.g., "60 km/h")
- `CalculateJourneyTime` (distance + speed → time)

**UI Updates**:

- Train data panels now show speed in km/h
- Distance calculations use real-world units

### Simulation Clock

**Time Tracking**: Seconds since midnight (e.g., 28800 = 8:00 AM)
**Start Time**: Configurable (default: 8:00 AM)
**Display**: Top-center of screen (HH:MM:SS format)
**Updates**: Every game tick (~60 Hz)
**Speed Multiplier**: `SimulationSpeed` config (default: 1.0 = real-time)

**Package**: `internal/clock/clock.go`

- `SimulationClock` struct (thread-safe)
- `GetCurrentTime()` - formatted string
- `GetCurrentTimeOfDay()` - seconds since midnight
- Wraps at midnight (24-hour cycle)

### Time Schedules

**Database Schema**: `schedule` table

- `train_id` - which train
- `station_id` - which station
- `scheduled_time` - seconds since midnight (e.g., 28800 = 8:00 AM)
- `sequence_order` - stop number (1, 2, 3...)
- Indexes for fast lookups by train, station, and time

**Schedule Generation**: `data/generate_schedules.go`

- Loads train routes from database
- Calculates realistic travel times:
  - Distance between stations (using coordinates)
  - Train speed and acceleration (physics-based)
  - Dwell time at each station (45 seconds)
- Generates full-day schedules (8:00 AM - 10:00 PM)
- Multiple loops per train (continuous service)

**Current Data**: 191 schedule entries across 5 trains

- Train Cha (Line 1): 28 stops
- Train Che (Line 4): 51 stops
- Train Chi (Line 1): 24 stops
- Train Cho (Line 3): 48 stops
- Train Chu (Line 2): 40 stops

**Example Schedule** (Train Cha):

```
08:00:00  Station 1
08:22:28  Station 2
08:48:32  Station 4
09:06:39  Station 12
10:09:33  Station 1 (loop complete)
...
22:03:57  Station 2 (last stop)
```

**Database Helpers**: `internal/baso/schedule.go`

- `GetScheduleForTrain()` - all stops for a train
- `GetScheduleForStation()` - all arrivals at a station
- `GetNextScheduledStop()` - next stop in sequence
- `GetScheduleByTrainAndStation()` - specific entry

### Punctuality Tracking ✅ (Steps 4 & 5)

**Arrival Tracking**:

- Trains emit `SimTime` (seconds since midnight) with each arrival event
- Train struct includes `ID` and `ClockInterface` for timing
- Events include `TrainID`, `StationID`, and `SimTime`

**Punctuality Calculation**:

- Tenjin looks up scheduled time from database for each arrival
- Calculates delay: `actual_time - scheduled_time`
- Categorizes: **Early** (>2min early), **On-Time** (±2min), **Late** (>2min late)
- Tracks average delay across all arrivals

**Metrics**:

- `TotalArrivalsChecked` - Total arrivals compared against schedule
- `OnTimeCount` & `OnTimePercentage` - Punctuality rate
- `EarlyArrivals` / `LateArrivals` - Breakdown by category
- `AverageDelay` - Mean delay in seconds (negative = early)

**Implementation**:

- `ScheduleDB` interface for database abstraction
- `BasoScheduleAdapter` connects analysis layer to schedule queries
- `trackPunctuality()` method processes each arrival
- Daily reset included in midnight rollover
- Non-blocking: if schedule not found, skip gracefully

**Output Example**:

```
--- PUNCTUALITY ---
Total Arrivals Checked: 42
On-Time: 35 (83.3%) | Early: 5 | Late: 2
Average Delay: 12 seconds (late)
```

**Files Modified**:

- `internal/models/train.go` - Added ID, ClockInterface, updated arrival event
- `data/load.go` - Pass train ID and clock to NewTrain
- `main.go` - Initialize clock before trains, pass to LoadTrains
- `internal/tenjin/analysis/metrics.go` - Added punctuality tracking
- `internal/tenjin/analysis/schedule_adapter.go` - New adapter for DB access
- `internal/tenjin/tenjin.go` - Initialize schedule adapter
- `tools/generate_schedules.go` - Moved from data/ (package conflict)

### UI & Newspaper Integration ✅ (Step 6)

**Train Data Panel Enhancements**:

- Shows scheduled arrival time for next stop
- Displays ETA in minutes (e.g., "Sched: 08:45 (ETA: 5m)")
- System-wide on-time percentage with color coding:
  - Green: ≥85% on-time
  - Yellow: 70-84% on-time
  - Red: <70% on-time
- Increased panel height to accommodate new info

**Newspaper Integration**:

- New story type: `StoryTypePunctuality`
- Automatically generated when ≥10 arrivals tracked
- Includes:
  - On-time percentage
  - Breakdown (on-time/early/late counts)
  - Average delay (formatted as "X seconds early/late")
  - Status assessment (excellent/good/fair/needs improvement)
- Playful, informative tone matching other stories

**Implementation Details**:

- `ScheduleDB` interface for UI layer
- `BasoScheduleAdapter` connects display to database
- `formatTime()` helper converts seconds-since-midnight to HH:MM
- Real-time schedule lookups on train selection
- Clock interface extended for `GetCurrentTimeOfDay()`

**Files Modified**:

- `display/game.go` - Enhanced train panel, added formatTime()
- `display/schedule_adapter.go` - **New**: Database adapter for UI
- `main.go` - Pass schedule adapter to NewGame
- `internal/newspaper/stories.go` - Added punctuality data collection and story type
- `internal/newspaper/generator.go` - Added punctuality story prompt

---

### Time Schedules & Punctuality: Complete! 🎉

All 6 steps completed:

1. ✅ Real-world metrics & clock
2. ✅ Schedule database & generation
3. ✅ Schedule seed data (191 entries)
4. ✅ Actual arrival tracking
5. ✅ Punctuality metrics in Tenjin
6. ✅ UI display & newspaper stories

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
