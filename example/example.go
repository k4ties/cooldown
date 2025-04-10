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
	cd.Set(time.Second * 10) // this line causing renew, must not be stop event

	<-time.After(time.Second * 5)
	cd.Reset()
	cd.Reset()

	select {
	case <-time.After(time.Second * 6):
		// end script
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
	if current%int64(cooldown.TicksPerSecond) == 0 {
		// second passed
		l("tick 1s")
	}
}

func (handler) HandleStop(_ *cooldown.CoolDown, cause cooldown.StopCause) {
	l("stop cause", cause)
}
