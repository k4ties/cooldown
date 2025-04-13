package main

import (
	"fmt"
	"github.com/k4ties/cooldown"
	"time"
)

func main() {
	cd := cooldown.New(&handler{})
	cd.Start(time.Second * 20)

	<-time.After(time.Second * 2)
	cd.Renew()

	<-time.After(time.Second * 5)
	cd.Stop() // stop cause cancelled
	// start
	cd.Start(time.Second)

	select {
	case <-time.After(time.Second * 2):
		// tick 1s
		// stop cause expired
	}
}

func l(a ...any) {
	fmt.Println(a...)
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
	l("stop cause", cause)
}
