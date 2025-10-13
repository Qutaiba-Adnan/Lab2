package main

type WaterRequest struct { // sent to every other truck
    FromID    int
    Timestamp int64 
}

type WaterReply struct { // other trucks send approval
    FromID    int
    ToID      int   // who requested
    Timestamp int64
    Granted   bool
}

type WaterRelease struct { // water is now available
    FromID    int
    Timestamp int64
}
