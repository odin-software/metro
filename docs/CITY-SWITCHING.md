# City Switching Guide

This guide explains how to switch between different city configurations in the Metro simulation.

## Quick Start

### Option 1: Test City (Quick Testing)

```bash
make seed_test_city
go build && ./metro
```

**Test City Features:**

- 12 stations across 4 lines
- 5 pre-configured trains
- Pre-generated schedules
- Perfect for quick testing and development

### Option 2: Santo Domingo (Real-World Data)

```bash
make setup_santo_domingo
go build && ./metro
```

**Santo Domingo Features:**

- 34 stations across 2 lines
- 69 trains (40 on Línea 1, 29 on Línea 2)
- Real coordinates from OpenStreetMap
- Auto-generated schedules based on real-world physics

## Available Commands

### Quick Setup Commands

```bash
# Test city (12 stations, 5 trains)
make seed_test_city

# Santo Domingo complete setup (34 stations, 69 trains)
make setup_santo_domingo
```

### Manual Santo Domingo Setup

If you want more control over the process:

```bash
# 1. Import fresh data from OpenStreetMap
make import_osm

# 2. Generate trains (40 for L1, 29 for L2)
make generate_santo_domingo_trains

# 3. Load everything (stations, lines, trains)
make seed_santo_domingo

# 4. Generate schedules for all trains
make generate_schedules
```

### Switching Between Cities

To switch from one city to another:

```bash
# Switch to test city
make seed_test_city
go build && ./metro

# Switch to Santo Domingo
make seed_santo_domingo  # or make setup_santo_domingo for fresh import
go build && ./metro
```

The `clean_city_data` step is automatically run before each seed, ensuring a clean slate.

## Database Maintenance

```bash
# Clear all city data (keeps migrations/schema)
make clean_city_data

# Run database migrations
make run_migrations
```

## Updating OSM Data

To refresh Santo Domingo data from OpenStreetMap:

```bash
# Full refresh (imports, generates trains, loads, creates schedules)
make setup_santo_domingo

# Or step by step:
make import_osm                    # Fetch fresh OSM data
make generate_santo_domingo_trains # Regenerate train fleet
make seed_santo_domingo            # Load into database
make generate_schedules            # Regenerate schedules
```

## Technical Details

### What Gets Cleaned

When you run `clean_city_data` or switch cities, these tables are cleared:

- `station` - All stations
- `line` - All metro lines
- `train` - All trains
- `edge` - Station connections
- `edge_point` - Path waypoints
- `station_line` - Station-to-line associations
- `passenger` - All passengers
- `schedule` - All train schedules

### What Stays

- Database schema (migrations)
- Train types (`make` table)
- Configuration settings

### File Locations

**Test City Seeds:**

- `data/sql/seeds/20240326000000_makes.sql` - Train types (shared)
- `data/sql/seeds/20240326002927_initial.sql` - Stations, lines, trains
- `data/sql/seeds/20240327153529_adding_edges.sql` - Station connections
- `data/sql/seeds/20251001200000_schedules.sql` - Train schedules

**Santo Domingo Seeds:**

- `data/sql/seeds/20240326000000_makes.sql` - Train types (shared)
- `data/sql/seeds/*_santo_domingo.sql` - Latest OSM import
- `data/sql/seeds/*_santo_domingo_trains.sql` - Latest train generation
- `data/sql/seeds/schedules_generated.sql` - Auto-generated schedules

### Schedule Generation

Schedules are automatically calculated based on:

- Station distances (real-world meters)
- Train physics (acceleration, top speed)
- Dwell time at stations (45 seconds)
- Operating hours (8:00 AM to 10:00 PM)

Each train runs continuous loops of its line throughout the day.

## Troubleshooting

**Problem:** Trains not moving

- **Solution:** Make sure schedules were generated: `make generate_schedules`

**Problem:** Wrong city data appearing

- **Solution:** Clean and reseed: `make clean_city_data && make seed_test_city`

**Problem:** Missing stations on Santo Domingo

- **Solution:** Reimport OSM data: `make setup_santo_domingo`

**Problem:** Database locked error

- **Solution:** Close the running simulation, then reseed

## Examples

### Daily Development Workflow

```bash
# Start with test city for quick testing
make seed_test_city
go build && ./metro

# Make changes to code...

# Switch to Santo Domingo for real-world testing
make seed_santo_domingo
go build && ./metro
```

### Updating Santo Domingo Data

```bash
# When you want fresh OpenStreetMap data
make setup_santo_domingo
go build && ./metro
```

### Creating Custom Train Configurations

```bash
# Import base data
make import_osm

# Edit tools/generate_santo_domingo_trains.go to customize
# (change train count, distribution, etc.)

# Regenerate everything
make generate_santo_domingo_trains
make seed_santo_domingo
make generate_schedules
```
