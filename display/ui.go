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

var TextFontSource *text.GoTextFaceSource

type FontSize int
type TitleSide int

const (
	XS_FONT_SIZE FontSize = 6
	S_FONT_SIZE  FontSize = 10
	M_FONT_SIZE  FontSize = 24
	L_FONT_SIZE  FontSize = 48

	TITLE_TOP_SIDE   TitleSide = 0
	TITLE_BOT_SIDE   TitleSide = 1
	TITLE_LEFT_SIDE  TitleSide = 2
	TITLE_RIGHT_SIDE TitleSide = 3

	TitleSeparation float64 = 5
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
	TextFontSource = s
}

func DrawText(screen *ebiten.Image, info string, position models.Vector, size FontSize) {
	// TODO: Make color configurable.
	op := &text.DrawOptions{}
	op.GeoM.Translate(position.X-14, position.Y-15)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, info, &text.GoTextFace{
		Source: TextFontSource,
		Size:   float64(size),
	}, op)
}

func DrawTitle(screen *ebiten.Image, info string, position models.Vector, size FontSize, fw int, fh int, side TitleSide) {
	// TODO: Make color configurable.
	f := &text.GoTextFace{
		Source: TextFontSource,
		Size:   float64(size),
	}
	w, h := text.Measure(info, f, 0)
	var xOff, yOff float64
	switch side {
	case TITLE_TOP_SIDE:
		yOff = -h - TitleSeparation
	case TITLE_BOT_SIDE:
		yOff = h + TitleSeparation
	case TITLE_LEFT_SIDE:
		xOff = -w - TitleSeparation
	case TITLE_RIGHT_SIDE:
		xOff = w + TitleSeparation
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(position.X-(w/2)+xOff, position.Y-(h/2)+yOff)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, info, f, op)
}
