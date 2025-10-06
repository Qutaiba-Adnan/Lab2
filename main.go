package main

import (
	"fmt"
)

const GridSize = 20

// -1 = truck, 0 = empty, >0 = fire intensity
var grid [GridSize][GridSize]int

// Position of the truck
var truckX, truckY int

func printGrid() {
	for i := 0; i < GridSize; i++ {
		for j := 0; j < GridSize; j++ {
			switch {
			case grid[i][j] == -1:
				fmt.Print("T ")
			case grid[i][j] == 0:
				fmt.Print(". ")
			default:
				fmt.Printf("%d ", grid[i][j])
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	// Set truck position
	truckX, truckY = 10, 10
	grid[truckX][truckY] = -1

	fmt.Println("=== Initial Grid State ===")
	printGrid()
}
