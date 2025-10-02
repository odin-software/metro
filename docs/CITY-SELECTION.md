# City Setup System

**Created**: October 1, 2025
**Status**: Complete ✅

## Overview

The Metro simulation runs with **ONE city at a time**. You choose which city to load when setting up the database. The simulation will then run with only that city's data - never both cities simultaneously.

## Quick Usage

```bash
# See all available commands
make help

# Setup with test city (choose this OR Santo Domingo, not both)
make seed_test_city

# Setup with Santo Domingo
make seed_santo_domingo
```

After running one of these commands, the database contains only that city's data.

## Available Cities

### Test City (Development)

**Stats**:

- 12 stations
- 4 lines
- 5 trains with schedules
- ~70 × 50 km (unrealistically large for testing edge cases)

**Use Case**: Development, testing features, debugging

**Command**: `make seed_test_city`

### Santo Domingo (Real Data)

**Stats**:

- 19 stations (real coordinates from OpenStreetMap)
- 2 lines (Línea 1 & Línea 2)
- ~12.7 × 7.4 km (realistic urban metro)
- No trains/schedules yet (generate with tools/generate_schedules.go)

**Use Case**: Real-world simulation, accurate timing, presentations

**Command**: `make seed_santo_domingo`

## How It Works

### Architecture

1. **Migrations** (schema): Shared across all cities
2. **Seeds** (data): City-specific with selective loading
3. **Makes** (train types): Shared across all cities

### Loading Process

When you run `make seed_test_city` or `make seed_santo_domingo`:

1. Runs migrations (ensures schema exists)
2. Cleans all city data (stations, lines, trains, schedules, edges)
3. Loads shared data (train types/makes)
4. Loads city-specific data (stations, lines, connections)

### File Structure

```
data/sql/seeds/
├── 20240326000000_makes.sql           # Shared train types
├── 20240326002927_initial.sql         # Test city data
├── 20240327153529_adding_edges.sql    # Test city edges
├── 20251001200000_schedules.sql       # Test city schedules
├── 20251001171034_santo_domingo.sql   # Santo Domingo data
├── test_city.sh                       # Test city loader script
└── santo_domingo.sh                   # Santo Domingo loader script
```

## Commands Reference

### Setup Commands

**Load Test City**:

```bash
make seed_test_city
```

Output: `✓ Test city loaded (12 stations, 4 lines, 5 trains)`

**Load Santo Domingo**:

```bash
make seed_santo_domingo
```

Output: `✓ Santo Domingo loaded (19 stations, 2 lines)`

**Import Fresh OSM Data**:

```bash
make import_osm
make seed_santo_domingo
```

### Management Commands

**Clean City Data** (keeps schema):

```bash
make clean_city_data
```

**Full Database Reset**:

```bash
make reset_db
```

**Run All Seeds** (loads both cities + duplicates):

```bash
make run_seeds  # Not recommended - loads everything
```

## Changing Cities

To use a different city, run the appropriate seed command. The database is automatically cleaned first:

```bash
# Want to use test city
make seed_test_city

# Want to use Santo Domingo
make seed_santo_domingo
```

Each command cleans all city data first, then loads only the selected city. You never have both cities in the database at once.

## Adding New Cities

To add a new city (e.g., New York, Tokyo):

### 1. Update OSM Import Tool

Edit `tools/import_osm.go`:

- Add city-specific bounding box coordinates
- Add reference point for coordinate conversion
- Add line color mappings
- Add known station configurations (if needed)

### 2. Generate Seed File

```bash
make import_osm
```

This creates: `data/sql/seeds/{timestamp}_city_name.sql`

### 3. Create Loader Script

Create `data/sql/seeds/city_name.sh`:

```bash
#!/bin/bash
DB_PATH="${1:-data/metro.db}"
echo "Loading City Name..."

sed -n '/-- +goose Up/,/-- +goose Down/p' data/sql/seeds/20240326000000_makes.sql | \
  sed '/-- +goose/d' | sqlite3 "$DB_PATH"

LATEST=$(ls -t data/sql/seeds/*city_name.sql 2>/dev/null | head -1)
sed -n '/-- +goose Up/,/-- +goose Down/p' "$LATEST" | \
  sed '/-- +goose/d' | sqlite3 "$DB_PATH"

echo "✓ City Name loaded"
```

Make it executable: `chmod +x data/sql/seeds/city_name.sh`

### 4. Add Makefile Target

Edit `Makefile`:

```makefile
seed_city_name: run_migrations clean_city_data
	@bash data/sql/seeds/city_name.sh $(GOOSE_DBSTRING)
```

### 5. Update Documentation

- Add city to `make help`
- Update README.md
- Update this document

## Technical Details

### ID Ranges

To prevent conflicts between cities, use different ID ranges:

- **Test City**: 1-99
- **Santo Domingo**: 1000-1999
- **Future City 1**: 2000-2999
- **Future City 2**: 3000-3999

### Shared Data

**Train Makes** (types) are shared:

- ID 1: Fast metro (70 km/h)
- ID 2: Standard metro (60 km/h)

These are loaded by all city seed scripts.

### Scale Consistency

All cities use: **1 pixel = 100 meters**

This ensures:

- Train speeds remain realistic across cities
- Distance calculations are consistent
- Journey times are accurate

### Database Tables Affected

**Cleaned on City Switch**:

- `station` - Station locations
- `line` - Metro lines
- `train` - Active trains
- `edge` - Station connections
- `edge_point` - Path waypoints
- `station_line` - Station-line associations
- `schedule` - Timetables
- `passenger` - Active passengers

**Preserved**:

- `make` - Train types
- Schema/migrations

## Troubleshooting

### "No Santo Domingo seed file found"

**Cause**: Haven't imported OSM data yet

**Fix**:

```bash
make import_osm
make seed_santo_domingo
```

### "No trains appear in Santo Domingo"

**Expected**: Santo Domingo data doesn't include trains by default

**Fix**: Generate schedules which will create trains:

```bash
cd tools
go run generate_schedules.go
```

### Seeds running but no data appears

**Check**: Verify data was loaded

```bash
sqlite3 data/metro.db "SELECT COUNT(*) FROM station;"
```

**Fix**: Ensure migrations ran first

```bash
make run_migrations
make seed_test_city
```

### Both cities appear at once

**Cause**: Used `make run_seeds` instead of city-specific command

**Fix**: Clean and reload specific city

```bash
make seed_test_city
```

## Future Enhancements

### Runtime City Selection

Add a city selector in the UI:

- Config setting for active city
- Load city on startup
- Switch without database wipe

### Multiple Cities Simultaneously

Support loading multiple cities with:

- City ID in each table
- Viewport selection
- Memory optimization

### City Metadata

Add `city` table:

```sql
CREATE TABLE city (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    country TEXT,
    real_world BOOLEAN,
    center_lat REAL,
    center_lon REAL,
    created_at INTEGER
);
```

## Summary

The city selection system provides:

- ✅ Easy switching between test and real cities
- ✅ Clean separation of city data
- ✅ Shared infrastructure (train types, schema)
- ✅ Extensible for new cities
- ✅ Simple commands (`make seed_X`)
- ✅ Automatic cleanup when switching

No manual database manipulation required - just run the appropriate `make` command!
