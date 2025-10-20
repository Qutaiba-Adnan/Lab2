package main

import (
	"firelab/internal/messaging"
	"fmt"
	"log"
	"math/rand"
	"time"
)

var totalTrucks = 3 // number of trucks
var logChannel = make(chan string, 200)

func main() {
	rand.Seed(time.Now().UnixNano())

	// ----------------------------------------------------
	// Connect to NATS
	// ----------------------------------------------------
	nc := messaging.Connect()
	defer messaging.Drain(nc)

	// ----------------------------------------------------
	// Initialize grid and ignite fires
	// ----------------------------------------------------
	g := NewGrid()
	for i := 0; i < 3; i++ {
		g.IgniteFire()
	}

	trucks := []*Firetruck{
		{X: 5, Y: 5, ID: 1, nc: nc, approvals: make(map[int]bool)},
		{X: 10, Y: 10, ID: 2, nc: nc, approvals: make(map[int]bool)},
		{X: 15, Y: 15, ID: 3, nc: nc, approvals: make(map[int]bool)},
	}
	totalTrucks = len(trucks)

	for _, t := range trucks {
		t.Place(g)
	}

	// ----------------------------------------------------
	// Subscribe to NATS topics
	// ----------------------------------------------------
	for _, t := range trucks {
		_, err := messaging.SubscribeJSON[WaterRequest](nc, "water.request", t.OnWaterRequest)
		if err != nil {
			log.Fatalf("subscribe error: %v", err)
		}
		_, err = messaging.SubscribeJSON[WaterReply](nc, "water.reply", t.OnWaterReply)
		if err != nil {
			log.Fatalf("subscribe error: %v", err)
		}
		_, err = messaging.SubscribeJSON[WaterRelease](nc, "water.release", t.OnWaterRelease)
		if err != nil {
			log.Fatalf("subscribe error: %v", err)
		}
	}

	// ----------------------------------------------------
	// Simulation loop
	// ----------------------------------------------------
	fmt.Println("Initial Grid:")
	g.Print()

	for step := 0; step < 30; step++ {
		fmt.Printf("\n=== Step %d ===\n", step)

		// Simulate truck failure (random)
		if rand.Float64() < 0.05 {
			failedTruck := trucks[rand.Intn(len(trucks))]
			failedTruck.Failed = true
			fmt.Printf("Truck %d failed at step %d!\n", failedTruck.ID, step)
		}

		g.SpreadFire()

		for _, t := range trucks {
			t.Move(g)
		}

		g.Print()

		for len(logChannel) > 0 {
			fmt.Println(<-logChannel)
		}

		time.Sleep(400 * time.Millisecond)
	}

	fmt.Println("\nSimulation ended.")
}
