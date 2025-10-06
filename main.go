package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	g := NewGrid()

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

	fmt.Println("Initial grid:")
	g.Print()

	for step := 0; step < 10; step++ {
		fmt.Printf("Step %d\n", step)

		g.SpreadFire()

		for _, t := range trucks {
			t.Move(g)
			t.Extinguish(g)
		}

		g.Print()
		time.Sleep(500 * time.Millisecond)
	}
}
