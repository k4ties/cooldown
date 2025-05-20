package main

import (
	"context"
	"github.com/k4ties/cooldown"
	"log"
	"time"
)

func main() {
	cooldown.Proc.StartTracking(context.Background(), time.Second/20)

	basic := cooldown.NewBasic()
	lf("started basic cooldown")
	basic.Set(time.Second)

	lf("active: %t", basic.Active())
	<-time.After(time.Second + time.Millisecond)
	lf("(timeout) active: %t", basic.Active())

	basic.Set(time.Second * 3)
	lf("started basic cooldown for 3 seconds")
	lf("active: %t", basic.Active())
	lf("remaining: %s", basic.Remaining().String())

	<-time.After(time.Second * 4)
}

func lf(f string, a ...any) {
	log.Printf(f, a...)
}
