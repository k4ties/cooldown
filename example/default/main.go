package main

import (
	"fmt"
	"github.com/k4ties/cooldown"
	"log"
	"time"
)

func main() {
	cd := cooldown.NewWithVal(&handler[string]{})
	cd.Start(time.Second*20, "start unique string")

	<-time.After(time.Second * 2)
	fmt.Println("Renewing")
	cd.Renew("renew unique string")

	<-time.After(time.Second * 5)
	cd.Stop("stop unique string") // stop cause cancelled
	// start
	cd.Start(time.Second, "start 2 unique string")

	<-time.After(time.Second * 2)
}

func l(f string, a ...any) {
	log.Printf(f, a...)
}

type handler[T any] struct{}

func (h handler[T]) HandleStart(_ *cooldown.ContextWithVal[T], val T) {
	l("started [(val=%#v)]", val)
}

func (h handler[T]) HandleRenew(_ *cooldown.ContextWithVal[T], val T) {
	l("renew [(val=%#v)]", val)
}

func (h handler[T]) HandleTick(_ *cooldown.WithVal[T], current int64, val T) {
	// Value is nil here because HandleTick called by cooldown internally, it stores
	// to HandleTick zero T value.

	if current == 1 {
		l("first tick [(val=%#v)]", val)
		return
	}

	// e.g: if your ticker calls 50 times in a second change 20 to 50
	if (current % 20) == 0 {
		// second passed
		l("tick second [(tick=%d), (val=%#v)]", current, val)
	}
}

func (h handler[T]) HandleStop(_ *cooldown.WithVal[T], cause cooldown.StopCause, val T) {
	// note that if cause is expired value will be nil because cooldown sets it to zero value
	l("stop cause=%v [(val=%#v)]", cause, val)
}
