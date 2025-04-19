package display

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/odin-software/metro/internal/assets"
	"github.com/odin-software/metro/internal/models"
)

var textFontSource *text.GoTextFaceSource

type FontSize int

const (
	XS_FONT_SIZE FontSize = 6
	S_FONT_SIZE  FontSize = 12
	M_FONT_SIZE  FontSize = 24
	L_FONT_SIZE  FontSize = 48
)

func Init() {
	data, err := assets.Assets.ReadFile("files/fonts/future.ttf")
	if err != nil {
		log.Fatal(err)
	}
	s, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	textFontSource = s
}

func DrawText(screen *ebiten.Image, info string, position models.Vector, size FontSize) {
	// TODO: Make color configurable.
	op := &text.DrawOptions{}
	op.GeoM.Translate(position.X-14, position.Y-15)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, info, &text.GoTextFace{
		Source: textFontSource,
		Size:   float64(size),
	}, op)
}
