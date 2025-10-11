package main

import (
	"fmt"
	"time"
	"firelab/internal/messaging" 
)

var totalTrucks int

func main() {
	var g Grid

	// Create and place two firetrucks
	nc := messaging.Connect()
	defer messaging.Drain(nc)

	trucks := []Firetruck{
		{X: 5, Y: 5, ID: 1, nc: nc},
		{X: 15, Y: 10, ID: 2, nc: nc},
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

	for step := 0; step < 20; step++ {
		fmt.Printf("Step %d\n", step)

		g.SpreadFire()

		for i := range trucks {
			trucks[i].Extinguish(&g)
			trucks[i].Move(&g)
			trucks[i].Extinguish(&g)
		}

		// Print grid state
		g.Print()

		time.Sleep(500 * time.Millisecond)
	}
}
