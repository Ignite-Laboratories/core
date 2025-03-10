package core

import (
	"github.com/ignite-laboratories/support/atomic"
	"sync"
	"time"
)

// Clock represents a source of time.  Once a clock is started, all kernels
// associated with it will be invoked as frequently as possible in batched waves.
// Kernels that run longer than a single tick are responsible for ensuring they
// don't run more frequently than they are invoked by the clock.
type Clock struct {
	period  int
	Kernels *atomic.Slice[Kernel]
}

// NewClock creates a new instance of a Clock.  The provided period defines how
// high to increment the current beat before looping back to 0.  If your clock
// should loop between 0-31, provide a period of 32 beats.
func NewClock(period int) Clock {
	return Clock{
		period:  period,
		Kernels: atomic.NewSlice[Kernel](),
	}
}

// Start is the entry point to begin ticking.
func (c Clock) Start() {
	var wg sync.WaitGroup
	beat := 0
	lastNow := time.Now()

	for KeepAlive {
		var ctx Context
		ctx.Now = time.Now()
		ctx.Delta = ctx.Now.Sub(lastNow)
		ctx.Beat = beat
		ctx.Period = c.period
		ctx.waitGroup = &wg

		// We retrieve all kernels first in case the
		// data changes during this loop cycle.
		kernels := c.Kernels.All()
		wg.Add(len(kernels))
		for _, k := range kernels {
			ctx.Kernel = k
			go k.Tick(ctx)
		}
		wg.Wait()

		beat++
		if beat >= c.period {
			beat = 0
		}
		lastNow = ctx.Now
	}
}
