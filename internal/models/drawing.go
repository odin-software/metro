package models

import "github.com/hajimehoshi/ebiten/v2"

type Drawing struct {
	Counter     int
	FrameWidth  int
	FrameHeight int
	FrameCount  int
	Sprite      *ebiten.Image
}
