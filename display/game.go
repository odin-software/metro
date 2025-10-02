package display

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/odin-software/metro/control"
	"github.com/odin-software/metro/internal/models"
	"github.com/odin-software/metro/internal/tenjin"
)

type SceneType int

const (
	SceneMap SceneType = iota
	SceneStation
	SceneNewspaper
)

type Game struct {
	trains             []models.Train
	stations           []*models.Station
	lines              []models.Line
	selectedTrain      *models.Train
	selectedStation    *models.Station
	currentScene       SceneType
	lastMouseClick     bool
	brain              *tenjin.Tenjin
	scoreBreakdownOpen bool
	clock              interface{ GetCurrentTime() string; GetCurrentTimeOfDay() int } // Simulation clock interface
	scheduleDB         ScheduleDB                                                       // Schedule database access

	// Camera controls
	cameraZoom   float64 // Zoom level (1.0 = normal, 2.0 = 2x zoom)
	cameraOffsetX float64 // Camera pan offset X
	cameraOffsetY float64 // Camera pan offset Y
}

// ScheduleDB interface for accessing schedule data
type ScheduleDB interface {
	GetScheduleByTrainAndStation(trainID, stationID int64) (Schedule, error)
	GetScheduleForTrain(trainID int64) ([]Schedule, error)
	GetScheduleForStation(stationID int64) ([]Schedule, error)
}

// Schedule represents a scheduled stop
type Schedule struct {
	TrainID       int64
	StationID     int64
	ScheduledTime int
	SequenceOrder int64
}

func NewGame(trains []models.Train, stations []*models.Station, lines []models.Line, brain *tenjin.Tenjin, clock interface{ GetCurrentTime() string; GetCurrentTimeOfDay() int }, scheduleDB ScheduleDB) *Game {
	Init()
	models.LineInit()
	return &Game{
		trains:       trains,
		stations:     stations,
		lines:        lines,
		currentScene: SceneMap,
		scheduleDB:   scheduleDB,
		brain:        brain,
		clock:        clock,
		cameraZoom:   1.0, // Start at 1x zoom
		cameraOffsetX: 0,
		cameraOffsetY: 0,
	}
}

func (g *Game) Update() error {
	for i := range g.trains {
		g.trains[i].Update()
	}
	for i := range g.stations {
		g.stations[i].Update()
	}

	// Handle camera controls (only in map scene)
	if g.currentScene == SceneMap {
		g.handleCameraControls()
	}

	// Handle mouse clicks
	g.handleMouseClick()

	return nil
}

// formatTime converts seconds since midnight to HH:MM format
func (g *Game) formatTime(secondsSinceMidnight int) string {
	hours := secondsSinceMidnight / 3600
	minutes := (secondsSinceMidnight % 3600) / 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

// handleCameraControls processes zoom and pan inputs
func (g *Game) handleCameraControls() {
	// Zoom controls: Mouse wheel or +/- keys
	_, wheelY := ebiten.Wheel()
	if wheelY != 0 {
		zoomDelta := wheelY * 0.1 // 0.1 zoom per wheel notch
		g.cameraZoom += zoomDelta
	}

	// Keyboard zoom controls
	if ebiten.IsKeyPressed(ebiten.KeyEqual) || ebiten.IsKeyPressed(ebiten.KeyKPAdd) {
		g.cameraZoom += 0.02 // Smooth zoom in
	}
	if ebiten.IsKeyPressed(ebiten.KeyMinus) || ebiten.IsKeyPressed(ebiten.KeyKPSubtract) {
		g.cameraZoom -= 0.02 // Smooth zoom out
	}

	// Clamp zoom between 0.5x and 10x
	if g.cameraZoom < 0.5 {
		g.cameraZoom = 0.5
	}
	if g.cameraZoom > 10.0 {
		g.cameraZoom = 10.0
	}

	// Pan controls: Arrow keys or WASD
	panSpeed := 5.0 / g.cameraZoom // Pan slower when zoomed in

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.cameraOffsetX += panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.cameraOffsetX -= panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.cameraOffsetY += panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.cameraOffsetY -= panSpeed
	}

	// Reset camera: R key
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.cameraZoom = 1.0
		g.cameraOffsetX = 0
		g.cameraOffsetY = 0
	}
}

// screenToWorld converts screen coordinates to world coordinates (accounting for camera)
func (g *Game) screenToWorld(screenX, screenY float64) (float64, float64) {
	// Center of screen
	centerX := float64(control.DefaultConfig.DisplayScreenWidth) / 2.0
	centerY := float64(control.DefaultConfig.DisplayScreenHeight) / 2.0

	// Transform from screen space to world space
	worldX := (screenX - centerX) / g.cameraZoom - g.cameraOffsetX + centerX
	worldY := (screenY - centerY) / g.cameraZoom - g.cameraOffsetY + centerY

	return worldX, worldY
}

// worldToScreen converts world coordinates to screen coordinates (accounting for camera)
func (g *Game) worldToScreen(worldX, worldY float64) (float64, float64) {
	// Center of screen
	centerX := float64(control.DefaultConfig.DisplayScreenWidth) / 2.0
	centerY := float64(control.DefaultConfig.DisplayScreenHeight) / 2.0

	// Transform from world space to screen space
	screenX := (worldX - centerX + g.cameraOffsetX) * g.cameraZoom + centerX
	screenY := (worldY - centerY + g.cameraOffsetY) * g.cameraZoom + centerY

	return screenX, screenY
}

func (g *Game) handleMouseClick() {
	// Detect mouse button press (on the frame it's pressed)
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	// Only handle click on the frame the button was just pressed
	if mousePressed && !g.lastMouseClick {
		mx, my := ebiten.CursorPosition()
		// Keep mouse position in screen space
		mousePos := models.NewVector(float64(mx), float64(my))

		if g.currentScene == SceneStation {
			// In station scene, check for back button
			if g.isPointInBackButton(mousePos) {
				g.currentScene = SceneMap
				g.selectedStation = nil
			}
		} else if g.currentScene == SceneNewspaper {
			// In newspaper scene, check for back button
			if g.isPointInBackButton(mousePos) {
				g.currentScene = SceneMap
			}
		} else {
			// In map scene
			// Check if clicked on score panel to toggle breakdown
			if g.isPointInScorePanel(mousePos) {
				g.scoreBreakdownOpen = !g.scoreBreakdownOpen
				return
			}

			// Check if clicked on newspaper button
			if g.isPointInNewspaperButton(mousePos) {
				g.currentScene = SceneNewspaper
				return
			}

			// Clear previous train selection
			g.selectedTrain = nil

			// Check if clicked on a train (compare in screen space)
			for i := range g.trains {
				// Convert train position to screen space
				screenX, screenY := g.worldToScreen(g.trains[i].Position.X, g.trains[i].Position.Y)
				screenPos := models.NewVector(screenX, screenY)

				if g.isPointInBounds(mousePos, screenPos, g.trains[i].FrameWidth, g.trains[i].FrameHeight) {
					g.selectedTrain = &g.trains[i]
					break
				}
			}

			// If no train clicked, check stations (switch to station scene)
			if g.selectedTrain == nil {
				for _, st := range g.stations {
					// Convert station position to screen space
					screenX, screenY := g.worldToScreen(st.Position.X, st.Position.Y)
					screenPos := models.NewVector(screenX, screenY)

					if g.isPointInBounds(mousePos, screenPos, st.FrameWidth, st.FrameHeight) {
						g.selectedStation = st
						g.currentScene = SceneStation
						break
					}
				}
			}
		}
	}

	g.lastMouseClick = mousePressed
}

func (g *Game) isPointInBackButton(point models.Vector) bool {
	// Back button in top-left corner: 10, 10, 80x30
	return point.X >= 10 && point.X <= 90 && point.Y >= 10 && point.Y <= 40
}

func (g *Game) isPointInScorePanel(point models.Vector) bool {
	// Score panel in top-left corner: 10, 10, 180x100 (or 180x180 when expanded)
	panelHeight := 100.0
	if g.scoreBreakdownOpen {
		panelHeight = 180.0
	}
	return point.X >= 10 && point.X <= 190 && point.Y >= 10 && point.Y <= 10+panelHeight
}

func (g *Game) isPointInNewspaperButton(point models.Vector) bool {
	// Newspaper button in top-right corner: screenWidth-70, 10, 60x60
	buttonX := float64(control.DefaultConfig.DisplayScreenWidth - 70)
	return point.X >= buttonX && point.X <= buttonX+60 && point.Y >= 10 && point.Y <= 70
}

func (g *Game) isPointInBounds(point, objPos models.Vector, width, height int) bool {
	halfW := float64(width) / 2
	halfH := float64(height) / 2
	return point.X >= objPos.X-halfW && point.X <= objPos.X+halfW &&
		point.Y >= objPos.Y-halfH && point.Y <= objPos.Y+halfH
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.currentScene == SceneStation {
		g.drawStationScene(screen)
	} else if g.currentScene == SceneNewspaper {
		g.drawNewspaperScene(screen)
	} else {
		g.drawMapScene(screen)
	}
}

func (g *Game) drawMapScene(screen *ebiten.Image) {
	// Draw lines first (background layer) - lines draw between stations
	for _, line := range g.lines {
		// Lines need to be drawn with transformed endpoints
		// For now, we'll draw them normally since Line.Draw handles its own rendering
		// TODO: Implement line drawing with camera transform
		g.drawLineTransformed(screen, line)
	}

	// Draw stations with camera transform
	for _, st := range g.stations {
		// Get screen position
		screenX, screenY := g.worldToScreen(st.Position.X, st.Position.Y)

		// Calculate current animation frame
		i := (st.Counter / st.FrameCount) % st.FrameCount
		sx, sy := 0+i*st.FrameWidth, 0
		frame := st.Sprite.SubImage(image.Rect(sx, sy, sx+st.FrameWidth, sy+st.FrameHeight)).(*ebiten.Image)

		// Draw station sprite at constant screen size (no zoom scaling)
		op := &ebiten.DrawImageOptions{}
		// Center the sprite on its screen position
		op.GeoM.Translate(-float64(st.FrameWidth)/2, -float64(st.FrameHeight)/2)
		op.GeoM.Translate(screenX, screenY)

		screen.DrawImage(frame, op)

		// Draw waiting passengers
		g.drawWaitingPassengersTransformed(screen, st)
	}

	// Draw trains with camera transform
	for i := range g.trains {
		tr := &g.trains[i]
		// Get screen position
		screenX, screenY := g.worldToScreen(tr.Position.X, tr.Position.Y)

		// Calculate current animation frame
		frameIdx := (tr.Counter / tr.FrameCount) % tr.FrameCount
		sx, sy := 0+frameIdx*tr.FrameWidth, 0
		frame := tr.Sprite.SubImage(image.Rect(sx, sy, sx+tr.FrameWidth, sy+tr.FrameHeight)).(*ebiten.Image)

		// Draw train sprite at constant screen size (no zoom scaling)
		op := &ebiten.DrawImageOptions{}
		// Center the sprite on its screen position
		op.GeoM.Translate(-float64(tr.FrameWidth)/2, -float64(tr.FrameHeight)/2)
		op.GeoM.Translate(screenX, screenY)

		screen.DrawImage(frame, op)
	}

	// Draw text labels in screen space (crisp text regardless of zoom)
	for _, st := range g.stations {
		// Transform station position to screen space
		screenX, screenY := g.worldToScreen(st.Position.X, st.Position.Y)
		screenPos := models.NewVector(screenX, screenY)

		// Draw station name
		DrawTitle(screen, st.Name, screenPos, XS_FONT_SIZE, st.FrameWidth, st.FrameHeight, TITLE_BOT_SIDE)

		// Draw passenger count
		waitingCount := st.GetWaitingPassengersCount()
		if waitingCount > 0 {
			countText := fmt.Sprintf("%d", waitingCount)
			DrawInfo(screen, countText, screenPos, S_FONT_SIZE, st.FrameWidth, st.FrameHeight)
		}
	}

	// Draw train labels in screen space
	for i := range g.trains {
		tr := &g.trains[i]
		// Transform train position to screen space
		screenX, screenY := g.worldToScreen(tr.Position.X, tr.Position.Y)
		screenPos := models.NewVector(screenX, screenY)

		// Draw train name
		DrawTitle(screen, tr.Name, screenPos, S_FONT_SIZE, tr.FrameWidth, tr.FrameHeight, TITLE_TOP_SIDE)

		// Draw passenger count
		passengerCount := tr.GetPassengerCount()
		if passengerCount > 0 {
			countText := fmt.Sprintf("%d", passengerCount)
			DrawInfo(screen, countText, screenPos, S_FONT_SIZE, tr.FrameWidth, tr.FrameHeight)
		}
	}

	// Draw UI elements on top (unaffected by camera)
	// Draw train data panel if train is selected
	if g.selectedTrain != nil {
		g.drawTrainDataPanel(screen)
	}

	// Draw score overlay (always visible)
	g.drawScoreOverlay(screen)

	// Draw newspaper button (top-right corner)
	g.drawNewspaperButton(screen)

	// Draw simulation clock (top-center)
	g.drawSimulationClock(screen)

	// Draw camera controls help (bottom-right)
	g.drawCameraHelp(screen)
}

// drawLineTransformed draws a line with camera transform applied
func (g *Game) drawLineTransformed(screen *ebiten.Image, line models.Line) {
	// Get all stations for this line (already in correct order from DB via station_line.odr)
	stations := line.Stations

	if len(stations) < 2 {
		return
	}

	// Use a dark gray color that blends with the black background
	lineColor := color.RGBA{60, 60, 60, 180} // Very dark gray, semi-transparent

	// Draw very thin dashed line segments connecting consecutive stations
	for i := 0; i < len(stations)-1; i++ {
		st1 := stations[i]
		st2 := stations[i+1]

		// Transform both endpoints to screen space
		x1, y1 := g.worldToScreen(st1.Position.X, st1.Position.Y)
		x2, y2 := g.worldToScreen(st2.Position.X, st2.Position.Y)

		// Draw as dashed line (draw short segments with gaps)
		g.drawDashedLine(screen, x1, y1, x2, y2, lineColor)
	}
}

// drawDashedLine draws a dashed line between two points
func (g *Game) drawDashedLine(screen *ebiten.Image, x1, y1, x2, y2 float64, lineColor color.Color) {
	// Calculate line length
	dx := x2 - x1
	dy := y2 - y1
	length := math.Sqrt(dx*dx + dy*dy)

	if length < 1 {
		return
	}

	// Normalize direction
	dirX := dx / length
	dirY := dy / length

	// Dash pattern: 5px dash, 3px gap
	dashLength := 5.0
	gapLength := 3.0
	patternLength := dashLength + gapLength

	// Draw dashes along the line
	currentPos := 0.0
	for currentPos < length {
		// Start of dash
		startX := x1 + dirX*currentPos
		startY := y1 + dirY*currentPos

		// End of dash (but not beyond line end)
		endPos := currentPos + dashLength
		if endPos > length {
			endPos = length
		}
		endX := x1 + dirX*endPos
		endY := y1 + dirY*endPos

		// Draw this dash segment (very thin, constant 1px width)
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY),
			1.0, lineColor, true)

		// Move to next dash
		currentPos += patternLength
	}
}

// drawWaitingPassengersTransformed draws passenger dots with camera transform
func (g *Game) drawWaitingPassengersTransformed(screen *ebiten.Image, st *models.Station) {
	passengers := st.GetWaitingPassengers()
	if len(passengers) == 0 {
		return
	}

	// Draw max 10 passenger dots to avoid clutter
	maxDots := 10
	if len(passengers) < maxDots {
		maxDots = len(passengers)
	}

	// Arrange dots in a circle around the station (constant screen-space radius)
	radius := 20.0 // Screen pixels
	stationScreenX, stationScreenY := g.worldToScreen(st.Position.X, st.Position.Y)

	for i := 0; i < maxDots; i++ {
		angle := float64(i) * (2.0 * math.Pi / float64(maxDots))
		screenX := stationScreenX + radius*math.Cos(angle)
		screenY := stationScreenY + radius*math.Sin(angle)

		// Draw passenger as a small colored circle (constant screen size)
		passengerColor := getPassengerColor(passengers[i].Sentiment)
		vector.DrawFilledCircle(screen, float32(screenX), float32(screenY), 2.0, passengerColor, true)
	}
}

func (g *Game) drawWaitingPassengers(screen *ebiten.Image, st *models.Station) {
	passengers := st.GetWaitingPassengers()
	if len(passengers) == 0 {
		return
	}

	// Draw max 10 passenger dots to avoid clutter
	maxDots := 10
	if len(passengers) < maxDots {
		maxDots = len(passengers)
	}

	// Arrange dots in a circle around the station
	radius := float32(20.0)
	for i := 0; i < maxDots; i++ {
		angle := float64(i) * (2.0 * math.Pi / float64(maxDots))
		x := st.Position.X + float64(radius)*math.Cos(angle)
		y := st.Position.Y + float64(radius)*math.Sin(angle)

		// Draw passenger as a small colored circle
		passengerColor := getPassengerColor(passengers[i].Sentiment)
		vector.DrawFilledCircle(screen, float32(x), float32(y), 2, passengerColor, true)
	}
}

func getPassengerColor(sentiment float64) color.Color {
	// Green (happy) -> Yellow (neutral) -> Red (angry)
	if sentiment >= 70 {
		return color.RGBA{0, 255, 0, 255} // Green
	} else if sentiment >= 40 {
		return color.RGBA{255, 255, 0, 255} // Yellow
	} else {
		return color.RGBA{255, 0, 0, 255} // Red
	}
}

func (g *Game) drawTrainDataPanel(screen *ebiten.Image) {
	tr := g.selectedTrain

	// Panel dimensions and position (top-right corner)
	panelX := float32(control.DefaultConfig.DisplayScreenWidth - 200)
	panelY := float32(10)
	panelW := float32(190)
	panelH := float32(160) // Increased height for schedule info

	// Draw panel background
	vector.DrawFilledRect(screen, panelX, panelY, panelW, panelH, color.RGBA{30, 30, 40, 230}, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 2, color.RGBA{100, 150, 200, 255}, false)

	// Draw train data
	yPos := float32(25.0)
	DrawDataText(screen, fmt.Sprintf("TRAIN: %s", tr.Name), panelX+10, yPos, M_FONT_SIZE)
	yPos += 20
	DrawDataText(screen, fmt.Sprintf("Speed: %.1f km/h", tr.GetSpeedKmH()), panelX+10, yPos, S_FONT_SIZE)
	yPos += 15
	DrawDataText(screen, fmt.Sprintf("Passengers: %d/%d", tr.GetPassengerCount(), tr.Capacity), panelX+10, yPos, S_FONT_SIZE)
	yPos += 15

	// Show destination info
	if tr.Next != nil {
		DrawDataText(screen, fmt.Sprintf("Next: %s", tr.Next.Name), panelX+10, yPos, S_FONT_SIZE)
		yPos += 15

		// Show scheduled arrival time if available
		if g.scheduleDB != nil && g.clock != nil {
			schedule, err := g.scheduleDB.GetScheduleByTrainAndStation(tr.ID, tr.Next.ID)
			if err == nil {
				scheduledTime := g.formatTime(schedule.ScheduledTime)
				currentTime := g.clock.GetCurrentTimeOfDay()

				// Calculate ETA
				timeDiff := schedule.ScheduledTime - currentTime
				etaText := ""
				if timeDiff > 0 {
					mins := timeDiff / 60
					etaText = fmt.Sprintf(" (ETA: %dm)", mins)
				}

				DrawDataText(screen, fmt.Sprintf("Sched: %s%s", scheduledTime, etaText), panelX+10, yPos, XS_FONT_SIZE)
			}
		}
	} else if tr.Current != nil {
		DrawDataText(screen, fmt.Sprintf("At: %s", tr.Current.Name), panelX+10, yPos, S_FONT_SIZE)
	}
	yPos += 15

	// Show punctuality status if brain is available
	if g.brain != nil {
		metrics := g.brain.GetMetrics()
		if metrics.TotalArrivalsChecked > 0 {
			statusText := fmt.Sprintf("On-Time: %.0f%%", metrics.OnTimePercentage)
			statusColor := color.RGBA{0, 200, 0, 255} // Green

			if metrics.OnTimePercentage < 70 {
				statusColor = color.RGBA{200, 0, 0, 255} // Red
			} else if metrics.OnTimePercentage < 85 {
				statusColor = color.RGBA{200, 200, 0, 255} // Yellow
			}

			DrawColoredText(screen, statusText, panelX+10, yPos, XS_FONT_SIZE, statusColor)
			yPos += 15
		}
	}

	// Capacity bar
	capacityPct := tr.GetCapacityPercentage() / 100.0
	barW := float32(170)
	barH := float32(10)
	vector.DrawFilledRect(screen, panelX+10, float32(yPos), barW, barH, color.RGBA{50, 50, 50, 255}, false)

	barColor := color.RGBA{0, 200, 0, 255}
	if capacityPct > 0.8 {
		barColor = color.RGBA{200, 0, 0, 255}
	} else if capacityPct > 0.5 {
		barColor = color.RGBA{200, 200, 0, 255}
	}
	vector.DrawFilledRect(screen, panelX+10, float32(yPos), barW*float32(capacityPct), barH, barColor, false)
}

func (g *Game) drawStationScene(screen *ebiten.Image) {
	if g.selectedStation == nil {
		return
	}
	st := g.selectedStation

	// Draw background (platform)
	bgColor := color.RGBA{40, 40, 50, 255}
	screen.Fill(bgColor)

	// Draw platform floor
	floorY := float32(control.DefaultConfig.DisplayScreenHeight - 100)
	vector.DrawFilledRect(screen, 0, floorY, float32(control.DefaultConfig.DisplayScreenWidth), 100, color.RGBA{60, 60, 70, 255}, false)

	// Draw back button (top-left)
	g.drawBackButton(screen)

	// Draw station name (centered at top)
	stationNameY := float32(30)
	DrawDataText(screen, st.Name, float32(control.DefaultConfig.DisplayScreenWidth/2-50), stationNameY, L_FONT_SIZE)

	// Draw waiting count
	waitingCount := st.GetWaitingPassengersCount()
	DrawDataText(screen, fmt.Sprintf("Waiting Passengers: %d", waitingCount), 20, 80, M_FONT_SIZE)

	// Draw passengers as sprites
	passengers := st.GetWaitingPassengers()
	if len(passengers) > 0 {
		// Arrange passengers in rows on the platform
		spacing := 60.0
		perRow := 10
		startX := 100.0
		startY := float64(floorY) - 60 // Position passengers just above the floor

		for i, p := range passengers {
			row := i / perRow
			col := i % perRow
			x := startX + float64(col)*spacing
			y := startY + float64(row)*spacing

			// Draw passenger as colored circle (by sentiment)
			sentimentColor := getPassengerColor(p.Sentiment)
			vector.DrawFilledCircle(screen, float32(x), float32(y), 8, sentimentColor, true)

			// Draw passenger outline
			vector.StrokeCircle(screen, float32(x), float32(y), 8, 1, color.White, true)

			// Draw passenger name below
			DrawDataText(screen, p.Name, float32(x-20), float32(y+15), XS_FONT_SIZE)
		}
	}
}

func (g *Game) drawBackButton(screen *ebiten.Image) {
	// Button background
	buttonX := float32(10)
	buttonY := float32(10)
	buttonW := float32(80)
	buttonH := float32(30)

	vector.DrawFilledRect(screen, buttonX, buttonY, buttonW, buttonH, color.RGBA{100, 100, 120, 255}, false)
	vector.StrokeRect(screen, buttonX, buttonY, buttonW, buttonH, 2, color.RGBA{150, 150, 170, 255}, false)

	// Button text
	DrawDataText(screen, "< Back", buttonX+10, buttonY+20, S_FONT_SIZE)
}

// drawScoreOverlay draws the system score in the top-left corner
func (g *Game) drawScoreOverlay(screen *ebiten.Image) {
	// Skip if Tenjin is not enabled or brain is nil
	if g.brain == nil {
		return
	}

	// Get current metrics with score
	metrics := g.brain.GetMetrics()
	score := metrics.Score

	// Panel dimensions and position (top-left corner)
	panelX := float32(10)
	panelY := float32(10)
	panelW := float32(180)
	panelH := float32(100)

	// Expand panel if breakdown is open
	if g.scoreBreakdownOpen {
		panelH = float32(180)
	}

	// Draw panel background
	vector.DrawFilledRect(screen, panelX, panelY, panelW, panelH, color.RGBA{30, 30, 40, 230}, false)

	// Get grade color
	gradeColor := g.getGradeColorRGBA(score.Grade)

	// Draw border with grade color
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 3, gradeColor, false)

	// Draw score title
	yPos := float32(25.0)
	DrawDataText(screen, "SYSTEM SCORE", panelX+10, yPos, S_FONT_SIZE)
	yPos += 20

	// Draw overall score with grade color
	scoreText := fmt.Sprintf("%.1f", score.Overall)
	DrawDataText(screen, scoreText, panelX+10, yPos, L_FONT_SIZE)

	// Draw grade next to score
	gradeX := panelX + 130
	DrawDataText(screen, score.Grade, gradeX, yPos-30, M_FONT_SIZE)
	yPos += 25

	if g.scoreBreakdownOpen {
		// Draw detailed component scores
		DrawDataText(screen, "--- Components ---", panelX+10, yPos, XS_FONT_SIZE)
		yPos += 15

		DrawDataText(screen, fmt.Sprintf("Satisfaction: %.1f", score.PassengerSatisfaction), panelX+10, yPos, S_FONT_SIZE)
		yPos += 15

		DrawDataText(screen, fmt.Sprintf("Efficiency: %.1f", score.ServiceEfficiency), panelX+10, yPos, S_FONT_SIZE)
		yPos += 15

		DrawDataText(screen, fmt.Sprintf("Capacity: %.1f", score.SystemCapacity), panelX+10, yPos, S_FONT_SIZE)
		yPos += 15

		DrawDataText(screen, fmt.Sprintf("Reliability: %.1f", score.Reliability), panelX+10, yPos, S_FONT_SIZE)
		yPos += 15

		// Draw hint
		DrawDataText(screen, "(click to collapse)", panelX+10, yPos, XS_FONT_SIZE)
	} else {
		// Draw component scores (compact)
		componentText := fmt.Sprintf("S:%.0f E:%.0f C:%.0f R:%.0f",
			score.PassengerSatisfaction,
			score.ServiceEfficiency,
			score.SystemCapacity,
			score.Reliability)
		DrawDataText(screen, componentText, panelX+10, yPos, XS_FONT_SIZE)
		yPos += 12

		// Draw hint
		DrawDataText(screen, "(click for details)", panelX+10, yPos, XS_FONT_SIZE)
	}
}

// getGradeColorRGBA returns the color for a grade
func (g *Game) getGradeColorRGBA(grade string) color.RGBA {
	switch grade {
	case "S":
		return color.RGBA{255, 215, 0, 255} // Gold
	case "A":
		return color.RGBA{0, 255, 0, 255} // Green
	case "B":
		return color.RGBA{144, 238, 144, 255} // Light Green
	case "C":
		return color.RGBA{255, 255, 0, 255} // Yellow
	case "D":
		return color.RGBA{255, 165, 0, 255} // Orange
	case "F":
		return color.RGBA{255, 0, 0, 255} // Red
	default:
		return color.RGBA{255, 255, 255, 255} // White
	}
}

// drawNewspaperButton draws the newspaper button in the top-right corner
func (g *Game) drawNewspaperButton(screen *ebiten.Image) {
	buttonX := float32(control.DefaultConfig.DisplayScreenWidth - 70)
	buttonY := float32(10)
	buttonW := float32(60)
	buttonH := float32(60)

	// Draw button background
	vector.DrawFilledRect(screen, buttonX, buttonY, buttonW, buttonH, color.RGBA{100, 100, 120, 255}, false)
	vector.StrokeRect(screen, buttonX, buttonY, buttonW, buttonH, 2, color.RGBA{200, 200, 220, 255}, false)

	// Draw newspaper icon (📰)
	DrawDataText(screen, "NEWS", buttonX+5, buttonY+35, S_FONT_SIZE)
}

// drawNewspaperScene draws the full newspaper view
func (g *Game) drawNewspaperScene(screen *ebiten.Image) {
	// Draw background (newspaper-like)
	bgColor := color.RGBA{240, 235, 220, 255} // Cream/beige newspaper color
	screen.Fill(bgColor)

	// Draw back button (top-left)
	g.drawBackButton(screen)

	// Get current newspaper edition
	newspaper := g.brain.GetNewspaper()
	edition := newspaper.GetCurrentEdition()

	// Title bar (darker background)
	titleBarHeight := float32(100)
	vector.DrawFilledRect(screen, 0, 0, float32(control.DefaultConfig.DisplayScreenWidth), titleBarHeight, color.RGBA{30, 30, 40, 255}, false)

	// Newspaper title (white on dark background)
	titleX := float32(control.DefaultConfig.DisplayScreenWidth/2 - 180)
	DrawDataText(screen, "METRO DAILY NEWS", titleX, 55, XXL_FONT_SIZE)

	// Check if newspaper is available
	if edition == nil {
		// Show "generating" message (black text on cream background)
		blackColor := color.RGBA{0, 0, 0, 255}
		if newspaper.IsGenerating() {
			message := "Generating today's edition..."
			messageX := float32(control.DefaultConfig.DisplayScreenWidth/2 - 150)
			DrawColoredText(screen, message, messageX, 200, L_FONT_SIZE, blackColor)
		} else {
			message := "No edition available yet."
			messageX := float32(control.DefaultConfig.DisplayScreenWidth/2 - 120)
			DrawColoredText(screen, message, messageX, 200, L_FONT_SIZE, blackColor)
		}
		return
	}

	// Display date (white on dark title bar)
	dateText := edition.Date.Format("Monday, January 2, 2006")
	dateX := float32(control.DefaultConfig.DisplayScreenWidth/2 - 120)
	DrawDataText(screen, dateText, dateX, 85, M_FONT_SIZE)

	// Draw stories in columns
	startY := float32(120)
	currentY := startY
	margin := float32(40)
	storySpacing := float32(50)
	blackColor := color.RGBA{0, 0, 0, 255}
	grayColor := color.RGBA{50, 50, 50, 255}

	for i, story := range edition.Stories {
		// Story divider
		if i > 0 {
			dividerY := currentY - 20
			vector.DrawFilledRect(screen, margin, dividerY, float32(control.DefaultConfig.DisplayScreenWidth)-margin*2, 2, color.RGBA{100, 100, 100, 255}, false)
			currentY += 10
		}

		// Story headline (bold, larger, black text)
		vector.DrawFilledRect(screen, margin-10, currentY-10, float32(control.DefaultConfig.DisplayScreenWidth)-margin*2+20, 50, color.RGBA{255, 250, 230, 255}, false)
		DrawColoredText(screen, fmt.Sprintf("📰 %s", story.Headline), margin, currentY+20, M_FONT_SIZE, blackColor)
		currentY += 60

		// Story article text (wrapped, gray text for body)
		articleLines := g.wrapText(story.Article, 80) // Wrap at ~80 characters for better readability
		for _, line := range articleLines {
			DrawColoredText(screen, line, margin+10, currentY, S_FONT_SIZE, grayColor)
			currentY += 20
		}

		currentY += storySpacing

		// Don't draw too many stories if they don't fit
		if currentY > float32(control.DefaultConfig.DisplayScreenHeight-80) {
			break
		}
	}
}

// wrapText wraps text to fit within a certain character width
func (g *Game) wrapText(text string, width int) []string {
	words := []string{}
	words = append(words, splitBySpaces(text)...)

	lines := []string{}
	currentLine := ""

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) > width {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		} else {
			currentLine = testLine
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// splitBySpaces splits a string by spaces
func splitBySpaces(s string) []string {
	result := []string{}
	current := ""
	for _, ch := range s {
		if ch == ' ' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// drawSimulationClock draws the current simulation time in the top-center
func (g *Game) drawSimulationClock(screen *ebiten.Image) {
	if g.clock == nil {
		return
	}

	// Clock panel in top-center
	panelW := float32(120)
	panelH := float32(40)
	panelX := float32(control.DefaultConfig.DisplayScreenWidth/2) - panelW/2
	panelY := float32(10)

	// Draw panel background
	vector.DrawFilledRect(screen, panelX, panelY, panelW, panelH, color.RGBA{30, 30, 40, 230}, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 2, color.RGBA{100, 150, 200, 255}, false)

	// Draw time
	currentTime := g.clock.GetCurrentTime()
	DrawDataText(screen, currentTime, panelX+15, panelY+25, L_FONT_SIZE)
}

// drawCameraHelp draws the camera control instructions in the bottom-right
func (g *Game) drawCameraHelp(screen *ebiten.Image) {
	// Panel in bottom-right corner
	panelW := float32(180)
	panelH := float32(90)
	panelX := float32(control.DefaultConfig.DisplayScreenWidth) - panelW - 10
	panelY := float32(control.DefaultConfig.DisplayScreenHeight) - panelH - 10

	// Draw panel background (semi-transparent)
	vector.DrawFilledRect(screen, panelX, panelY, panelW, panelH, color.RGBA{30, 30, 40, 200}, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 1, color.RGBA{100, 150, 200, 255}, false)

	// Draw controls text
	textX := panelX + 10
	textY := panelY + 15
	lineHeight := float32(14)

	DrawDataText(screen, "Camera Controls:", textX, textY, XS_FONT_SIZE)
	textY += lineHeight
	DrawDataText(screen, "Wheel/+/-: Zoom", textX, textY, XS_FONT_SIZE)
	textY += lineHeight
	DrawDataText(screen, "Arrows/WASD: Pan", textX, textY, XS_FONT_SIZE)
	textY += lineHeight
	DrawDataText(screen, "R: Reset", textX, textY, XS_FONT_SIZE)
	textY += lineHeight
	DrawDataText(screen, fmt.Sprintf("Zoom: %.1fx", g.cameraZoom), textX, textY, XS_FONT_SIZE)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return control.DefaultConfig.DisplayScreenWidth, control.DefaultConfig.DisplayScreenHeight
}
