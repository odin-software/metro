package models

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/odin-software/metro/internal/assets"
)

type Station struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Position Vector `json:"position"`
	Drawing
}

func NewStation(id int64, name string, location Vector) *Station {
	img, frameWidth, frameHeight, frameCount := assets.GetStationSprite()
	return &Station{
		ID:       id,
		Name:     name,
		Position: location,
		Drawing: Drawing{
			Counter:     0,
			FrameWidth:  frameWidth,
			FrameHeight: frameHeight,
			FrameCount:  frameCount,
			Sprite:      img,
		},
	}
}

// Drawing methods

func (st *Station) Update() {
	st.Drawing.Counter++
}

func (st *Station) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(st.FrameWidth)/2, -float64(st.FrameHeight)/2)
	op.GeoM.Translate(st.Position.X, st.Position.Y)
	i := (st.Counter / st.FrameCount) % st.FrameCount
	sx, sy := 0+i*st.FrameWidth, 0
	screen.DrawImage(st.Sprite.SubImage(image.Rect(sx, sy, sx+st.FrameWidth, sy+st.FrameHeight)).(*ebiten.Image), op)
}
