package main

import (
	"math"
)

type Firetruck struct {
	X, Y   int
	ID     int
	Failed bool // new flag
}

func (t *Firetruck) Place(g *Grid) {
	g.Cells[t.Y][t.X] = Truck
}

func (t *Firetruck) Move(g *Grid) {
	if t.Failed {
		return
	}

	g.Cells[t.Y][t.X] = Empty

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

// Extinguish nearby fires (if not failed)
func (t *Firetruck) Extinguish(g *Grid) {
	if t.Failed {
		return
	}

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
						println("Truck", t.ID, "extinguished fire at", nx, ny)
					}
				}
			}
		}
	}
}
