package main

import (
	"math"
)

type Firetruck struct {
	X, Y int
	ID   int
}

// Place the truck on the grid
func (t *Firetruck) Place(g *Grid) {
	g.Cells[t.Y][t.X] = Truck
}

// Find the closest fire using Manhattan distance
func (t *Firetruck) findClosestFire(g *Grid) (int, int, bool) {
	minDist := math.MaxFloat64
	targetX, targetY := -1, -1
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			if g.Cells[y][x] == Fire {
				dist := math.Abs(float64(t.X-x)) + math.Abs(float64(t.Y-y))
				if dist < minDist {
					minDist = dist
					targetX, targetY = x, y
				}
			}
		}
	}
	if targetX == -1 {
		return 0, 0, false
	}
	return targetX, targetY, true
}

// Move one step toward the closest fire
func (t *Firetruck) Move(g *Grid) {
	g.Cells[t.Y][t.X] = Empty // clear old position

	targetX, targetY, found := t.findClosestFire(g)
	if !found {

		g.Cells[t.Y][t.X] = Truck
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

// Instantly extinguish fire if present
func (t *Firetruck) Extinguish(g *Grid) {
	if g.Cells[t.Y][t.X] == Fire {
		g.Cells[t.Y][t.X] = Extinguished
		g.Intensity[t.Y][t.X] = 0
	}
}
