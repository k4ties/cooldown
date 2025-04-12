package main

import (
	"github.com/k4ties/cooldown"
	"log"
	"time"
)

func main() {
	cd := cooldown.New(&handler{})
	cd.Set(time.Second * 20)

	<-time.After(time.Second * 2)
	cd.Renew()

	<-time.After(time.Second * 5)
	cd.Reset() // stop cause cancelled
	// start
	cd.Set(time.Second)

	select {
	case <-time.After(time.Second * 2):
		// tick 1s
		// stop cause expired
	}
}

func l(a ...any) {
	log.Print(a...)
}

type handler struct{}

func (handler) HandleStart(*cooldown.Context) {
	l("start")
}

func (handler) HandleRenew(*cooldown.Context) {
	l("renew")
}

func (handler) HandleTick(_ *cooldown.CoolDown, current int64) {
	if current == 1 {
		l("first tick second")
		return
	}

	if current%int64(cooldown.TicksPerSecond) == 0 {
		// second passed
		l("tick second")
	}
}

func (handler) HandleStop(_ *cooldown.CoolDown, cause cooldown.StopCause) {
	l("stop cause ", cause)
}
