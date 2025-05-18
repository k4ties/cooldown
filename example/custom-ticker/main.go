package main

import (
	"context"
	"fmt"
	"github.com/k4ties/cooldown"
	"log"
	"math/rand"
	"time"
)

func main() {
	opt := tickerOption()

	cd := cooldown.NewWithVal(&handler[*exampleContextValue]{}, opt)
	cd.Start(time.Second*20, newVal())

	<-time.After(time.Second*2 + (time.Second / 2))
	fmt.Println("Renewing")
	cd.Renew(newVal())

	<-time.After(time.Second * 5)
	cd.Stop(newVal()) // stop cause cancelled
	// start
	cd.Start(time.Second, newVal())

	<-time.After(time.Second * 2)
}

func newVal() *exampleContextValue {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &exampleContextValue{uniqueID: int64(r.Intn(100))}
}

func tickerOption() cooldown.Option[*exampleContextValue] {
	return func(cd *cooldown.WithVal[*exampleContextValue]) {
		cooldown.StartFunc[*exampleContextValue](func(data *cooldown.TickData) {
			ctx, cancel := context.WithCancelCause(context.Background())
			data.Context = ctx
			cd.SetCancel(cancel)

			go func() {
				ticker := time.NewTicker(time.Second / 20)
				defer ticker.Stop()

				for {
					select {
					case <-ticker.C:
						ok := cd.Tick(data, newVal()) // we're setting value by ourselves
						if !ok {
							return
						}
					}
				}
			}()
		})(cd)
	}
}

type exampleContextValue struct {
	uniqueID int64
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
