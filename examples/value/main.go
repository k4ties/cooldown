package main

import (
	"log"
	"time"

	"github.com/k4ties/cooldown"
)

func main() {
	cd := cooldown.NewValued[string](cooldown.ValuedOptionHandler(new(handler)))

	cd.Start(time.Second, "unique value for start 1")
	lf("started cooldown")
	lf("remaining: %s", cd.Remaining().String())
	<-time.After(time.Millisecond * 200)
	lf("remaining (timeout 200 miliseconds): %s", cd.Remaining().String())
	cd.Renew("unique value for renew")
	lf("remaining (after renew): %s", cd.Remaining().String())
	lf("active: %t", cd.Active())
	cd.Stop("unique value for stop")
	<-time.After(time.Second + time.Millisecond)
	lf("(timeout after second) active: %t, remaining: %s", cd.Active(), cd.Remaining().String())

	cd.Start(time.Second*3, "unique value for start 2")
	lf("started cooldown for 3 seconds")
	lf("active: %t", cd.Active())
	lf("remaining: %s", cd.Remaining().String())

	<-time.After(time.Second * 4)
}

func lf(f string, a ...any) {
	log.Printf(f, a...)
}

type handler struct{}

func (handler) HandleStart(_ *cooldown.ValuedContext[string], dur time.Duration, val string) {
	lf("handle start [val='%v',dur='%s']", val, dur)
}

func (handler) HandleRenew(_ *cooldown.ValuedContext[string], dur time.Duration, val string) {
	lf("handle renew [val='%v',dur='%s']", val, dur)
}

func (handler) HandleStop(_ *cooldown.Valued[string], cause cooldown.StopCause, val string) {
	// val may be zero, because if it is expiration stop cause, library sets zero value to it
	lf("handle stop [cause='%v', val='%v']", cause, val)
}
