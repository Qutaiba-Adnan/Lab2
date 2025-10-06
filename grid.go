package main

import (
	"fmt"
	"math/rand"
)

const GridSize = 20

type Cell int

const (
	Empty Cell = iota
	Fire
	Extinguished
	Truck
)

type Grid struct {
	Cells [GridSize][GridSize]Cell
}

// Ignite a random fire
func (g *Grid) IgniteFire() {
	x := rand.Intn(GridSize)
	y := rand.Intn(GridSize)
	g.Cells[y][x] = Fire

}

// spread fire to neighboring cells
func (g *Grid) SpreadFire() {
	var newFire [][2]int
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			if g.Cells[y][x] == Fire {
				if rand.Float64() < 0.2 { // "20% chance to spread"
					dx := rand.Intn(3) - 1
					dy := rand.Intn(3) - 1
					nx, ny := x+dx, y+dy
					if nx >= 0 && nx < GridSize && ny >= 0 && ny < GridSize && g.Cells[ny][nx] == Empty {
						newFire = append(newFire, [2]int{nx, ny})
					}
				}
			}
		}
	}
	for _, f := range newFire {
		g.Cells[f[1]][f[0]] = Fire
	}
}

// Print the grid
func (g *Grid) Print() {
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			switch g.Cells[y][x] {
			case Empty:
				fmt.Print(".")
			case Fire:
				fmt.Print("F")
			case Extinguished:
				fmt.Print("E")
			case Truck:
				fmt.Print("T")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	var g Grid
	g.IgniteFire()
	g.Print()
}
