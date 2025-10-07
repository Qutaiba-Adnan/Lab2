package main

import (
	"fmt"
	"math"
)

// ChiefTruck coordinates trucks and assigns fires
type ChiefTruck struct {
	ID int
}

// Assign each fire to the nearest truck
func (c *ChiefTruck) AssignFires(trucks []*Firetruck, fires [][2]int) map[int][2]int {
	assignments := make(map[int][2]int)
	assigned := make(map[[2]int]bool)

	for _, t := range trucks {
		minDist := math.MaxFloat64
		var closest [2]int
		found := false

		for _, f := range fires {
			if assigned[f] {
				continue
			}
			dist := math.Abs(float64(t.X-f[0])) + math.Abs(float64(t.Y-f[1]))
			if dist < minDist {
				minDist = dist
				closest = f
				found = true
			}
		}

		if found {
			assignments[t.ID] = closest
			assigned[closest] = true
			fmt.Printf("ChiefTruck assigned fire-%d-%d to truck-%d\n", closest[0], closest[1], t.ID)
		}
	}

	return assignments
}
