package display

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/odin-software/metro/control"
	"github.com/odin-software/metro/internal/models"
)

type Game struct {
	trains   []models.Train
	stations []*models.Station
	lines    []models.Line
}

func NewGame(trains []models.Train, stations []*models.Station, lines []models.Line) *Game {
	Init()
	models.LineInit()
	return &Game{
		trains:   trains,
		stations: stations,
		lines:    lines,
	}
}

func (g *Game) Update() error {
	for i := range g.trains {
		g.trains[i].Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, tr := range g.trains {
		tr.Draw(screen)
		DrawTitle(screen, tr.Name, tr.Position, S_FONT_SIZE, tr.FrameWidth, tr.FrameHeight, TITLE_TOP_SIDE)
	}
	for _, st := range g.stations {
		st.Draw(screen)
		DrawTitle(screen, st.Name, st.Position, XS_FONT_SIZE, st.FrameWidth, st.FrameHeight, TITLE_BOT_SIDE)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return control.DefaultConfig.DisplayScreenWidth, control.DefaultConfig.DisplayScreenHeight
}
