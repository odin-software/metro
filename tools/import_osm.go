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

	var stations []Station
	var lines []Line

	// Query Overpass API for real coordinates
	fmt.Println("Querying Overpass API for station coordinates...")
	query := buildOverpassQuery()
	response, err := queryOverpass(query)

	if err != nil {
		fmt.Printf("Warning: Overpass API failed: %v\n", err)
		fmt.Println("Falling back to curated data with estimated coordinates...\n")
		stations, lines = getCuratedData()
	} else {
		fmt.Printf("Found %d elements from OSM API\n", len(response.Elements))
		fmt.Println("Using OSM coordinates with curated line ordering...\n")

		// Get stations with real OSM coordinates
		stations = parseStations(response)

		// Apply curated line ordering (OSM relations have wrong order)
		lines = applyCuratedLineOrdering(stations)
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

	if err := generateSQL(filename, stations, lines); err != nil {
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

area["name"="Santo Domingo"]->.sd;
area["name"="Distrito Nacional"]->.dn;

(
  node["railway"="station"]["station"="subway"][!"opening_date"]["name"!="San Felipe"](area.sd);
  node["railway"="station"]["subway"="yes"](area.dn);
);
out body;`,
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
		// Línea 1 (North to South) - 16 stations
		{"Mamá Tingó", 18.5100, -69.8950, "1"},
		{"Gregorio Urbano Gilbert", 18.5050, -69.8970, "1"},
		{"Gregorio Luperón", 18.5000, -69.8990, "1"},
		{"José Francisco Peña Gómez", 18.4950, -69.9010, "1"},
		{"Hermanas Mirabal", 18.4900, -69.9030, "1"},
		{"Máximo Gómez", 18.4850, -69.9050, "1"},
		{"Los Taínos", 18.4800, -69.9070, "1"},
		{"Pedro Livio Cedeño", 18.4750, -69.9090, "1"},
		{"Manuel Arturo Peña Batlle", 18.4710, -69.9110, "1"},
		{"Juan Pablo Duarte", 18.4670, -69.9130, "1"},
		{"Profesor Juan Bosch", 18.4630, -69.9150, "1"},
		{"Casandra Damirón", 18.4590, -69.9170, "1"},
		{"Joaquín Balaguer", 18.4550, -69.9190, "1"},
		{"Amín Abel Hasbún", 18.4510, -69.9210, "1"},
		{"Francisco Alberto Caamaño Deñó", 18.4470, -69.9230, "1,2"}, // Transfer station
		{"Centro de los Héroes", 18.4430, -69.9250, "1"},

		// Línea 2 (East to West) - 18 stations
		{"María Montez", 18.4800, -69.8800, "2"},
		{"Pedro Francisco Bonó", 18.4790, -69.8850, "2"},
		{"Francisco Gregorio Billini", 18.4780, -69.8900, "2"},
		{"Ulises Francisco Espaillat", 18.4770, -69.8950, "2"},
		{"Pedro Mir", 18.4760, -69.9000, "2"},
		{"Freddy Beras-Goico", 18.4750, -69.9050, "2"},
		{"Juan Ulises García Saleta", 18.4740, -69.9100, "2"},
		{"Juan Pablo Duarte", 18.4670, -69.9130, "2"}, // Transfer station with L1
		{"Coronel Rafael Tomás Fernández Domínguez", 18.4700, -69.9180, "2"},
		{"Mauricio Báez", 18.4690, -69.9230, "2"},
		{"Ramón Cáceres", 18.4680, -69.9280, "2"},
		{"Horacio Vásquez", 18.4670, -69.9330, "2"},
		{"Manuel de Jesús Galván", 18.4660, -69.9380, "2"},
		{"Eduardo Brito", 18.4650, -69.9430, "2"},
		{"Ercilia Pepín", 18.4640, -69.9480, "2"},
		{"Rosa Duarte", 18.4630, -69.9530, "2"},
		{"Trina de Moya de Vásquez", 18.4620, -69.9580, "2"},
		{"Concepción Bona", 18.4610, -69.9630, "2"},
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
				stationIDs["Mamá Tingó"],
				stationIDs["Gregorio Urbano Gilbert"],
				stationIDs["Gregorio Luperón"],
				stationIDs["José Francisco Peña Gómez"],
				stationIDs["Hermanas Mirabal"],
				stationIDs["Máximo Gómez"],
				stationIDs["Los Taínos"],
				stationIDs["Pedro Livio Cedeño"],
				stationIDs["Manuel Arturo Peña Batlle"],
				stationIDs["Juan Pablo Duarte"],
				stationIDs["Profesor Juan Bosch"],
				stationIDs["Casandra Damirón"],
				stationIDs["Joaquín Balaguer"],
				stationIDs["Amín Abel Hasbún"],
				stationIDs["Francisco Alberto Caamaño Deñó"],
				stationIDs["Centro de los Héroes"],
			},
		},
		{
			Name:  "Línea 2",
			Color: "#0066B3",
			Stations: []int64{
				stationIDs["María Montez"],
				stationIDs["Pedro Francisco Bonó"],
				stationIDs["Francisco Gregorio Billini"],
				stationIDs["Ulises Francisco Espaillat"],
				stationIDs["Pedro Mir"],
				stationIDs["Freddy Beras-Goico"],
				stationIDs["Juan Ulises García Saleta"],
				stationIDs["Juan Pablo Duarte"],
				stationIDs["Coronel Rafael Tomás Fernández Domínguez"],
				stationIDs["Mauricio Báez"],
				stationIDs["Ramón Cáceres"],
				stationIDs["Horacio Vásquez"],
				stationIDs["Manuel de Jesús Galván"],
				stationIDs["Eduardo Brito"],
				stationIDs["Ercilia Pepín"],
				stationIDs["Rosa Duarte"],
				stationIDs["Trina de Moya de Vásquez"],
				stationIDs["Concepción Bona"],
			},
		},
	}

	return stations, lines
}

// applyCuratedLineOrdering takes stations from OSM (with real coordinates)
// and applies the correct line ordering as specified by local knowledge
func applyCuratedLineOrdering(stations []Station) []Line {
	// Build a map from station name to station
	stationMap := make(map[string]*Station)
	for i := range stations {
		stationMap[stations[i].Name] = &stations[i]
	}

	// Define correct line ordering
	linea1Order := []string{
		"Mamá Tingó",
		"Gregorio Urbano Gilbert",
		"Gregorio Luperón",
		"José Francisco Peña Gómez",
		"Hermanas Mirabal",
		"Máximo Gómez",
		"Los Taínos",
		"Pedro Livio Cedeño",
		"Manuel Arturo Peña Batlle",
		"Juan Pablo Duarte",
		"Profesor Juan Bosch",
		"Casandra Damirón",
		"Joaquín Balaguer",
		"Amín Abel Hasbún",
		"Francisco Alberto Caamaño Deñó",
		"Centro de los Héroes",
	}

	linea2Order := []string{
		"María Montez",
		"Pedro Francisco Bonó",
		"Francisco Gregorio Billini",
		"Ulises Francisco Espaillat",
		"Pedro Mir",
		"Freddy Beras-Goico",
		"Juan Ulises García Saleta",
		"Juan Pablo Duarte",
		"Coronel Rafael Tomás Fernández Domínguez",
		"Mauricio Báez",
		"Ramón Cáceres",
		"Horacio Vásquez",
		"Manuel de Jesús Galván",
		"Eduardo Brito",
		"Ercilia Pepín",
		"Rosa Duarte",
		"Trina de Moya de Vásquez",
		"Concepción Bona",
	}

	// Build lines with correct ordering
	lines := []Line{}

	// Línea 1
	linea1Stations := []int64{}
	for _, name := range linea1Order {
		if st, ok := stationMap[name]; ok {
			linea1Stations = append(linea1Stations, st.OSMID)
		} else {
			fmt.Printf("Warning: Station '%s' not found in OSM data for Línea 1\n", name)
		}
	}
	if len(linea1Stations) > 0 {
		lines = append(lines, Line{
			Name:     "Línea 1",
			Color:    "#E84B28",
			Stations: linea1Stations,
		})
	}

	// Línea 2
	linea2Stations := []int64{}
	for _, name := range linea2Order {
		if st, ok := stationMap[name]; ok {
			linea2Stations = append(linea2Stations, st.OSMID)
		} else {
			fmt.Printf("Warning: Station '%s' not found in OSM data for Línea 2\n", name)
		}
	}
	if len(linea2Stations) > 0 {
		lines = append(lines, Line{
			Name:     "Línea 2",
			Color:    "#0066B3",
			Stations: linea2Stations,
		})
	}

	return lines
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
