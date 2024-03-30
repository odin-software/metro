package baso

import (
	"log"

	models "github.com/odin-software/metro/internal/models"
)

func (bs *Baso) ListMakes() []models.Make {
	makes, err := bs.queries.ListMakes(bs.ctx)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]models.Make, 0)
	for _, make := range makes {
		result = append(
			result,
			models.NewMake(make.Name, make.Description, make.Acceleration.Float64, make.TopSpeed.Float64),
		)
	}
	return result
}
