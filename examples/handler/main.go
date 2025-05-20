package main

import (
	"context"
	"github.com/k4ties/cooldown"
	"log"
	"time"
)

func main() {
	cooldown.Proc.StartTracking(context.Background(), time.Second/20)

	cd := cooldown.New()
	cd.Handle(&handler{})

	cd.Start(time.Second)
	lf("started cooldown")
	lf("remaining: %s", cd.Remaining().String())
	<-time.After(time.Millisecond * 200)
	lf("remaining (timeout 200 miliseconds): %s", cd.Remaining().String())
	cd.Renew()
	lf("remaining (after renew): %s", cd.Remaining().String())
	lf("active: %t", cd.Active())
	<-time.After(time.Second + time.Millisecond)
	lf("(timeout after second) active: %t", cd.Active())

	cd.Start(time.Second * 3)
	lf("started cooldown for 3 seconds")
	lf("active: %t", cd.Active())
	lf("remaining: %s", cd.Remaining().String())

	<-time.After(time.Second * 4)
}

func lf(f string, a ...any) {
	log.Printf(f, a...)
}

type handler struct{}

func (h handler) HandleStart(_ *cooldown.Context) {
	lf("handle start")
}

func (h handler) HandleRenew(_ *cooldown.Context) {
	lf("handle renew")
}

func (h handler) HandleStop(_ *cooldown.CoolDown, cause cooldown.StopCause) {
	lf("handle stop, cause=%v", cause)
}
