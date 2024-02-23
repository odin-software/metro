package main

import (
	"fmt"
	"internal/model"
	"os"
	"os/exec"
	"strings"
)

const (
	StationColor = "\033[1;36m%s\033[0m"
	TrainColor   = "\033[1;31m%s\033[0m"
)

func getIdNumber(st model.Station) string {
	a := strings.Split(st.ID, "-")
	return a[1]
}

func PrintMap(width, height int, sts []*model.Station, trs []model.Train) {
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
						t := fmt.Sprintf(" %s ", centerString(getIdNumber(*st), 3))
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
	fmt.Println(len(sts[1].Trains))
}

func ClearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func centerString(str string, width int) string {
	spaces := int(float64(width-len(str)) / 2)
	return strings.Repeat(" ", spaces) + str + strings.Repeat(" ", width-(spaces+len(str)))
}
