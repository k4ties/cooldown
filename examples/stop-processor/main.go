package main

import (
	"context"
	"github.com/k4ties/cooldown"
	"time"
)

func init() {
	cooldown.StartProcessorOnInit.Store(false)
}

func main() {
	go cooldown.Proc.StartTracking(context.Background(), time.Second/20)
	select {
	case <-time.After(time.Second):
	}
}
