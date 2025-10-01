package display

import (
	"fmt"
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
}

func NewGame(trains []models.Train, stations []*models.Station, lines []models.Line, brain *tenjin.Tenjin) *Game {
	Init()
	models.LineInit()
	return &Game{
		trains:       trains,
		stations:     stations,
		lines:        lines,
		currentScene: SceneMap,
		brain:        brain,
	}
}

func (g *Game) Update() error {
	for i := range g.trains {
		g.trains[i].Update()
	}
	for i := range g.stations {
		g.stations[i].Update()
	}

	// Handle mouse clicks
	g.handleMouseClick()

	return nil
}

func (g *Game) handleMouseClick() {
	// Detect mouse button press (on the frame it's pressed)
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	// Only handle click on the frame the button was just pressed
	if mousePressed && !g.lastMouseClick {
		mx, my := ebiten.CursorPosition()
		mousePos := models.NewVector(float64(mx), float64(my))

		if g.currentScene == SceneStation {
			// In station scene, check for back button
			if g.isPointInBackButton(mousePos) {
				g.currentScene = SceneMap
				g.selectedStation = nil
			}
		} else {
			// In map scene
			// Check if clicked on score panel to toggle breakdown
			if g.isPointInScorePanel(mousePos) {
				g.scoreBreakdownOpen = !g.scoreBreakdownOpen
				return
			}

			// Clear previous train selection
			g.selectedTrain = nil

			// Check if clicked on a train
			for i := range g.trains {
				if g.isPointInBounds(mousePos, g.trains[i].Position, g.trains[i].FrameWidth, g.trains[i].FrameHeight) {
					g.selectedTrain = &g.trains[i]
					break
				}
			}

			// If no train clicked, check stations (switch to station scene)
			if g.selectedTrain == nil {
				for _, st := range g.stations {
					if g.isPointInBounds(mousePos, st.Position, st.FrameWidth, st.FrameHeight) {
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

func (g *Game) isPointInBounds(point, objPos models.Vector, width, height int) bool {
	halfW := float64(width) / 2
	halfH := float64(height) / 2
	return point.X >= objPos.X-halfW && point.X <= objPos.X+halfW &&
		point.Y >= objPos.Y-halfH && point.Y <= objPos.Y+halfH
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.currentScene == SceneStation {
		g.drawStationScene(screen)
	} else {
		g.drawMapScene(screen)
	}
}

func (g *Game) drawMapScene(screen *ebiten.Image) {
	// Draw stations
	for _, st := range g.stations {
		st.Draw(screen)
		DrawTitle(screen, st.Name, st.Position, XS_FONT_SIZE, st.FrameWidth, st.FrameHeight, TITLE_BOT_SIDE)

		// Draw waiting passengers as small dots around the station
		g.drawWaitingPassengers(screen, st)

		// Draw passenger count
		waitingCount := st.GetWaitingPassengersCount()
		if waitingCount > 0 {
			countText := fmt.Sprintf("%d", waitingCount)
			DrawInfo(screen, countText, st.Position, S_FONT_SIZE, st.FrameWidth, st.FrameHeight)
		}
	}

	// Draw trains
	for _, tr := range g.trains {
		tr.Draw(screen)
		DrawTitle(screen, tr.Name, tr.Position, S_FONT_SIZE, tr.FrameWidth, tr.FrameHeight, TITLE_TOP_SIDE)

		// Draw passenger count on train
		passengerCount := tr.GetPassengerCount()
		if passengerCount > 0 {
			countText := fmt.Sprintf("%d", passengerCount)
			DrawInfo(screen, countText, tr.Position, S_FONT_SIZE, tr.FrameWidth, tr.FrameHeight)
		}
	}

	// Draw train data panel if train is selected
	if g.selectedTrain != nil {
		g.drawTrainDataPanel(screen)
	}

	// Draw score overlay (always visible)
	g.drawScoreOverlay(screen)
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
	panelH := float32(120)

	// Draw panel background
	vector.DrawFilledRect(screen, panelX, panelY, panelW, panelH, color.RGBA{30, 30, 40, 230}, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 2, color.RGBA{100, 150, 200, 255}, false)

	// Draw train data
	yPos := float32(25.0)
	DrawDataText(screen, fmt.Sprintf("TRAIN: %s", tr.Name), panelX+10, yPos, M_FONT_SIZE)
	yPos += 20
	DrawDataText(screen, fmt.Sprintf("Speed: %.2f", tr.GetSpeed()), panelX+10, yPos, S_FONT_SIZE)
	yPos += 15
	DrawDataText(screen, fmt.Sprintf("Passengers: %d/%d", tr.GetPassengerCount(), tr.Capacity), panelX+10, yPos, S_FONT_SIZE)
	yPos += 15

	// Show destination info
	if tr.Next != nil {
		DrawDataText(screen, fmt.Sprintf("Next: %s", tr.Next.Name), panelX+10, yPos, S_FONT_SIZE)
	} else if tr.Current != nil {
		DrawDataText(screen, fmt.Sprintf("At: %s", tr.Current.Name), panelX+10, yPos, S_FONT_SIZE)
	}
	yPos += 15

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

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return control.DefaultConfig.DisplayScreenWidth, control.DefaultConfig.DisplayScreenHeight
}
