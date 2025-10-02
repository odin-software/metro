# Metro Tools

## import_osm.go - OpenStreetMap Data Importer

Imports real metro data from OpenStreetMap for Santo Domingo, Dominican Republic.

### Usage

```bash
make import_osm
```

This will:

1. Query OpenStreetMap Overpass API for Santo Domingo metro stations
2. Fall back to curated data if API is unavailable
3. Convert lat/lon coordinates to pixels (1 pixel = 100 meters)
4. Generate SQL seed file with stations, lines, and edges
5. Save to `data/sql/seeds/{timestamp}_santo_domingo.sql`

### Apply the Data

After generation, apply to database:

```bash
make run_seeds
```

### Data Generated

**Línea 1** (14 stations):

- Villa Mella → Hermanas Mirabal → Máximo Gómez → Bartolomé Colón → Mama Tingó → Juan Pablo Duarte → María Montez → Casandra Damirón → Francisco del Rosario Sánchez → Pedro Livio Cedeño → Los Taínos → Centro de Los Héroes → Joaquín Balaguer → Eduardo Brito

**Línea 2** (5 stations):

- Francisco Alberto Caamaño Deñó → La Julia → Antonio Duvergé → Hermanas Mirabal 2 → Manuel María Valencia

### Real-World Coordinates

- Reference point: 18.47°N, 69.91°W (approximate center of Santo Domingo)
- Scale: 1 pixel = 100 meters
- Map size: ~10km × 10km
- Station IDs: 1000-1018 (to avoid conflicts with test data)
- Line IDs: 100-101

### Refreshing Data

The tool can be run multiple times to refresh data. It uses `INSERT OR IGNORE` so duplicate runs are safe.

### Customization

Edit `import_osm.go` to:

- Change bounding box (MinLat, MaxLat, MinLon, MaxLon)
- Adjust reference point for coordinate conversion
- Modify line colors
- Add more cities

### Technical Details

**Coordinate Conversion:**

- Uses Equirectangular projection
- Converts lat/lon degrees to meters using:
  - 111,000 meters per degree latitude
  - 111,000 × cos(latitude) meters per degree longitude
- Converts meters to pixels using PixelsPerMeter = 0.01

**Line Assignment:**

- Tries OSM relations first
- Falls back to known station configuration
- Properly handles transfer stations (Centro de Los Héroes)

**Edge Generation:**

- Automatically connects adjacent stations on each line
- Creates bidirectional connectivity for train pathfinding
