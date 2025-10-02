package models

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Line struct {
	ID       int64
	Name     string
	Stations []*Station // Pointers to share state across system
}

func LineInit() {
	whiteImage.Fill(color.White)
}

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func (ln Line) Draw(screen *ebiten.Image) {
	var path vector.Path
	path.MoveTo(float32(ln.Stations[0].Position.X), float32(ln.Stations[0].Position.Y))
	for _, st := range ln.Stations {
		path.LineTo(float32(st.Position.X), float32(st.Position.Y))
	}

	// Draw the main line in white.
	op := &vector.StrokeOptions{}
	op.LineCap = vector.LineCapRound
	op.Width = float32(10 / 2)

	var verts []ebiten.Vertex
	var idxs []uint16

	vs, is := path.AppendVerticesAndIndicesForStroke(verts, idxs[:0], op)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = 0.4
		vs[i].ColorG = 1
		vs[i].ColorB = 1
		vs[i].ColorA = 1
	}
	screen.DrawTriangles(vs, is, whiteSubImage, &ebiten.DrawTrianglesOptions{
		AntiAlias: true,
	})
}
