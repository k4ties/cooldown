# cooldown
A [Go](https://go.dev/) library implementing cooldown.

### Get started

```go
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/k4ties/cooldown"
)

func main() {
	cd := cooldown.New(cooldown.OptionHandler(new(Handler)))
	cd.Start(time.Second)
	cd.Active() // true

	<-time.After(time.Millisecond)
	cd.Renew()

	<-time.After(time.Second * 3)
	cd.Stop()
}

type Handler struct {
	cooldown.NopHandler // implements cooldown.Handler
	first               sync.Once
}

func (h *Handler) HandleStart(ctx *cooldown.Context, _ time.Duration) {
	h.first.Do(func() {
		// The event can be canceled
		ctx.Cancel()
	})
	fmt.Println("handle start")
}
func (h *Handler) HandleStop(_ *cooldown.CoolDown, cause cooldown.StopCause) {
	// This event cannot be canceled
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