# cooldown
A [Go](https://go.dev/) library, that implements cooldowns in any way.

### Get started

```go
package main

import (
	"fmt"
	"github.com/k4ties/cooldown"
	"sync/atomic"
	"time"
)

func main() {
	cd := cooldown.New(cooldown.OptionHandler(new(handler)))
	cd.Start(time.Second)
	cd.Active() // true
	
	<-time.After(time.Millisecond)
	cd.Renew()
	
	<-time.After(time.Second * 3)
	cd.Stop()
}

type handler struct {
	cooldown.NopHandler // implements cooldown.Handler
	canceled            atomic.Bool
}

func (h *handler) HandleStart(ctx *cooldown.Context) {
	if h.canceled.CompareAndSwap(false, true) {
		// The event can be cancelled
		ctx.Cancel()
	}
	fmt.Println("handle start")
}
func (h *handler) HandleStop(cd *cooldown.CoolDown, cause cooldown.StopCause) {
	// This event cannot be cancelled
	fmt.Printf("handle stop, cause=%v\n", cause)
}

// HandleRenew is implemented by cooldown.NopHandler
``` 

You can see more examples in `/examples` folder.

---


# Credits
[dragonfly](https://github.com/df-mc/dragonfly)

# License
[MIT](LICENSE)