package city

import (
	"math"
	"math/rand/v2"

	"github.com/odin-software/metro/internal/baso"
	"github.com/odin-software/metro/internal/models"
)

func getNearest(position models.Vector, stations []baso.GetStation, radius float64) []baso.GetStation {
	near := make([]baso.GetStation, 0)
	for _, st := range stations {
		if st.Position.Dist(position) < radius {
			near = append(near, st)
		}
	}

	return near
}

func poissonDiskSampling(radius float64, width, height, k int) []models.Vector {
	s1 := rand.NewPCG(13, 20)
	r1 := rand.New(s1)

	N := 2.0
	points := make([]models.Vector, 0)
	active := make([]models.Vector, 0)

	p0 := models.NewVector(float64(r1.IntN(width)), float64(rand.IntN(height)))

	cellSize := math.Floor(radius / math.Sqrt(N))

	gridWidth := math.Ceil(float64(width)/cellSize) + 1
	gridHeight := math.Ceil(float64(height)/cellSize) + 1

	grid := make([][]models.Vector, int(gridWidth))
	for i := 0; i < int(gridWidth); i++ {
		grid[i] = make([]models.Vector, int(gridHeight))
		for j := 0; j < int(gridHeight); j++ {
			grid[i][j] = models.NewVector(float64(-1), float64(-1))
		}
	}

	insertPoint(grid, p0, cellSize)
	points = append(points, p0)
	active = append(active, p0)

	for len(active) > 0 {
		ridx := r1.IntN(len(active))
		p := active[ridx]

		found := false
		for tries := 0; tries < k; tries++ {
			// step 1
			theta := (r1.Float64() * 360) * (math.Pi / 180)
			newRadius := radius + (r1.Float64() * ((radius * 2) - radius))
			newX := p.X + newRadius*math.Cos(theta)
			newY := p.Y + newRadius*math.Sin(theta)
			newPoint := models.NewVector(newX, newY)

			if !isValidPoint(grid, int(gridWidth), int(gridHeight), width, height, newPoint, radius, cellSize) {
				continue
			}

			// step 3
			points = append(points, newPoint)
			insertPoint(grid, newPoint, cellSize)
			active = append(active, newPoint)
			found = true
			break
		}

		if !found {
			active[ridx] = active[len(active)-1]
			active = active[:len(active)-1]
		}

	}

	return points
}

func insertPoint(grid [][]models.Vector, point models.Vector, cellsize float64) {
	xIndex := int(math.Floor(point.X / cellsize))
	yIndex := int(math.Floor(point.Y / cellsize))
	grid[xIndex][yIndex] = point
}

func isValidPoint(
	grid [][]models.Vector,
	gwidth, gheight, width, height int,
	p models.Vector,
	radius, cellsize float64,
) bool {
	/* Make sure the point is on the screen */
	if p.X < 0 || p.X >= float64(width) || p.Y < 0 || p.Y >= float64(height) {
		return false
	}

	/* Check neighboring eight cells */
	xindex := math.Floor(p.X / cellsize)
	yindex := math.Floor(p.Y / cellsize)
	i0 := math.Max(xindex-1, 0)
	i1 := math.Min(xindex+1, float64(gwidth)-1)
	j0 := math.Max(yindex-1, 0)
	j1 := math.Min(yindex+1, float64(gheight)-1)

	for i := i0; i <= i1; i++ {
		for j := j0; j <= j1; j++ {
			gp := grid[int(i)][int(j)]
			if gp.X != -1 && gp.Y != -1 {
				if gp.Dist(p) < radius {
					return false
				}
			}
		}
	}

	/* If we get here, return true */
	return true
}
