package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGrid()

	// Start with three fires
	for i := 0; i < 3; i++ {
		g.IgniteFire()
	}

	trucks := []*Firetruck{
		{X: 5, Y: 5, ID: 1},
		{X: 15, Y: 10, ID: 2},
	}

	for _, t := range trucks {
		t.Place(g)
	}

	chief := ChiefTruck{ID: 0}

	fmt.Println("Initial grid:")
	g.Print()

	for step := 0; step < 10; step++ {
		fmt.Printf("Step %d\n", step)

		// Find active fires
		var fires [][2]int
		for y := 0; y < GridSize; y++ {
			for x := 0; x < GridSize; x++ {
				if g.Cells[y][x] == Fire {
					fires = append(fires, [2]int{x, y})
				}
			}
		}

		if len(fires) == 0 {
			fmt.Println("All fires extinguished.")
			break
		}

		// Chief assigns fires
		assignments := chief.AssignFires(trucks, fires)

		// Trucks act
		for _, t := range trucks {
			if target, ok := assignments[t.ID]; ok {
				fmt.Printf("%s heading toward fire-%d-%d\n", t.Name(), target[0], target[1])
				t.Move(g)
				t.Extinguish(g)
			}
		}

		// Spread and print
		g.SpreadFire()
		g.Print()

		time.Sleep(600 * time.Millisecond)
	}
}
