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
	Cells     [GridSize][GridSize]Cell
	Intensity [GridSize][GridSize]int
	Water     int
}

// Initialize the grid with water supply
func NewGrid() *Grid {
	g := &Grid{Water: 100}
	return g
}

func (g *Grid) IgniteFire() {
	x := rand.Intn(GridSize)
	y := rand.Intn(GridSize)
	g.Cells[y][x] = Fire
	g.Intensity[y][x] = 1
}

// Spread fires and increase intensity of existing ones
func (g *Grid) SpreadFire() {
	var newFire [][2]int
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			if g.Cells[y][x] == Fire {

				if rand.Float64() < 0.2 {
					g.Intensity[y][x]++
					if g.Intensity[y][x] > 5 {
						g.Intensity[y][x] = 5
					}
				}

				// Slow spread: 5% base + 1% per intensity level
				spreadChance := 0.05 + 0.01*float64(g.Intensity[y][x])
				if rand.Float64() < spreadChance {
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
		g.Intensity[f[1]][f[0]] = 1
	}
}

// Request water from shared tank
func (g *Grid) RequestWater(amount int) bool {
	if g.Water >= amount {
		g.Water -= amount
		fmt.Printf(" Water requested: %d used, %d remaining.\n", amount, g.Water)
		return true
	}
	fmt.Println(" Not enough water! Truck must wait for refill.")
	return false
}

// Refill shared water tank (for testing)
func (g *Grid) RefillWater() {
	g.Water = 100
	fmt.Println(" Water tank refilled to 100.")
}

func (g *Grid) Print() {
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			switch g.Cells[y][x] {
			case Empty:
				fmt.Print(".")
			case Fire:
				fmt.Print(g.Intensity[y][x])
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
