package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/odin-software/metro/internal/models"
)

const (
	StationColor = "\033[1;36m%s\033[0m"
	TrainColor   = "\033[1;31m%s\033[0m"
)

func PrintMap(width, height int, sts []*models.Station, trs []models.Train) {
	ClearScreen()
	fmt.Printf("Mapa de New Metro\n\n")
	for y := 0; y <= height; y += 50 {
		for x := 0; x <= width; x += 50 {
			foundStation := false
			foundTrain := false
			for _, tr := range trs {
				if tr.Position.CloseTo(float64(x), float64(y), 20) {
					t := fmt.Sprintf(" %s ", centerString(tr.Name, 3))
					fmt.Printf(TrainColor, t)
					foundTrain = true
					break
				}
			}
			if !foundTrain {
				for _, st := range sts {
					if int(st.Position.X) == x && int(st.Position.Y) == y {
						t := fmt.Sprintf(" %s ", centerString(strconv.FormatInt(st.ID, 10), 3))
						fmt.Printf(StationColor, t)
						foundStation = true
						break
					}
				}
			}
			if !foundStation && !foundTrain {
				fmt.Printf("  •  ")
			}
		}
		fmt.Println()
	}
}

func ClearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func StartMap(tick <-chan time.Time, sts []*models.Station, trs []models.Train) {
	for range tick {
		PrintMap(600, 600, sts, trs)
	}
}

func centerString(str string, width int) string {
	spaces := int(float64(width-len(str)) / 2)
	return strings.Repeat(" ", spaces) + str + strings.Repeat(" ", width-(spaces+len(str)))
}
