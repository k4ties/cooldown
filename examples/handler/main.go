package main

import (
	"github.com/k4ties/cooldown"
	"log"
	"time"
)

func main() {
	cd := cooldown.New(cooldown.OptionHandler(new(handler)))

	cd.Start(time.Second)
	lf("started cooldown")
	lf("remaining: %s", cd.Remaining().String())
	<-time.After(time.Millisecond * 200)
	lf("remaining (timeout 200 miliseconds): %s", cd.Remaining().String())
	cd.Renew()
	lf("remaining (after renew): %s", cd.Remaining().String())
	lf("active: %t", cd.Active())
	lf("stoping the cooldown")
	cd.Stop()
	<-time.After(time.Second + time.Millisecond)
	lf("(timeout after second) active: %t, remaining: %s", cd.Active(), cd.Remaining().String())

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
