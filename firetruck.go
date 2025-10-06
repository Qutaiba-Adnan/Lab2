package main

import (
	"fmt"
	"math"
)

type Firetruck struct {
	X, Y int
	ID   int
}

func (t *Firetruck) Place(g *Grid) {
	g.Cells[t.Y][t.X] = Truck
}

func (t *Firetruck) Move(g *Grid) {
	g.Cells[t.Y][t.X] = Empty // clear previous position

	if g.Cells[t.Y][t.X] == Fire {
		t.Extinguish(g)
		return
	}

	// Find nearest fire
	targetX, targetY := -1, -1
	minDist := math.MaxFloat64
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			if g.Cells[y][x] == Fire {
				dist := math.Sqrt(math.Pow(float64(t.X-x), 2) + math.Pow(float64(t.Y-y), 2))
				if dist < minDist {
					minDist = dist
					targetX, targetY = x, y
				}
			}
		}
	}

	if targetX == -1 {
		return
	}

	fmt.Printf("Truck %d targeting fire at (%d,%d) from (%d,%d)\n",
		t.ID, targetX, targetY, t.X, t.Y)

	// Move one step toward fire
	if t.X < targetX {
		t.X++
	} else if t.X > targetX {
		t.X--
	}
	if t.Y < targetY {
		t.Y++
	} else if t.Y > targetY {
		t.Y--
	}

	g.Cells[t.Y][t.X] = Truck
}

// Extinguish any nearby fire (within 1-cell radius)
func (t *Firetruck) Extinguish(g *Grid) {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			nx := t.X + dx
			ny := t.Y + dy
			if nx >= 0 && nx < GridSize && ny >= 0 && ny < GridSize {
				if g.Cells[ny][nx] == Fire {
					waterCost := 5 * g.Intensity[ny][nx]
					if g.RequestWater(waterCost) {
						g.Cells[ny][nx] = Extinguished
						g.Intensity[ny][nx] = 0
						fmt.Printf("Truck %d extinguished fire at (%d,%d)! Water left: %d\n",
							t.ID, nx, ny, g.Water)
					}
				}
			}
		}
	}
}
