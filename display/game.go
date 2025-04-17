package display

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/odin-software/metro/control"
	"github.com/odin-software/metro/internal/models"
)

type Game struct {
	trains []models.Train
}

func NewGame(trains []models.Train) *Game {
	return &Game{
		trains: trains,
	}
}

func (g *Game) Update() error {
	for _, tr := range g.trains {
		tr.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, tr := range g.trains {
		tr.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return control.DefaultConfig.DisplayScreenWidth, control.DefaultConfig.DisplayScreenHeight
}
