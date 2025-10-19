package main

import "fmt"

func fireName(x, y int) string {
	return fmt.Sprintf("fire-%d-%d", x, y)
}

func (t *Firetruck) Name() string {
	return fmt.Sprintf("truck-%d", t.ID)
}
