# Tenjin Implementation Progress

**Last Updated**: September 30, 2025
**Current Phase**: Phase 1 Complete - Foundation & Basic Metrics
**Next Phase**: Phase 2 - Passenger Integration

---

## Overview

Tenjin is the "brain" of the Metro simulation - a layered service architecture that observes, analyzes, and will eventually manage all aspects of the transit system.

### Architecture Layers

```
┌─────────────────────────────────────────────────────┐
│                     TENJIN                          │
├─────────────────────────────────────────────────────┤
│  Observation Layer (Eyes & Ears) ✅                 │
│  ├─ Listens to all train events                     │
│  ├─ Monitors station states (future)                │
│  └─ Tracks passenger activities (future)            │
├─────────────────────────────────────────────────────┤
│  Analysis Layer (Understanding) ✅                   │
│  ├─ Calculates metrics (avg wait, satisfaction)     │
│  ├─ Aggregates data over time                       │
│  └─ Scores the system health                        │
├─────────────────────────────────────────────────────┤
│  Intelligence Layer (Thinking) 🚧                   │
│  ├─ Evaluates current state                         │
│  ├─ Applies decision strategies                     │
│  └─ Plans interventions                             │
├─────────────────────────────────────────────────────┤
│  Action Layer (Hands) 🚧                            │
│  ├─ Executes commands on trains                     │
│  ├─ Modifies station states                         │
│  └─ Spawns/removes actors                           │
├─────────────────────────────────────────────────────┤
│  Memory Layer (Memory) 🚧                           │
│  ├─ Current simulation state                        │
│  ├─ Historical time-series data                     │
│  └─ Snapshots for replay/debugging                  │
├─────────────────────────────────────────────────────┤
│  Reporting Layer (Communication) 🚧                 │
│  ├─ Daily summaries                                 │
│  ├─ Event logs                                      │
│  └─ Newspaper generation                            │
└─────────────────────────────────────────────────────┘

Legend: ✅ Complete | 🚧 Planned | ⚠️ In Progress
```

---

## Phase 1: Foundation & Basic Metrics ✅

### What Was Built

#### 1. Event System

**File**: `/internal/tenjin/events.go`

Defined event types for train activities:

- `TrainArrivalEvent` - When train reaches a station
- `TrainDepartureEvent` - When train leaves a station
- `TrainTickEvent` - Periodic state snapshot (every 60 ticks = 1 second)
- `TrainErrorEvent` - When errors occur

**Key Design**: Events use anonymous structs with `models.Vector` types to avoid import cycles.

#### 2. Observation Layer

**File**: `/internal/tenjin/observation/collector.go`

- Receives events from trains via buffered channel (500 capacity)
- Runs in its own goroutine continuously collecting events
- Provides batched collection every second via `Collect()`
- Thread-safe with mutex protection

#### 3. Analysis Layer

**Files**:

- `/internal/tenjin/analysis/metrics.go` - Metrics calculation engine
- `/internal/tenjin/analysis/logger.go` - File logging with 5MB rotation

**Metrics Tracked**:

- Total trains in system
- Arrivals per station (count map)
- Departures per station (count map)
- Average train speed (across all trains)
- Total distance traveled (cumulative)
- Error count

**Features**:

- Thread-safe metrics updates
- Real-time aggregation
- Formatted output for logging/display
- Automatic log file rotation at 5MB

#### 4. Main Tenjin Coordinator

**File**: `/internal/tenjin/tenjin.go`

Orchestrates all layers:

- Initializes observation, analysis, and logging layers
- Runs on its own tick rate (configurable, default: 1 second)
- Collects events → Processes through analysis → Logs to file
- Graceful shutdown with context cancellation

#### 5. Configuration

**File**: `/control/config.go`

Added config values:

```go
TenjinEnabled   bool          // Toggle Tenjin on/off (default: true)
TenjinTickRate  time.Duration // How often Tenjin processes (default: 1 second)
```

#### 6. Train Integration

**File**: `/internal/models/train.go`

Modified trains to emit events:

- Added `eventChannel chan<- interface{}` field
- Added `tickCounter int` for periodic events
- Emit arrival events in `logArrival()`
- Emit departure events in `logDeparture()`
- Emit tick events every 60 ticks in `Tick()`
- Emit error events when issues occur
- Non-blocking event sends (channel full = skip event)

#### 7. Data Loading Integration

**File**: `/data/load.go`

Updated `LoadTrains()` to accept and pass event channel to train constructors.

#### 8. Main Integration

**File**: `/main.go`

- Initialize Tenjin after loading stations/lines
- Pass event channel to all trains via `LoadTrains()`
- Start Tenjin's goroutines
- Graceful cleanup on shutdown

---

## How It Works

### Event Flow

```
1. Train Movement
   ↓
2. Train.logArrival/logDeparture/Tick()
   ↓
3. Event sent to eventChannel
   ↓
4. Observation.Collector buffers event
   ↓
5. Tenjin.run() collects events every second
   ↓
6. Analysis.ProcessEvents() updates metrics
   ↓
7. Metrics logged to file & stdout
```

### Type Assertion Pattern

Events are sent as `interface{}` to avoid import cycles. The analysis layer uses type assertions:

```go
if e, ok := event.(struct {
    Type        string
    Train       string
    StationID   int64
    StationName string
    Time        time.Time
    Position    models.Vector
}); ok && e.Type == "train_arrival" {
    // Process arrival
}
```

**Critical**: Types must match exactly, including `models.Vector` (not `interface{}`).

---

## Files Created

```
/internal/tenjin/
  ├── events.go                    # Event type definitions
  ├── tenjin.go                    # Main coordinator
  ├── observation/
  │   └── collector.go             # Event collection
  └── analysis/
      ├── metrics.go               # Metrics calculation
      └── logger.go                # File logging with rotation

/docs/
  └── tenjin-implementation.md     # This file
```

## Files Modified

```
/control/config.go                 # Added TenjinEnabled, TenjinTickRate
/internal/models/train.go          # Added event emission
/data/load.go                      # Updated LoadTrains signature
/main.go                           # Integrated Tenjin initialization
```

---

## Current Output Example

### Metrics (every 1 second)

```
=== TENJIN METRICS (Updated: 14:29:42) ===
Total Trains: 5
Average Speed: 0.11
Total Distance Traveled: 59.70
Total Errors: 0

Station Arrivals (9 stations):
  Station 4: 5 arrivals
  Station 9: 2 arrivals
  Station 12: 2 arrivals
  ...

Station Departures (10 stations):
  Station 4: 5 departures
  Station 11: 2 departures
  ...
=====================================
```

### Log Files

- Location: `/logs/tenjin/tenjin-metrics-[timestamp].log`
- Rotation: Automatic at 5MB
- Format: Human-readable text with timestamps

---

## Technical Decisions & Patterns

### 1. Why `interface{}` for Events?

- Avoids import cycles between packages
- Allows flexibility in event structure
- Trade-off: Requires careful type assertions

### 2. Why Buffered Channel (500)?

- Prevents blocking train goroutines
- Handles burst traffic (many arrivals at once)
- Non-blocking sends with `select/default`

### 3. Why 1 Second Tick Rate?

- Balances real-time updates vs. CPU usage
- Matches train tick events (every 60 ticks)
- Can be adjusted via config for different needs

### 4. Why Separate Observation/Analysis Layers?

- **Observation**: Fast, non-blocking collection
- **Analysis**: Slower processing, aggregation
- Clean separation of concerns
- Easy to test independently

### 5. Thread Safety

- All shared state protected by mutexes
- Read locks for getters, write locks for updates
- Goroutines communicate via channels

---

## Configuration

### Enable/Disable Tenjin

```go
// In control/config.go
TenjinEnabled: true,  // Set to false to disable
```

### Adjust Tick Rate

```go
// In control/config.go
TenjinTickRate: time.Second,  // Process metrics every second
```

### Enable Console Output

```go
// In control/config.go
StdLogs: true,  // Print metrics to stdout
```

---

## Known Issues & Solutions

### Issue: Type Assertion Failures

**Problem**: Events not matching in `ProcessEvents()`
**Solution**: Ensure event struct types match exactly, including `models.Vector`

### Issue: Channel Full

**Problem**: Events dropped when channel capacity exceeded
**Solution**: Increase buffer size or reduce event frequency

### Issue: High CPU Usage

**Problem**: Tenjin processing too frequently
**Solution**: Increase `TenjinTickRate` (e.g., `2 * time.Second`)

---

## Phase 2: Passenger Integration ⚠️

**Status**: In Progress
**Started**: September 30, 2025

### Implementation Progress

#### ✅ Step 1: Passenger Model Created

**File**: `/internal/models/passenger.go`

Created comprehensive `Passenger` struct with:

**Fields**:

- `ID` - Unique identifier
- `Name` - Passenger name
- `Position` - Current Vector position (for rendering)
- `CurrentStation` - Station where passenger is located
- `DestinationStation` - Where passenger wants to go
- `CurrentTrain` - Train passenger is riding (nil if not on train)
- `Sentiment` - Float64 from 0-100 (satisfaction score)
- `State` - PassengerState enum
- `WaitStartTime` - When they started waiting at station
- `JourneyStartTime` - When their journey began
- `eventChannel` - For emitting events to Tenjin

**States** (PassengerState enum):

- `waiting` - Waiting at station for train
- `boarding` - In process of getting on train
- `riding` - On the train traveling
- `disembarking` - Getting off the train
- `arrived` - Reached destination

**Key Methods**:

- `NewPassenger()` - Constructor with event emission
- `UpdateSentiment()` - Decreases sentiment based on wait/ride time
- `StartWaiting()` - Sets to waiting state
- `BoardTrain()` - Puts passenger on train
- `DisembarkTrain()` - Removes from train, checks if destination reached
- `GetSentimentCategory()` - Returns "Happy", "Satisfied", "Neutral", "Frustrated", or "Angry"

**Event Emissions**:

- `emitSpawnEvent()` - When passenger is created
- `emitWaitEvent()` - When waiting at station
- `emitBoardEvent()` - When boarding train
- `emitDisembarkEvent()` - When leaving train
- `emitArriveEvent()` - When reaching destination
- `emitFrustrationEvent()` - When sentiment drops

**Sentiment Model**:

- Starts at 100 (perfect satisfaction)
- Decreases by ~1 point per 5 seconds of waiting
- Minor decrease for long journeys (>2 minutes)
- Future: Will decrease more in crowded trains

**Design Decisions**:

1. Each passenger tracks their own state independently
2. Events emitted for all major state changes
3. Non-blocking event sends (drops if channel full)
4. Implements `Drawing` interface for future visualization
5. Uses pointers to Station/Train to avoid copies

#### ✅ Step 2: Station Passenger Queues

**File**: `/internal/models/station.go` (modified)

Added passenger management to stations:

**New Fields**:

- `WaitingPassengers []*Passenger` - Slice of passengers at station
- `passengerMutex sync.RWMutex` - Thread safety for concurrent access

**New Methods**:

- `AddPassenger(passenger *Passenger)` - Adds passenger to waiting queue
- `RemovePassenger(passenger *Passenger) bool` - Removes passenger from queue
- `GetWaitingPassengersCount() int` - Returns number waiting
- `GetWaitingPassengers() []*Passenger` - Thread-safe copy of queue
- `GetPassengersForDestination(destinationID int64) []*Passenger` - Filter by destination

**Updated Methods**:

- `Update()` - Now also updates all waiting passengers' sentiment

**Thread Safety**:

- All passenger operations protected by RWMutex
- Read operations use RLock (multiple readers OK)
- Write operations use Lock (exclusive access)
- Prevents data races in concurrent goroutines

#### ✅ Step 3: Train Passenger Capacity

**File**: `/internal/models/train.go` (modified)

Added passenger capacity tracking to trains:

**New Fields**:

- `Capacity int` - Maximum passengers (default: 50)
- `Passengers []*Passenger` - Current passengers on board
- `passengerMutex sync.RWMutex` - Thread safety

**New Methods**:

- `AddPassenger(passenger *Passenger) bool` - Adds passenger, returns false if full
- `RemovePassenger(passenger *Passenger) bool` - Removes passenger from train
- `GetPassengerCount() int` - Returns current passenger count
- `IsFull() bool` - Checks if at capacity
- `IsCrowded() bool` - Checks if over 80% capacity
- `GetCapacityPercentage() float64` - Returns capacity usage percentage
- `GetPassengers() []*Passenger` - Thread-safe copy of passengers
- `GetPassengersForStation(stationID int64) []*Passenger` - Filter by destination

**Important Refactoring**:

- Changed `Current Station` to `Current *Station` (pointer to avoid mutex copying)
- Renamed field `make` to `model` to avoid conflict with Go's built-in `make()`
- Renamed parameter `make` to `trainMake` in `NewTrain()`
- Updated `/data/load.go` to pass station pointers

**Note**: Some linter warnings about copying locks in ranges remain (performance warnings, not errors)

### Goals

1. **Add Passenger Model** (`/internal/models/passenger.go`)

   - ID, Name, Position
   - Current station, Destination station
   - Sentiment/satisfaction score
   - State (waiting, boarding, riding, arrived)

2. **Passenger Events**

   - `PassengerSpawnEvent` - New passenger enters system
   - `PassengerWaitEvent` - Waiting at station
   - `PassengerBoardEvent` - Boards a train
   - `PassengerDisembarkEvent` - Leaves train
   - `PassengerArriveEvent` - Reaches destination
   - `PassengerFrustrationEvent` - Sentiment changes

3. **Station Modifications**

   - Add passenger queue/list to stations
   - Track waiting passengers
   - Board/disembark logic during train stops

4. **Train Modifications**

   - Add passenger capacity
   - Track current passengers onboard
   - Board/disembark during wait at stations

5. **Tenjin Analysis Updates**
   - Track passenger count metrics
   - Calculate average wait times
   - Calculate satisfaction scores
   - Identify congested stations

### Implementation Plan

#### Step 1: Create Passenger Model

- Define `Passenger` struct
- Add spawning logic
- Basic movement/state machine

#### Step 2: Integrate with Stations

- Add passenger queues
- Implement waiting logic
- Generate wait events

#### Step 3: Integrate with Trains

- Add capacity tracking
- Implement boarding/disembarking
- Update train wait logic

#### Step 4: Update Tenjin

- Add passenger event types
- Update metrics for passengers
- Track passenger-specific KPIs

#### Step 5: Testing

- Spawn test passengers
- Verify movement
- Check metrics accuracy

---

## Future Phases (Not Started)

### Phase 3: Intelligence Layer

- Decision strategies
- Congestion detection
- Intervention recommendations

### Phase 4: Action Layer

- Direct train control
- Dynamic scheduling
- Spawning/removing trains

### Phase 5: Memory & Persistence

- Database schema for events/metrics
- Historical data storage
- Snapshot/replay system

### Phase 6: Reporting & Newspaper

- Daily summaries
- LLM integration for narrative generation
- Key events highlighting

---

## Development Notes

### Building & Running

```bash
go build && ./metro
```

### Viewing Metrics

```bash
# Real-time in console (set StdLogs: true)
./metro

# From log files
tail -f logs/tenjin/tenjin-metrics-*.log
```

### Testing Tenjin Separately

```bash
# Disable visual display, observe metrics only
# (Future: Add --headless flag)
```

---

## Questions for Discussion

1. **Passenger Spawning**: Random vs. scheduled? Peak hours?
2. **Passenger AI**: Simple random destinations or route optimization?
3. **Sentiment Model**: What factors affect satisfaction?
4. **Capacity Handling**: What happens when trains are full?
5. **Database**: Should passenger events be persisted?

---

## Phase 2: Passenger Integration (COMPLETED)

### Overview

Added passenger simulation with sentiment tracking, boarding/disembarking, and Tenjin metrics.

### Files Created/Modified

- **NEW**: `/data/passengers.go` - Spawning logic (initial + random)
- **NEW**: `/internal/models/passenger.go` - Passenger model with states & sentiment
- **MOD**: `/internal/models/train.go` - Boarding/disembarking logic, capacity tracking
- **MOD**: `/internal/models/station.go` - Waiting passenger queues
- **MOD**: `/internal/tenjin/analysis/metrics.go` - Passenger metrics tracking
- **MOD**: `/control/config.go` - PassengerSpawnRate, PassengersPerStation
- **MOD**: `/main.go` - Context, spawnTick, passenger spawning goroutine

### Passenger Model

- **States**: Waiting, Boarding, Riding, Disembarking, Arrived
- **Sentiment**: 0-100 score, decreases with wait time/crowding
- **Events**: spawn, wait, board, disembark, arrive, frustration

### Train Integration

- **Capacity**: 50 passengers per train
- **Boarding**: Automatic when train arrives at station
- **Disembarking**: Automatic at destination
- **Thread-safe**: RWMutex for passenger list

### Tenjin Metrics Added

- Total passengers, waiting, riding, arrived
- Boardings/disembarkments counters
- Average sentiment across all passengers
- Real-time tracking via event stream

### Configuration

- `PassengerSpawnRate: 5 * time.Second` - New passengers every 5s
- `PassengersPerStation: 3` - Initial passengers per station

### Status

✅ Spawning with reachable destinations (same line only)
✅ Boarding, disembarking all functional
✅ Tenjin tracking passenger metrics
✅ Passengers now arrive at destinations correctly
✅ Database persistence (completed)

### Bug Fixes

1. **Unreachable destinations**: Passengers spawn with destinations only on lines that serve their origin station
2. **Arrival tracking**: Fixed event field names (StationID→DestinationID) so arrivals are correctly tracked
3. **Sentiment calculation**: Arrived passengers excluded from average (only waiting/riding counted)
4. **Journey timing**: JourneyStartTime now set when boarding (not spawning) for accurate journey duration
5. **Station updates**: Added station.Update() to game loop so passengers emit events

## Phase 3: Visualization & UI (IN PROGRESS)

### Overview

Added interactive visualization and click-based data panels for trains and stations.

### Files Created/Modified

- **MOD**: `/display/game.go` - Added passenger visualization, click detection, data panels
- **MOD**: `/display/ui.go` - Added `DrawDataText()` helper function
- **MOD**: `/internal/models/train.go` - Added `GetSpeed()` method

### Features Implemented

#### 1. Passenger Visualization

- Waiting passengers shown as colored dots around stations (max 10 visible)
- Dots arranged in circle around each station
- Color-coded by sentiment:
  - 🟢 Green (70-100): Happy
  - 🟡 Yellow (40-69): Neutral
  - 🔴 Red (0-39): Frustrated/Angry

#### 2. Passenger Count Displays

- Stations: Shows count of waiting passengers
- Trains: Shows count of passengers on board
- Displayed in light blue text next to objects

#### 3. Click Interactions

- **Click Train**: Shows data panel overlay with:

  - Train name
  - Current speed
  - Passenger count (X/50)
  - Next destination or current station
  - Capacity bar (green/yellow/red based on load)

- **Click Station**: Switches to dedicated station interior scene:
  - Station name (centered at top)
  - Platform background with floor
  - Individual passengers displayed as colored circles (sentiment-based)
  - Passenger names below each circle
  - Waiting passenger count
  - Back button (top-left) to return to map

### Technical Details

- Scene management system (SceneMap vs SceneStation)
- Mouse click detection using `ebiten.IsMouseButtonPressed()`
- Click handled only on frame button is pressed (not held)
- Back button bounds checking for scene switching
- Visual bounds checking for train/station hitboxes
- Passengers arranged in rows (10 per row) on platform

### Status

✅ Passenger visualization complete (map view)
✅ Train click → data panel overlay
✅ Station click → dedicated interior scene
✅ Scene management system (map ↔ station)
✅ Back button navigation
✅ Sentiment timing adjustments
⏳ Scoring system (next)
⏳ Newspaper summary (next)

---

## Phase 4: Sentiment Tuning

### Sentiment Drop Rate Fix (2025-09-30)

**Problem**: Sentiment was dropping too rapidly (every frame, 60x/sec)

**Solution**: Implemented time-based sentiment decay with `lastSentimentDrop` tracking

#### Timing Parameters (Adjusted)

- **Waiting**: -2 points every 5 seconds
- **Riding**: -0.5 points every 15 seconds
- **Crowded Train**: Additional -1.0 penalty per interval
- **Frustration Events**: Only emit when sentiment < 50

#### Implementation

Added `lastSentimentDrop time.Time` field to `Passenger` struct

- Initialized in `NewPassenger()`
- Reset on `StartWaiting()`, `BoardTrain()`, and transfer in `DisembarkTrain()`
- Checked in `UpdateSentiment()` to throttle drops to 5-second intervals

#### Result

Passengers can wait ~4 minutes before reaching 0 sentiment (100 points ÷ 2 points per 5 sec = 250 sec). This creates urgency while giving the system reasonable time to process them through the network.

---

## Phase 2b: Database Persistence (2025-09-30)

### Overview

Added SQLite persistence for passengers and their journey events, synced every 2 seconds via reflexTick.

### Database Schema

#### `passenger` Table

Stores active passengers (non-arrived):

- `id` (TEXT PRIMARY KEY) - Unique passenger identifier
- `name` (TEXT) - Passenger name
- `current_station_id` (INTEGER) - Current station, FK to station
- `destination_station_id` (INTEGER) - Destination station, FK to station
- `current_train_id` (INTEGER NULL) - Train passenger is riding (NULL if waiting)
- `state` (TEXT) - waiting/boarding/riding/disembarking/arrived
- `sentiment` (REAL) - Satisfaction score (0-100)
- `spawn_time` (TIMESTAMP) - When passenger entered system
- `created_at`, `updated_at` (TIMESTAMP)

**Indexes**: current_station_id, destination_station_id, current_train_id, state

#### `passenger_event` Table

Historical log of all passenger events:

- `id` (INTEGER PRIMARY KEY AUTOINCREMENT)
- `passenger_id` (TEXT) - FK to passenger
- `event_type` (TEXT) - spawn/wait/board/disembark/arrive/frustration
- `station_id` (INTEGER NULL) - Related station
- `train_id` (INTEGER NULL) - Related train
- `sentiment` (REAL) - Sentiment at time of event
- `metadata` (TEXT) - Additional event data
- `created_at` (TIMESTAMP)

**Indexes**: passenger_id, event_type, created_at

### Files Created/Modified

- **NEW**: `/data/sql/migrations/20250930210333_add_passenger_tables.sql` - Migration
- **NEW**: `/data/sql/queries/passenger.sql` - CRUD queries
- **NEW**: `/internal/dbstore/passenger.sql.go` - Generated sqlc code
- **NEW**: `/internal/baso/passenger.go` - Database wrapper functions
- **MOD**: `/data/sql/schema.sql` - Updated with passenger tables
- **MOD**: `/data/dump.go` - Added `DumpPassengersData()` function
- **MOD**: `/main.go` - Added passenger DB sync to reflexTick goroutine

### Key Functions

#### `/internal/baso/passenger.go`

- `SavePassenger(passenger)` - Insert new passenger
- `UpdatePassengerState(id, state, sentiment)` - Update state/sentiment
- `UpdatePassengerBoarding(id)` - Mark passenger as boarding
- `UpdatePassengerDisembarking(id, stationID, state)` - Mark passenger as disembarking
- `DeletePassenger(id)` - Remove passenger (on arrival)
- `SyncPassengers(passengers)` - Batch sync all active passengers
- `CreatePassengerEvent(...)` - Log passenger event
- `GetPassengerEvents(id)` - Retrieve passenger's event history

#### `/data/dump.go`

- `DumpPassengersData(stations, trains)` - Collects all active passengers from stations/trains and syncs to DB

### Sync Strategy

**Timing**: Every 2 seconds (reflexTick)

**Method**: Batch sync

1. Delete all existing passenger records
2. Collect active passengers from stations (waiting) and trains (riding)
3. Insert all active passengers
4. Arrived passengers are excluded (deleted on arrival)

**Rationale**: Simple "snapshot" approach - DB always reflects current in-memory state

### Event Logging

Passenger events are logged separately to `passenger_event` table for historical analysis:

- Spawn events when passengers enter system
- Wait events during station waiting
- Board events when entering trains
- Disembark events when leaving trains
- Arrive events when reaching destination
- Frustration events when sentiment drops

### Limitations & Future Enhancements

1. **Train ID**: Currently set to NULL - `Train` struct doesn't expose database ID (only `Name`)
   - **Fix**: Add `ID int64` field to `Train` struct in future refactor
2. **Delete on Arrival**: Arrived passengers are removed from DB
   - **Enhancement**: Could keep for X hours with TTL cleanup
3. **Event Logging**: Events logged in separate table but not yet integrated with Tenjin event flow
   - **Enhancement**: Wire up passenger events to `baso.CreatePassengerEvent()` calls

### Technical Notes

- Used `sql.NullInt64` for nullable foreign keys (current_train_id)
- Transactions ensure atomic batch updates
- Indexes on frequently queried columns for performance
- Schema matches Go `models.Passenger` struct closely

## References

- Project Vision: `/PROJECT.md`
- Architecture Diagram: `/docs/tenjin.txt`
- Cursor Rules: `/.cursorrules`
- Main Entry: `/main.go`

---

**Document Status**: Living document, updated as implementation progresses
