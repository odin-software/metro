package display

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/odin-software/metro/control"
	"github.com/odin-software/metro/internal/models"
)

type Game struct {
	trains   []models.Train
	stations []models.Station
}

func NewGame(trains []models.Train, stations []models.Station) *Game {
	Init()
	return &Game{
		trains:   trains,
		stations: stations,
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
		DrawText(screen, tr.Name, tr.Position, S_FONT_SIZE)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return control.DefaultConfig.DisplayScreenWidth, control.DefaultConfig.DisplayScreenHeight
}
