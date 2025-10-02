# Metro Simulation

A real-time metro train simulation using Go, Ebiten graphics engine, and SQLite for data persistence. Features realistic train physics, passenger management, and an AI "brain" (Tenjin) that monitors and analyzes system performance.

## Quick Start

### Prerequisites

- Go 1.21 or higher
- SQLite3
- Make

### Setup & Run

**Test City (Quick Testing):**

```bash
make seed_test_city
go build && ./metro
```

**Santo Domingo (Real-World Data):**

```bash
make setup_santo_domingo  # Full setup with 69 trains
go build && ./metro
```

## Features

- **Real-time train simulation** with physics-based movement
- **Passenger system** with destinations, boarding, and sentiment
- **Schedule-based operation** (8:00 AM to 10:00 PM)
- **Multiple cities** (Test City, Santo Domingo from OpenStreetMap)
- **Camera controls** for zooming and panning the map
- **AI monitoring** (Tenjin) that tracks performance metrics
- **Daily newspaper** with auto-generated stories about system performance

## City Configurations

### Test City

- 12 stations across 4 lines
- 5 trains with pre-configured schedules
- Ideal for development and testing

### Santo Domingo

- 34 real metro stations (from OpenStreetMap)
- 2 lines (Línea 1: 16 stations, Línea 2: 18 stations)
- 69 trains (40 on L1, 29 on L2) matching real-world fleet
- Auto-generated schedules based on actual distances

## Controls

### Camera

- **Zoom:** Mouse wheel or `+`/`-` keys
- **Pan:** Arrow keys or `WASD`
- **Reset:** `R` key

### UI

- **Click station:** View station details
- **Click train:** View train information
- **Click score panel:** Toggle detailed metrics
- **Click newspaper button:** View daily report

## Available Commands

### Quick Setup

```bash
make seed_test_city          # Setup test city (5 trains)
make setup_santo_domingo     # Complete Santo Domingo setup (69 trains)
```

### Manual Santo Domingo Configuration

```bash
make import_osm              # Import fresh data from OpenStreetMap
make generate_santo_domingo_trains # Generate 69 trains
make seed_santo_domingo      # Load into database
make generate_schedules      # Generate schedules
```

### Switching Between Cities

```bash
# Switch to test city
make seed_test_city
go build && ./metro

# Switch to Santo Domingo
make seed_santo_domingo
go build && ./metro
```

### Database

```bash
make run_migrations          # Run database migrations
make clean_city_data         # Clear city data (keeps schema)
make generate_sqlc           # Regenerate Go types from SQL
```

## Project Structure

```
/control/          - Configuration, logging, constants
/data/             - Database operations, loading/dumping
  /sql/migrations/ - Database schema migrations
  /sql/queries/    - SQL queries (sqlc-generated Go code)
  /sql/seeds/      - Initial data for cities
/display/          - Ebiten game engine & UI rendering
/internal/
  /baso/           - Database abstraction layer
  /broadcast/      - Event messaging system
  /clock/          - Simulation time management
  /dbstore/        - Generated sqlc code
  /models/         - Domain models (Station, Train, Line, etc.)
  /newspaper/      - Auto-generated daily reports
  /tenjin/         - AI brain for system monitoring
  /assets/         - Embedded fonts & images
/tools/            - Utility scripts (OSM import, schedule generation)
```

## How It Works

### Train Movement

- Trains follow scheduled routes with realistic physics
- Acceleration, top speed, and deceleration calculated from train specs
- 45-second dwell time at each station

### Passengers

- Spawn at random stations with destinations
- Board trains heading toward their destination
- Sentiment degrades while waiting or if delayed
- Disembark upon reaching destination

### Schedules

- Auto-generated based on:
  - Real distances between stations
  - Train acceleration and top speed
  - Operating hours (8 AM - 10 PM)
- Each train runs continuous loops

### Tenjin (AI Brain)

- Monitors all train movements and passenger activity
- Calculates system-wide metrics:
  - **Satisfaction:** Based on passenger sentiment
  - **Efficiency:** Distance traveled vs. ideal routes
  - **Capacity:** Train utilization
  - **Reliability:** Schedule adherence (punctuality)
- Generates daily newspaper reports

## Database Schema

Key tables:

- `station` - Metro stations with coordinates
- `line` - Metro lines
- `station_line` - Station-to-line associations with ordering
- `edge` - Connections between stations
- `train` - Train fleet with current positions
- `make` - Train types/models
- `schedule` - Train timetables
- `passenger` - Active passengers in the system

## Configuration

Edit `control/config.go` for settings:

- `PixelsPerMeter` - Map scale (default: 0.01 = 1px = 100m)
- `DisplayScreenWidth/Height` - Window dimensions
- `TerminalMapEnabled` - Console ASCII map (for debugging)
- `StdLogs` - Enable console logging

## Development

### Creating New Migrations

```bash
make create_migration name="add_new_feature"
```

### Creating New Seeds

```bash
make create_seed name="new_city_data"
```

### Adding New Queries

1. Add SQL to `data/sql/queries/`
2. Run `make generate_sqlc`
3. Use generated functions in `internal/baso/`

## OpenStreetMap Integration

Santo Domingo data is imported from OpenStreetMap:

- Queries Overpass API for metro stations
- Extracts real coordinates (lat/lon)
- Converts to pixel positions with accurate scaling
- Applies manually curated line ordering (OSM relations can be incorrect)

See `tools/import_osm.go` and `docs/OSM-IMPORT.md` for details.

## Performance

- 60 FPS target
- Handles 69 trains simultaneously
- ~30,000 schedule entries for full-day operation
- Efficient spatial calculations with vectors
- Concurrent train updates via goroutines

## Documentation

- [`docs/PROJECT.md`](docs/PROJECT.md) - Project overview & roadmap
- [`docs/TENJIN-SUMMARY.md`](docs/TENJIN-SUMMARY.md) - AI brain architecture
- [`docs/CITY-SELECTION.md`](docs/CITY-SELECTION.md) - City data management
- [`docs/CITY-SWITCHING.md`](docs/CITY-SWITCHING.md) - Switching between cities
- [`docs/OSM-IMPORT.md`](docs/OSM-IMPORT.md) - OpenStreetMap integration
- [`docs/pointer-architecture-review.md`](docs/pointer-architecture-review.md) - Technical decisions

## License

[Add your license here]

## Credits

- Train simulation based on real-world Santo Domingo Metro system
- Data sourced from OpenStreetMap contributors
- Built with [Ebiten](https://ebiten.org/) game engine
