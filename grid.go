package main

import (
	"fmt"
	"math/rand"
)

const (
	GridSize     = 20
	Empty        = '.'
	Extinguished = 'E'
	Truck        = 'T'
	FireMin      = 1
	FireMax      = 5
)

// Grid represents the simulation world
type Grid struct {
	Cells     [GridSize][GridSize]rune
	Intensity [GridSize][GridSize]int
	StepCount int
}

// ----------------------------------------------------
// Grid initialization
// ----------------------------------------------------
func NewGrid() *Grid {
	g := &Grid{}
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			g.Cells[y][x] = Empty
		}
	}
	return g
}

// ----------------------------------------------------
// Fire ignition
// ----------------------------------------------------
func (g *Grid) IgniteFire() {
	x := rand.Intn(GridSize)
	y := rand.Intn(GridSize)
	g.Cells[y][x] = rune('0' + FireMin)
	g.Intensity[y][x] = FireMin
}

// ----------------------------------------------------
// Fire spread and intensity logic (scaled down)
// ----------------------------------------------------
func (g *Grid) SpreadFire() {
	g.StepCount++

	var newFire [][2]int

	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			if g.Intensity[y][x] >= FireMin {
				intensity := g.Intensity[y][x]

				// Small chance to increase intensity
				if rand.Float64() < 0.15 {
					g.Intensity[y][x]++
					if g.Intensity[y][x] > FireMax {
						g.Intensity[y][x] = FireMax
					}
				}

				// Chance to weaken slightly (cooling effect)
				if rand.Float64() < 0.15 {
					g.Intensity[y][x]--
					if g.Intensity[y][x] < FireMin {
						g.Intensity[y][x] = FireMin
					}
				}

				// Spread only to adjacent cells (very local)
				spreadRange := 1
				spreadChance := 0.02 * float64(intensity) // reduced spread

				for dy := -spreadRange; dy <= spreadRange; dy++ {
					for dx := -spreadRange; dx <= spreadRange; dx++ {
						if dx == 0 && dy == 0 {
							continue
						}
						nx, ny := x+dx, y+dy
						if nx >= 0 && nx < GridSize && ny >= 0 && ny < GridSize {
							if g.Cells[ny][nx] == Empty && rand.Float64() < spreadChance {
								newFire = append(newFire, [2]int{nx, ny})
							}
						}
					}
				}

				// Rare long-distance jump fire (extremely rare)
				if rand.Float64() < 0.003 {
					jx := x + rand.Intn(7) - 3
					jy := y + rand.Intn(7) - 3
					if jx >= 0 && jx < GridSize && jy >= 0 && jy < GridSize && g.Cells[jy][jx] == Empty {
						newFire = append(newFire, [2]int{jx, jy})
					}
				}

				// Slightly higher burnout chance (scaled with intensity)
				if rand.Float64() < 0.07*float64(intensity) {
					g.Cells[y][x] = Extinguished
					g.Intensity[y][x] = 0
				}
			}
		}
	}

	// Add new fires
	for _, f := range newFire {
		g.Cells[f[1]][f[0]] = '1'
		g.Intensity[f[1]][f[0]] = 1
	}

	// Ignite a random new fire less frequently (every 10 steps)
	if g.StepCount%10 == 0 {
		x := rand.Intn(GridSize)
		y := rand.Intn(GridSize)
		if g.Cells[y][x] == Empty {
			g.Cells[y][x] = '1'
			g.Intensity[y][x] = 1
		}
	}
}

// ----------------------------------------------------
// Grid visualization
// ----------------------------------------------------
func (g *Grid) Print() {
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			if g.Intensity[y][x] > 0 {
				fmt.Printf("%d", g.Intensity[y][x])
			} else {
				fmt.Printf("%c", g.Cells[y][x])
			}
		}
		fmt.Println()
	}
}
