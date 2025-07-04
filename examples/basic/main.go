package main

import (
	"github.com/k4ties/cooldown"
	"log"
	"time"
)

func main() {
	basic := new(cooldown.Basic)
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
