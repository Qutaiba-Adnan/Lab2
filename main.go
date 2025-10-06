package main

import (
	"fmt"
	"time"
)

func main() {
	var g Grid

	// Create and place two firetrucks
	trucks := []Firetruck{
		{X: 5, Y: 5, ID: 1},
		{X: 15, Y: 10, ID: 2},
	}

	for i := range trucks {
		trucks[i].Place(&g)
	}

	// Ignite an initial fire
	g.IgniteFire()

	// Run simulation for 20 steps
	for step := 0; step < 20; step++ {
		fmt.Printf("Step %d\n", step)

		// Spread fires
		g.SpreadFire()

		// Each truck moves and tries to extinguish fires
		for i := range trucks {
			trucks[i].Move(&g)
			trucks[i].Extinguish(&g)
		}

		// Print grid state
		g.Print()

		// Wait a bit between steps
		time.Sleep(500 * time.Millisecond)
	}
}
