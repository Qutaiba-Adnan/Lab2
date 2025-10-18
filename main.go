package main

import (
	"firelab/internal/messaging"
	"fmt"
	"time"
)

var totalTrucks int

func main() {
	var g Grid

	// Create and place two firetrucks
	nc := messaging.Connect()
	defer messaging.Drain(nc)

	trucks := []Firetruck{
		{X: 5, Y: 5, ID: 1, nc: nc, approvals: make(map[int]bool)},
		{X: 15, Y: 10, ID: 2, nc: nc, approvals: make(map[int]bool)},
	}
	totalTrucks = len(trucks)

	for i := range trucks {
		// each truck subscribes to incoming messages from other trucks
		messaging.SubscribeJSON[WaterRequest](nc, "water.request", trucks[i].OnWaterRequest)
		messaging.SubscribeJSON[WaterReply](nc, "water.reply", trucks[i].OnWaterReply)
		messaging.SubscribeJSON[WaterRelease](nc, "water.release", trucks[i].OnWaterRelease)

		trucks[i].Place(&g)
	}

	g.IgniteFire()

	for step := 0; step < 10; step++ {
		fmt.Printf("Step %d\n", step)

		g.SpreadFire()

		for i := range trucks {
			trucks[i].Move(&g)
		}

		// Print grid state
		g.Print()

	logDrain:
		for {
			select {
			case log := <-logChannel:
				fmt.Println(log)
			default:
				break logDrain
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\n Simulation ended.")
}
