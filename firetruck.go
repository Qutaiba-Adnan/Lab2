package main

import (
	"math"
	"fmt"
	"firelab/internal/messaging" 
    "github.com/nats-io/nats.go"
)

type Firetruck struct {
<<<<<<< HEAD
	X, Y   int
	ID     int
	Failed bool // new flag
=======
	X, Y int
	ID   int
	Clock int

	nc *nats.Conn // points to connection in main.go
	usingWater      bool   // true after approval, before WaterRelease
    requestedWater  bool   // true after sending WaterRequest
    requestTimestamp int64
	deferredRequests []WaterRequest // requests this truck denies
	approvals map[int]bool 
	waitCounter int
}

func (t *Firetruck) incrementClock() int {
	t.Clock++
	v := t.Clock
	return v
}

var logChannel = make(chan string, 100) // store logs here

func (t *Firetruck) updateClock(received int) int {
    if received > t.Clock { t.Clock = received }
    t.Clock++
    v := t.Clock
    return v
}

func (t *Firetruck) OnWaterRequest(req WaterRequest) {
    if req.FromID == t.ID {
        return // ignore own request
    }

    t.updateClock(int(req.Timestamp))

    grant := true

    switch {
    case t.usingWater:
        grant = false
    case t.requestedWater:
        // both trucks want water, compare timestamps 
        if req.Timestamp > t.requestTimestamp ||
            (req.Timestamp == t.requestTimestamp && req.FromID > t.ID) {
            grant = false
        }
    }

	if !grant {
        t.deferredRequests = append(t.deferredRequests, req)     // remember to grant later
    }

    messaging.PublishJSON(t.nc, "water.reply", WaterReply{
        FromID:    t.ID,
        ToID:      req.FromID,
        Timestamp: int64(t.Clock),
        Granted:   grant,
    })

    if grant {
        logChannel <- fmt.Sprintf("Truck %d granted water to truck %d at Ts:%d\n", t.ID, req.FromID, t.Clock)
    } else {
        logChannel <- fmt.Sprintf("Truck %d denied water to truck %d at Ts:%d\n", t.ID, req.FromID, t.Clock)
    }
}

func (t *Firetruck) OnWaterReply(reply WaterReply) {
    if reply.ToID != t.ID {
        return // not for this truck
    }

    t.updateClock(int(reply.Timestamp))

    if reply.Granted {
        t.approvals[reply.FromID] = true
        logChannel <- fmt.Sprintf("Truck %d received approval from Truck %d at Ts:%d\n", t.ID, reply.FromID, t.Clock)
    } else {
        logChannel <- fmt.Sprintf("Truck %d was denied by Truck %d at Ts:%d\n", t.ID, reply.FromID, t.Clock)
    }

    // check if we have all approvals
    if len(t.approvals) == totalTrucks-1 {
        t.usingWater = true
        t.requestedWater = false
        logChannel <- fmt.Sprintf("Truck %d has access to water at Ts:%d\n", t.ID, t.Clock)
    }
}

func (t *Firetruck) OnWaterRelease(msg WaterRelease) {
    if msg.FromID == t.ID {
        return
    }

    t.updateClock(int(msg.Timestamp)) // rest of the firetrucks update their clocks

	for _, req := range t.deferredRequests {
        messaging.PublishJSON(t.nc, "water.reply", WaterReply{
            FromID:    t.ID,
            ToID:      req.FromID,
            Timestamp: int64(t.Clock),
            Granted:   true,
        })
        logChannel <- fmt.Sprintf("Truck %d granted water to truck %d at Ts:%d\n", t.ID, req.FromID, t.Clock)
    }
	t.deferredRequests = nil
>>>>>>> origin/3-task-2-coordination-with-logical-clocks
}

func (t *Firetruck) Place(g *Grid) {
	g.Cells[t.Y][t.X] = Truck
}

func (t *Firetruck) Move(g *Grid) {
	if t.Failed {
		return
	}

	g.Cells[t.Y][t.X] = Empty

	// Find nearest fire
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

	if targetX == -1 {
<<<<<<< HEAD
=======
		return 0, 0, false
	}
	return targetX, targetY, true
}

// Move one step toward the closest fire
func (t *Firetruck) Move(g *Grid) {
	g.Cells[t.Y][t.X] = Empty // clear old position
	t.incrementClock() // increase local timestamp by 1

	targetX, targetY, found := t.findClosestFire(g)
	if !found {

		g.Cells[t.Y][t.X] = Truck
>>>>>>> origin/3-task-2-coordination-with-logical-clocks
		return
	}

	if math.Abs(float64(t.X - targetX)) <= 1 && math.Abs(float64(t.Y - targetY)) <= 1 {
		// if the truck is right next to the fire, extinguish 
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

	logChannel <- fmt.Sprintf("Truck %d moved at Ts:%d\n", t.ID, t.Clock)
}

<<<<<<< HEAD
// Extinguish nearby fires (if not failed)
func (t *Firetruck) Extinguish(g *Grid) {
	if t.Failed {
		return
	}

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			nx := t.X + dx
			ny := t.Y + dy
			if nx >= 0 && nx < GridSize && ny >= 0 && ny < GridSize {
				if g.Cells[ny][nx] == Fire {
					waterCost := 5 * g.Intensity[ny][nx]
					if g.RequestWater(waterCost) {
						g.Cells[ny][nx] = Extinguished
						g.Intensity[ny][nx] = 0
						println("Truck", t.ID, "extinguished fire at", nx, ny)
					}
				}
			}
		}
=======
func (t *Firetruck) Extinguish(g *Grid, fireX int, fireY int) {
	if !t.requestedWater && !t.usingWater {
		logChannel <- fmt.Sprintf("Truck %d encountered a fire at Ts:%d\n", t.ID, t.Clock)
		t.requestedWater = true
		t.requestTimestamp = int64(t.incrementClock())

		req := WaterRequest{
			FromID:    t.ID,
			Timestamp: t.requestTimestamp,
		}
		messaging.PublishJSON(t.nc, "water.request", req)
		logChannel <- fmt.Sprintf("Truck %d sent WaterRequest at Ts:%d\n", t.ID, t.Clock)
		return 
>>>>>>> origin/3-task-2-coordination-with-logical-clocks
	}

	if t.requestedWater && !t.usingWater {
		t.waitCounter++
		if t.waitCounter >= totalTrucks - 1 {
			t.waitCounter = 0
			t.requestTimestamp = int64(t.incrementClock())
			messaging.PublishJSON(t.nc, "water.request", WaterRequest{
				FromID:    t.ID,
				Timestamp: t.requestTimestamp,
			})
			logChannel <- fmt.Sprintf("Truck %d re-sent WaterRequest at Ts:%d\n", t.ID, t.Clock)
		}
		return
	}

	t.incrementClock()
	g.Cells[fireY][fireX] = Extinguished
	g.Intensity[fireY][fireX] = 0
	logChannel <- fmt.Sprintf("Truck %d Extinguished fire at Ts:%d\n", t.ID, t.Clock)

	t.usingWater = false
	t.requestedWater = false
	t.approvals = make(map[int]bool)
	t.incrementClock()

	messaging.PublishJSON(t.nc, "water.release", WaterRelease{
		FromID:    t.ID,
		Timestamp: int64(t.Clock),
	})
	logChannel <- fmt.Sprintf("Truck %d released water at Ts:%d\n", t.ID, t.Clock)
    
}
