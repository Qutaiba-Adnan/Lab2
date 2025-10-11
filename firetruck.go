package main

import (
	"math"
	"fmt"
	"firelab/internal/messaging" 
    "github.com/nats-io/nats.go"
)

type Firetruck struct {
	X, Y int
	ID   int
	Clock int

	nc *nats.Conn // points to connection in main.go
	usingWater      bool   // true after approval, before WaterRelease
    requestedWater  bool   // true after sending WaterRequest
    requestTimestamp int64
	deferredRequests []WaterRequest // requests this truck denies
}

func (t *Firetruck) incrementClock() int {
	t.Clock++
	v := t.Clock
	return v
}

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
        fmt.Printf("Truck %d granted water to truck %d at %d\n", t.ID, req.FromID, t.Clock)
    } else {
        fmt.Printf("Truck %d denied water to truck %d at %d\n", t.ID, req.FromID, t.Clock)
    }
}

func (t *Firetruck) OnWaterReply(reply WaterReply) {
    if reply.ToID != t.ID {
        return // not for this truck
    }

    t.updateClock(int(reply.Timestamp))

	approvals := 0
    if reply.Granted {
        approvals++
    } else {
        
    }

	if approvals == totalTrucks -1 {
		t.requestedWater = false
		t.usingWater = true
    	fmt.Printf("Truck %d is accessing water\n", t.ID)
	} else {

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
        fmt.Printf("Truck %d granted water to truck %d at %d\n", t.ID, req.FromID, t.Clock)
    }
	t.deferredRequests = nil
}

// Place the truck on the grid
func (t *Firetruck) Place(g *Grid) {
	g.Cells[t.Y][t.X] = Truck
}

// Find the closest fire using Manhattan distance
func (t *Firetruck) findClosestFire(g *Grid) (int, int, bool) {
	minDist := math.MaxFloat64
	targetX, targetY := -1, -1
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			if g.Cells[y][x] == Fire {
				dist := math.Abs(float64(t.X-x)) + math.Abs(float64(t.Y-y))
				if dist < minDist {
					minDist = dist
					targetX, targetY = x, y
				}
			}
		}
	}
	if targetX == -1 {
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

	fmt.Printf("Truck %d moved. Ts: %d\n", t.ID, t.Clock)
}

func (t *Firetruck) Extinguish(g *Grid) {
    if g.Cells[t.Y][t.X] == Fire {

        if !t.requestedWater && !t.usingWater {
            t.requestedWater = true
            t.requestTimestamp = int64(t.incrementClock())

            req := WaterRequest{
                FromID:    t.ID,
                Timestamp: t.requestTimestamp,
            }
            messaging.PublishJSON(t.nc, "water.request", req)
            fmt.Printf("[Truck %d sent WaterRequest at Ts:%d\n", t.ID, t.Clock)
            return 
        }

		if t.requestedWater && !t.usingWater {
			// waiting for approvals
			return
		}

        t.incrementClock()
        g.Cells[t.Y][t.X] = Extinguished
        g.Intensity[t.Y][t.X] = 0
        fmt.Printf("[Truck %d Extinguished fire at Ts: %d\n", t.ID, t.Clock)

        t.usingWater = false
        t.requestedWater = false
        t.incrementClock()

        messaging.PublishJSON(t.nc, "water.release", WaterRelease{
            FromID:    t.ID,
            Timestamp: int64(t.Clock),
        })
        fmt.Printf("[Truck %d released water at Ts:%d\n", t.ID, t.Clock)
    }
}
