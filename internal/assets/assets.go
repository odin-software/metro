package assets

import (
	"bytes"
	"embed"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed files/*
var Assets embed.FS

func GetTrainSprite() (ebitenImage *ebiten.Image, frameWidth int, frameHeight int, frameCount int) {
	t, err := Assets.ReadFile("files/t2.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(bytes.NewReader(t))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img), 13, 13, 13
}
