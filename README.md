# Metro Simulation

Real-time metro train simulation with Go, Ebiten, and SQLite. Features realistic physics, passenger management, and AI performance monitoring.

## Prerequisites

- Go 1.21+
- SQLite3
- Make

## Run

**Test City:**

```bash
make seed_test_city
go build && ./metro
```

**Santo Domingo (real-world data):**

```bash
make setup_santo_domingo
go build && ./metro
```

## Controls

- **Zoom:** Mouse wheel or `+`/`-`
- **Pan:** Arrow keys or `WASD`
- **Reset:** `R`
- **Click:** Stations/trains for details, score panel for metrics, newspaper button for reports

## Key Commands

```bash
make seed_test_city          # 12 stations, 5 trains
make setup_santo_domingo     # 34 stations, 69 trains (from OSM)
make clean_city_data         # Clear database
make run_migrations          # Setup schema
```

## Features

- Real-time physics-based train movement
- Passenger system with sentiment tracking
- Schedule-based operation (8 AM - 10 PM)
- Santo Domingo data from OpenStreetMap
- Camera zoom and pan
- AI monitoring (Tenjin) with performance metrics
- Auto-generated daily newspaper

## Docs

- [`docs/CITY-SWITCHING.md`](docs/CITY-SWITCHING.md) - City management
- [`docs/TENJIN-SUMMARY.md`](docs/TENJIN-SUMMARY.md) - AI architecture
- [`docs/OSM-IMPORT.md`](docs/OSM-IMPORT.md) - OpenStreetMap integration
