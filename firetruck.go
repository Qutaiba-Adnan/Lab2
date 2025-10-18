package main

import (
	"firelab/internal/messaging"
	"fmt"
	"math"

	"github.com/nats-io/nats.go"
)

// ----------------------------------------------------
// Firetruck struct and state
// ----------------------------------------------------
type Firetruck struct {
	X, Y             int
	ID               int
	Failed           bool
	Clock            int
	nc               *nats.Conn
	usingWater       bool
	requestedWater   bool
	requestTimestamp int64
	deferredRequests []WaterRequest
	approvals        map[int]bool
	waitCounter      int
}

// ----------------------------------------------------
// Clock management
// ----------------------------------------------------
func (t *Firetruck) incrementClock() int {
	t.Clock++
	return t.Clock
}

func (t *Firetruck) updateClock(received int) int {
	if received > t.Clock {
		t.Clock = received
	}
	t.Clock++
	return t.Clock
}

// ----------------------------------------------------
// Message Handlers
// ----------------------------------------------------
func (t *Firetruck) OnWaterRequest(req WaterRequest) {
	if req.FromID == t.ID {
		return
	}

	t.updateClock(int(req.Timestamp))
	grant := true

	switch {
	case t.usingWater:
		grant = false
	case t.requestedWater:
		if req.Timestamp > t.requestTimestamp ||
			(req.Timestamp == t.requestTimestamp && req.FromID > t.ID) {
			grant = false
		}
	}

	if !grant {
		t.deferredRequests = append(t.deferredRequests, req)
	}

	messaging.PublishJSON(t.nc, "water.reply", WaterReply{
		FromID:    t.ID,
		ToID:      req.FromID,
		Timestamp: int64(t.Clock),
		Granted:   grant,
	})

	if grant {
		logChannel <- fmt.Sprintf("Truck %d granted water to truck %d at Ts:%d", t.ID, req.FromID, t.Clock)
	} else {
		logChannel <- fmt.Sprintf("Truck %d denied water to truck %d at Ts:%d", t.ID, req.FromID, t.Clock)
	}
}

func (t *Firetruck) OnWaterReply(reply WaterReply) {
	if reply.ToID != t.ID {
		return
	}

	t.updateClock(int(reply.Timestamp))
	if reply.Granted {
		t.approvals[reply.FromID] = true
		logChannel <- fmt.Sprintf("Truck %d got approval from Truck %d at Ts:%d", t.ID, reply.FromID, t.Clock)
	}

	if len(t.approvals) == totalTrucks-1 {
		t.usingWater = true
		t.requestedWater = false
		logChannel <- fmt.Sprintf("Truck %d now has water access at Ts:%d", t.ID, t.Clock)
	}
}

func (t *Firetruck) OnWaterRelease(msg WaterRelease) {
	if msg.FromID == t.ID {
		return
	}

	t.updateClock(int(msg.Timestamp))
	for _, req := range t.deferredRequests {
		messaging.PublishJSON(t.nc, "water.reply", WaterReply{
			FromID:    t.ID,
			ToID:      req.FromID,
			Timestamp: int64(t.Clock),
			Granted:   true,
		})
		logChannel <- fmt.Sprintf("Truck %d granted water to truck %d after release at Ts:%d", t.ID, req.FromID, t.Clock)
	}
	t.deferredRequests = nil
}

// ----------------------------------------------------
// Movement and Extinguish
// ----------------------------------------------------
func (t *Firetruck) Place(g *Grid) {
	g.Cells[t.Y][t.X] = Truck
}

func (t *Firetruck) findClosestFire(g *Grid) (int, int, bool) {
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
	return targetX, targetY, targetX != -1
}

func (t *Firetruck) Move(g *Grid) {
	if t.Failed {
		return
	}

	g.Cells[t.Y][t.X] = Empty
	t.incrementClock()

	targetX, targetY, found := t.findClosestFire(g)
	if !found {
		g.Cells[t.Y][t.X] = Truck
		return
	}

	if math.Abs(float64(t.X-targetX)) <= 1 && math.Abs(float64(t.Y-targetY)) <= 1 {
		g.Cells[t.Y][t.X] = Truck
		t.Extinguish(g, targetX, targetY)
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
	logChannel <- fmt.Sprintf("Truck %d moved to (%d,%d) at Ts:%d", t.ID, t.X, t.Y, t.Clock)
}

// ----------------------------------------------------
// Extinguish logic with water coordination
// ----------------------------------------------------
func (t *Firetruck) Extinguish(g *Grid, fireX int, fireY int) {
	if t.Failed {
		return
	}

	if !t.requestedWater && !t.usingWater {
		t.requestedWater = true
		t.requestTimestamp = int64(t.incrementClock())
		req := WaterRequest{FromID: t.ID, Timestamp: t.requestTimestamp}
		messaging.PublishJSON(t.nc, "water.request", req)
		logChannel <- fmt.Sprintf("Truck %d sent WaterRequest at Ts:%d", t.ID, t.Clock)
		return
	}

	if t.requestedWater && !t.usingWater {
		t.waitCounter++
		if t.waitCounter >= totalTrucks-1 {
			t.waitCounter = 0
			t.requestTimestamp = int64(t.incrementClock())
			messaging.PublishJSON(t.nc, "water.request", WaterRequest{
				FromID:    t.ID,
				Timestamp: t.requestTimestamp,
			})
			logChannel <- fmt.Sprintf("Truck %d re-sent WaterRequest at Ts:%d", t.ID, t.Clock)
		}
		return
	}

	// Extinguish fire if water access granted
	t.incrementClock()
	g.Cells[fireY][fireX] = Extinguished
	g.Intensity[fireY][fireX] = 0
	logChannel <- fmt.Sprintf("Truck %d extinguished fire at (%d,%d) Ts:%d", t.ID, fireX, fireY, t.Clock)

	t.usingWater = false
	t.requestedWater = false
	t.approvals = make(map[int]bool)
	t.incrementClock()

	messaging.PublishJSON(t.nc, "water.release", WaterRelease{
		FromID:    t.ID,
		Timestamp: int64(t.Clock),
	})
	logChannel <- fmt.Sprintf("Truck %d released water at Ts:%d", t.ID, t.Clock)
}
