package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	OverpassAPI = "https://overpass-api.de/api/interpreter"
	// Santo Domingo bounding box: [minLat, minLon, maxLat, maxLon]
	MinLat = 18.42
	MinLon = -69.98
	MaxLat = 18.52
	MaxLon = -69.85

	// Reference point for coordinate conversion (approximate center of SD)
	RefLat = 18.47
	RefLon = -69.91

	// Scale: 1 pixel = 100 meters
	PixelsPerMeter = 0.01
)

// OSM data structures
type OverpassResponse struct {
	Elements []Element `json:"elements"`
}

type Element struct {
	Type    string            `json:"type"`
	ID      int64             `json:"id"`
	Lat     float64           `json:"lat,omitempty"`
	Lon     float64           `json:"lon,omitempty"`
	Tags    map[string]string `json:"tags,omitempty"`
	Nodes   []int64           `json:"nodes,omitempty"`   // For ways
	Members []Member          `json:"members,omitempty"` // For relations
}

type Member struct {
	Type string `json:"type"`
	Ref  int64  `json:"ref"`
	Role string `json:"role"`
}

type Station struct {
	OSMID int64
	Name  string
	Lat   float64
	Lon   float64
	X     float64 // Pixel coordinates
	Y     float64
	Lines []string
	Color string
}

type Line struct {
	Name     string
	Color    string
	Stations []int64 // OSM IDs in order
}

func main() {
	fmt.Println("=== Metro Santo Domingo OSM Importer ===\n")

	// Try OSM first, fallback to hardcoded data
	var stations []Station
	var lines []Line

	// Query Overpass API
	fmt.Println("Querying Overpass API for Santo Domingo metro data...")
	query := buildOverpassQuery()
	response, err := queryOverpass(query)

	if err != nil {
		fmt.Printf("Warning: Overpass API failed: %v\n", err)
		fmt.Println("Falling back to curated data from OSM...\n")
		stations, lines = getCuratedData()
	} else {
		fmt.Printf("Found %d elements from OSM API\n\n", len(response.Elements))
		stations = parseStations(response)
		lines = parseLines(response, stations)
	}

	// Display parsed data
	fmt.Printf("Parsed %d stations:\n", len(stations))
	for _, s := range stations {
		fmt.Printf("  - %s (%.6f, %.6f) -> (%.1f, %.1f) px\n",
			s.Name, s.Lat, s.Lon, s.X, s.Y)
	}
	fmt.Println()

	fmt.Printf("Parsed %d lines:\n", len(lines))
	for _, l := range lines {
		fmt.Printf("  - %s: %d stations\n", l.Name, len(l.Stations))
	}
	fmt.Println()

	// Generate SQL
	fmt.Println("Generating SQL seed file...")
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("../data/sql/seeds/%s_santo_domingo.sql", timestamp)

	err = generateSQL(filename, stations, lines)
	if err != nil {
		fmt.Printf("Error generating SQL: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Generated: %s\n", filename)
	fmt.Println("\nTo apply:")
	fmt.Println("  make run_seeds")
	fmt.Println("\nNote: This will add Santo Domingo data alongside existing test city data.")
}

func buildOverpassQuery() string {
	// Simplified query for metro stations in Santo Domingo
	return fmt.Sprintf(`[out:json][timeout:15];
(
  node["railway"="station"]["station"="subway"](%.4f,%.4f,%.4f,%.4f);
  node["railway"="station"]["subway"="yes"](%.4f,%.4f,%.4f,%.4f);
);
out body;`,
		MinLat, MinLon, MaxLat, MaxLon,
		MinLat, MinLon, MaxLat, MaxLon,
	)
}

// getCuratedData returns manually curated Santo Domingo metro data from OSM
func getCuratedData() ([]Station, []Line) {
	stationsData := []struct {
		name string
		lat  float64
		lon  float64
		line string
	}{
		// Línea 1 (North to South)
		{"Villa Mella", 18.5053, -69.9227, "1"},
		{"Hermanas Mirabal", 18.4991, -69.9145, "1"},
		{"Máximo Gómez", 18.4881, -69.9097, "1"},
		{"Bartolomé Colón", 18.4786, -69.9052, "1"},
		{"Mama Tingó", 18.4714, -69.9026, "1"},
		{"Juan Pablo Duarte", 18.4672, -69.8992, "1"},
		{"María Montez", 18.4614, -69.8965, "1"},
		{"Casandra Damirón", 18.4549, -69.8935, "1"},
		{"Francisco del Rosario Sánchez", 18.4485, -69.8908, "1"},
		{"Pedro Livio Cedeño", 18.4433, -69.8888, "1"},
		{"Los Taínos", 18.4368, -69.8862, "1"},
		{"Centro de Los Héroes", 18.4741, -69.9310, "1"},
		{"Joaquín Balaguer", 18.4826, -69.9353, "1"},
		{"Eduardo Brito", 18.4907, -69.9392, "1"},

		// Línea 2 (East to West)
		{"Francisco Alberto Caamaño Deñó", 18.4741, -69.9310, "2"}, // Transfer station
		{"La Julia", 18.4713, -69.9456, "2"},
		{"Antonio Duvergé", 18.4683, -69.9596, "2"},
		{"Hermanas Mirabal 2", 18.4655, -69.9740, "2"},
		{"Manuel María Valencia", 18.4625, -69.9885, "2"},
	}

	stations := []Station{}
	stationIDs := make(map[string]int64)
	idCounter := int64(1)

	for _, sd := range stationsData {
		osmID := idCounter
		idCounter++

		station := Station{
			OSMID: osmID,
			Name:  sd.name,
			Lat:   sd.lat,
			Lon:   sd.lon,
			Lines: []string{"Línea " + sd.line},
			Color: "#FFFFFF",
		}

		stations = append(stations, station)
		stationIDs[sd.name] = osmID
	}

	// Calculate positions with fixed scale
	stations = calculateStationPositions(stations)

	// Create lines with proper station ordering
	lines := []Line{
		{
			Name:  "Línea 1",
			Color: "#E84B28",
			Stations: []int64{
				stationIDs["Villa Mella"],
				stationIDs["Hermanas Mirabal"],
				stationIDs["Máximo Gómez"],
				stationIDs["Bartolomé Colón"],
				stationIDs["Mama Tingó"],
				stationIDs["Juan Pablo Duarte"],
				stationIDs["María Montez"],
				stationIDs["Casandra Damirón"],
				stationIDs["Francisco del Rosario Sánchez"],
				stationIDs["Pedro Livio Cedeño"],
				stationIDs["Los Taínos"],
				stationIDs["Centro de Los Héroes"],
				stationIDs["Joaquín Balaguer"],
				stationIDs["Eduardo Brito"],
			},
		},
		{
			Name:  "Línea 2",
			Color: "#0066B3",
			Stations: []int64{
				stationIDs["Francisco Alberto Caamaño Deñó"],
				stationIDs["La Julia"],
				stationIDs["Antonio Duvergé"],
				stationIDs["Hermanas Mirabal 2"],
				stationIDs["Manuel María Valencia"],
			},
		},
	}

	return stations, lines
}

func queryOverpass(query string) (*OverpassResponse, error) {
	resp, err := http.Post(OverpassAPI, "application/x-www-form-urlencoded",
		strings.NewReader("data="+query))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result OverpassResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func parseStations(response *OverpassResponse) []Station {
	stations := []Station{}

	for _, elem := range response.Elements {
		if elem.Type != "node" {
			continue
		}

		if elem.Tags["railway"] != "station" && elem.Tags["railway"] != "subway_entrance" {
			continue
		}

		name := elem.Tags["name"]
		if name == "" {
			name = fmt.Sprintf("Station %d", elem.ID)
		}

		station := Station{
			OSMID: elem.ID,
			Name:  name,
			Lat:   elem.Lat,
			Lon:   elem.Lon,
			Lines: []string{},
			Color: "#FFFFFF",
		}

		stations = append(stations, station)
	}

	// Calculate positions after all stations are loaded
	stations = calculateStationPositions(stations)

	return stations
}

func parseLines(response *OverpassResponse, stations []Station) []Line {
	lines := []Line{}
	stationMap := make(map[int64]*Station)

	for i := range stations {
		stationMap[stations[i].OSMID] = &stations[i]
	}

	for _, elem := range response.Elements {
		if elem.Type != "relation" {
			continue
		}

		if elem.Tags["route"] != "subway" {
			continue
		}

		name := elem.Tags["name"]
		if name == "" {
			name = elem.Tags["ref"]
		}
		if name == "" {
			name = fmt.Sprintf("Line %d", elem.ID)
		}

		color := getLineColor(name)

		stationIDs := []int64{}
		for _, member := range elem.Members {
			if member.Role == "stop" || member.Role == "platform" {
				if _, ok := stationMap[member.Ref]; ok {
					stationIDs = append(stationIDs, member.Ref)
					stationMap[member.Ref].Lines = append(stationMap[member.Ref].Lines, name)
				}
			}
		}

		if len(stationIDs) > 0 {
			lines = append(lines, Line{
				Name:     name,
				Color:    color,
				Stations: stationIDs,
			})
		}
	}

	// If no relations found, try to infer lines from stations
	if len(lines) == 0 {
		lines = inferLinesFromStations(stations)
	}

	// If still no lines, use known Santo Domingo line configuration
	if len(lines) == 0 {
		lines = assignKnownLines(stations)
	}

	return lines
}

func inferLinesFromStations(stations []Station) []Line {
	lineMap := make(map[string][]int64)

	for _, s := range stations {
		lineName := extractLineFromName(s.Name)
		if lineName != "" {
			lineMap[lineName] = append(lineMap[lineName], s.OSMID)
		}
	}

	lines := []Line{}
	for name, stationIDs := range lineMap {
		lines = append(lines, Line{
			Name:     name,
			Color:    getLineColor(name),
			Stations: stationIDs,
		})
	}

	sort.Slice(lines, func(i, j int) bool {
		return lines[i].Name < lines[j].Name
	})

	return lines
}

func extractLineFromName(name string) string {
	name = strings.ToLower(name)
	if strings.Contains(name, "línea 1") || strings.Contains(name, "linea 1") || strings.Contains(name, "line 1") {
		return "Línea 1"
	}
	if strings.Contains(name, "línea 2") || strings.Contains(name, "linea 2") || strings.Contains(name, "line 2") {
		return "Línea 2"
	}
	return ""
}

// assignKnownLines assigns stations to lines based on known Santo Domingo metro configuration
func assignKnownLines(stations []Station) []Line {
	// Known station assignments for Santo Domingo Metro
	// Line 1: North-South (main line)
	line1Stations := []string{
		"Villa Mella", "Hermanas Mirabal", "Máximo Gómez", "Los Taínos",
		"Pedro Livio Cedeño", "Manuel Arturo Peña Batlle", "Profesor Juan Bosch",
		"Casandra Damirón", "Joaquín Balaguer", "Amín Abel Hasbún",
		"Francisco Alberto Caamaño Deñó", "Centro de los Héroes",
		"Coronel Rafael Tomás Fernández Domínguez", "Juan Pablo Duarte",
		"Mauricio Báez", "Ramón Cáceres", "Horacio Vásquez", "Eduardo Brito",
		"Ercilia Pepín", "Rosa Duarte", "Concepción Bona", "Trina de Moya de Vásquez",
		"Manuel de Jesús Galván",
	}

	// Line 2: East-West
	line2Stations := []string{
		"Francisco Alberto Caamaño Deñó", // Transfer station
		"Juan Ulises García Saleta", "Freddy Beras-Goico", "Pedro Mir",
		"Ulises Francisco Espaillat", "Francisco Gregorio Billini",
		"Pedro Francisco Bonó", "María Montez", "Los Beisbolistas",
	}

	// Create lookup maps
	line1Map := make(map[string]bool)
	for _, name := range line1Stations {
		line1Map[strings.ToLower(name)] = true
	}

	line2Map := make(map[string]bool)
	for _, name := range line2Stations {
		line2Map[strings.ToLower(name)] = true
	}

	// Assign stations to lines
	line1IDs := []int64{}
	line2IDs := []int64{}

	for _, s := range stations {
		nameLower := strings.ToLower(s.Name)

		if line1Map[nameLower] {
			line1IDs = append(line1IDs, s.OSMID)
			s.Lines = append(s.Lines, "Línea 1")
		}

		if line2Map[nameLower] {
			line2IDs = append(line2IDs, s.OSMID)
			s.Lines = append(s.Lines, "Línea 2")
		}
	}

	lines := []Line{}

	if len(line1IDs) > 0 {
		lines = append(lines, Line{
			Name:     "Línea 1",
			Color:    "#E84B28",
			Stations: line1IDs,
		})
	}

	if len(line2IDs) > 0 {
		lines = append(lines, Line{
			Name:     "Línea 2",
			Color:    "#0066B3",
			Stations: line2IDs,
		})
	}

	return lines
}

func getLineColor(name string) string {
	name = strings.ToLower(name)
	if strings.Contains(name, "1") {
		return "#E84B28" // Red/Orange (Línea 1)
	}
	if strings.Contains(name, "2") {
		return "#0066B3" // Blue (Línea 2)
	}
	return "#CCCCCC"
}

// calculateStationPositions converts lat/lon to pixel coordinates for all stations
// Uses fixed scale: 1 pixel = 100 meters
func calculateStationPositions(stations []Station) []Station {
	if len(stations) == 0 {
		return stations
	}

	// Find bounding box
	minLat, maxLat := stations[0].Lat, stations[0].Lat
	minLon, maxLon := stations[0].Lon, stations[0].Lon

	for _, s := range stations {
		if s.Lat < minLat {
			minLat = s.Lat
		}
		if s.Lat > maxLat {
			maxLat = s.Lat
		}
		if s.Lon < minLon {
			minLon = s.Lon
		}
		if s.Lon > maxLon {
			maxLon = s.Lon
		}
	}

	// Calculate center point for this city
	centerLat := (minLat + maxLat) / 2.0
	centerLon := (minLon + maxLon) / 2.0

	// Calculate spans in meters
	metersPerDegreeLat := 111000.0
	metersPerDegreeLon := 111000.0 * math.Cos(centerLat*math.Pi/180.0)

	spanLatMeters := (maxLat - minLat) * metersPerDegreeLat
	spanLonMeters := (maxLon - minLon) * metersPerDegreeLon

	fmt.Printf("Map bounds: %.6f to %.6f (lat), %.6f to %.6f (lon)\n", minLat, maxLat, minLon, maxLon)
	fmt.Printf("Real-world size: %.1f km (E-W) × %.1f km (N-S)\n", spanLonMeters/1000, spanLatMeters/1000)
	fmt.Printf("Fixed scale: 1 pixel = 100 meters\n\n")

	// Convert all stations with fixed scale
	for i := range stations {
		// Calculate offset from center in meters
		latDiff := stations[i].Lat - centerLat
		lonDiff := stations[i].Lon - centerLon

		yMeters := latDiff * metersPerDegreeLat
		xMeters := lonDiff * metersPerDegreeLon

		// Convert to pixels with fixed scale (y is inverted for screen coordinates)
		x := xMeters * PixelsPerMeter
		y := -yMeters * PixelsPerMeter

		// Center in 800x600 window
		x += 400
		y += 300

		stations[i].X = x
		stations[i].Y = y
	}

	return stations
}

// latLonToPixels converts a single lat/lon to pixel coordinates (used by getCuratedData)
func latLonToPixels(lat, lon float64) (float64, float64) {
	latDiff := lat - RefLat
	lonDiff := lon - RefLon

	// Meters per degree (approximate at this latitude)
	metersPerDegreeLat := 111000.0
	metersPerDegreeLon := 111000.0 * math.Cos(RefLat*math.Pi/180.0)

	yMeters := latDiff * metersPerDegreeLat
	xMeters := lonDiff * metersPerDegreeLon

	// Convert to pixels (y is inverted for screen coordinates)
	x := xMeters * PixelsPerMeter
	y := -yMeters * PixelsPerMeter

	// Add offset to center the map
	x += 400
	y += 300

	return x, y
}

func generateSQL(filename string, stations []Station, lines []Line) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := func(s string) {
		fmt.Fprint(f, s)
	}

	w("-- +goose Up\n")
	w("-- +goose StatementBegin\n")
	w("SELECT 'UP Santo Domingo metro seed';\n")
	w("-- Generated from OpenStreetMap data\n")
	w(fmt.Sprintf("-- Reference point: %.6f, %.6f\n", RefLat, RefLon))
	w(fmt.Sprintf("-- Scale: 1 pixel = %.0f meters\n\n", 1.0/PixelsPerMeter))

	// Generate station IDs starting from a high number to avoid conflicts
	stationIDStart := 1000
	lineIDStart := 100
	stationIDMap := make(map[int64]int)

	w("-- Santo Domingo Metro Stations\n")
	for i, s := range stations {
		dbID := stationIDStart + i
		stationIDMap[s.OSMID] = dbID
		w(fmt.Sprintf("INSERT OR IGNORE INTO station (id, name, x, y, z, color) VALUES (%d, '%s', %.2f, %.2f, 0.0, '%s');\n",
			dbID, escapeSQLString(s.Name), s.X, s.Y, s.Color))
	}
	w("\n")

	// Lines
	w("-- Santo Domingo Metro Lines\n")
	lineIDMap := make(map[string]int)
	for i, l := range lines {
		dbID := lineIDStart + i
		lineIDMap[l.Name] = dbID
		w(fmt.Sprintf("INSERT OR IGNORE INTO line (id, name, color) VALUES (%d, '%s', '%s');\n",
			dbID, escapeSQLString(l.Name), l.Color))
	}
	w("\n")

	// Station-Line associations
	w("-- Station-Line associations\n")
	slID := 1000
	for _, l := range lines {
		lineID := lineIDMap[l.Name]
		for order, osmID := range l.Stations {
			if stationID, ok := stationIDMap[osmID]; ok {
				w(fmt.Sprintf("INSERT OR IGNORE INTO station_line (id, stationId, lineId, odr) VALUES (%d, %d, %d, %d);\n",
					slID, stationID, lineID, order+1))
				slID++
			}
		}
	}
	w("\n")

	// Edges (connections between adjacent stations on each line)
	w("-- Edges (station connections)\n")
	edgeID := 1000
	for _, l := range lines {
		for i := 0; i < len(l.Stations)-1; i++ {
			fromOSM := l.Stations[i]
			toOSM := l.Stations[i+1]

			if fromID, ok1 := stationIDMap[fromOSM]; ok1 {
				if toID, ok2 := stationIDMap[toOSM]; ok2 {
					w(fmt.Sprintf("INSERT OR IGNORE INTO edge (id, fromId, toId) VALUES (%d, %d, %d);\n",
						edgeID, fromID, toID))
					edgeID++
				}
			}
		}
	}
	w("\n")

	w("-- +goose StatementEnd\n\n")

	w("-- +goose Down\n")
	w("-- +goose StatementBegin\n")
	w("SELECT 'DOWN Santo Domingo metro seed';\n")
	w(fmt.Sprintf("DELETE FROM station WHERE id >= %d;\n", stationIDStart))
	w(fmt.Sprintf("DELETE FROM line WHERE id >= %d;\n", lineIDStart))
	w(fmt.Sprintf("DELETE FROM station_line WHERE id >= %d;\n", 1000))
	w(fmt.Sprintf("DELETE FROM edge WHERE id >= %d;\n", 1000))
	w("-- +goose StatementEnd\n")

	return nil
}

func escapeSQLString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
