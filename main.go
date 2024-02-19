package metro

import (
	// "bufio"

	"fmt"
	"internal/model"

	// "os"
	"time"
)

var stationHashFunction = func(station model.Station) string {
	return station.ID
}

func main() {
	// Timing and configuration
	// scnr := bufio.NewScanner(os.Stdin)
	ticker := time.NewTicker(20 * time.Millisecond)
	quit := make(chan struct{})

	// Filling graph data.
	g := model.NewNetwork[model.Station](stationHashFunction)
	sts, lines := GenerateTestData()
	g.InsertVertices(sts)
	g.InsertEdge(sts[0], sts[1], []model.Vector{model.NewVector(50.0, 250.0), model.NewVector(150.0, 200.0)})
	g.InsertEdge(sts[1], sts[2], []model.Vector{model.NewVector(250.0, 100.0)})
	g.InsertEdge(sts[1], sts[5], []model.Vector{model.NewVector(300.0, 300.0)})
	g.InsertEdge(sts[1], sts[3], []model.Vector{model.NewVector(350.0, 200.0), model.NewVector(400.0, 150.0), model.NewVector(400.0, 50.0)})
	g.InsertEdge(sts[3], sts[4], []model.Vector{model.NewVector(550.0, 100.0), model.NewVector(600.0, 100.0)})
	g.InsertEdge(sts[3], sts[10], []model.Vector{})
	g.InsertEdge(sts[3], sts[11], []model.Vector{model.NewVector(600.0, 50.0)})
	g.InsertEdge(sts[5], sts[6], []model.Vector{model.NewVector(100.0, 500.0)})
	g.InsertEdge(sts[7], sts[8], []model.Vector{model.NewVector(500.0, 450.0)})
	g.InsertEdge(sts[8], sts[9], []model.Vector{model.NewVector(500.0, 250.0), model.NewVector(550.0, 200.0)})

	// Creating the train and queing some destinations.
	make := model.NewMake("4-Legged-chu", "A type of fast train.", 0.01, 4)
	train := model.NewTrain("Chu", make, sts[0].Location, sts[0], lines[0], &g)

	go func() {
		for {
			select {
			case <-ticker.C:
				train.Update()
				fmt.Println(train.Position.X, train.Position.Y)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// Adding the reporter server.
	ReporterServer()

	// for {
	// 	// getting input
	// 	fmt.Print("metro > ")
	// 	scnr.Scan()
	// 	if scnr.Text() == "stop" {
	// 		quit <- struct{}{}
	// 	}
	// 	if scnr.Text() == "exit" {
	// 		os.Exit(0)
	// 	}
	// }
}
