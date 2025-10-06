package main

import (
	"math/rand"
)

type Firetruck struct {
	X, Y int
	ID   int
}

// Place the truck on the grid
func (t *Firetruck) Place(g *Grid) {
	g.Cells[t.Y][t.X] = Truck
}

// MOve randomly (N,S,E,W)
func (t *Firetruck) Move(g *Grid) {
	g.Cells[t.Y][t.X] = Empty // clear old position

	dx, dy := 0, 0
	switch rand.Intn(4) {
	case 0:
		dx = -1 // left
	case 1:
		dx = 1 // right
	case 2:
		dy = -1 // Up
	case 3:
		dy = 1 // down
	}

	newX := t.X + dx
	newY := t.Y + dy

	if newX >= 0 && newX < GridSize && newY >= 0 && newY < GridSize {
		t.X, t.Y = newX, newY
	}

	g.Cells[t.Y][t.X] = Truck

}

// Put out fire if truck is on a fire cell
func (t *Firetruck) Extinguish(g *Grid) {
	if g.Cells[t.Y][t.X] == Fire {
		g.Cells[t.Y][t.X] = Extinguished
	}
}
