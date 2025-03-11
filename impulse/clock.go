package impulse

import (
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/support/atomic"
	"sync"
	"time"
)

// Clock represents a source of time.  Once a impulse is started, all kernels
// associated with it will be invoked as frequently as possible in batched waves.
// Kernels that run longer than a single tick are responsible for ensuring they
// don't run more frequently than they are invoked by the impulse.
type Clock struct {
	period  int
	kernels *atomic.Slice[Kernel]
}

// NewClock creates a new instance of a Clock.  The provided period defines how
// high to increment the current beat before looping back to 0.  If your impulse
// should loop between 0-31, provide a period of 32 beats.
func NewClock(period int) Clock {
	return Clock{
		period:  period,
		kernels: atomic.NewSlice[Kernel](),
	}
}

// AddKernel adds a new Kernel to the Clock.
func (c *Clock) AddKernel(kernel Kernel) {
	c.kernels.Add(kernel)
}

// RemoveKernel removes a Kernel from the Clock by its ID.
func (c *Clock) RemoveKernel(kernel Kernel) {
	id := kernel.GetID()
	c.kernels.RemoveIf(func(k Kernel) bool {
		return k.GetID() == id
	})
}

// Start is the entry point to begin ticking.
func (c *Clock) Start() {
	var wg sync.WaitGroup
	beat := 0
	lastNow := time.Now()

	for core.KeepAlive {
		var ctx Context
		ctx.Now = time.Now()
		ctx.Delta = ctx.Now.Sub(lastNow)
		ctx.Beat = beat
		ctx.Period = c.period
		ctx.waitGroup = &wg

		// We retrieve all kernels first in case the
		// data changes during this loop cycle.
		ks := c.kernels.All()
		wg.Add(len(ks))
		for _, k := range ks {
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

// On creates a Kernel that fires the provided action whenever the invoked potential function returns true.
func (c *Clock) On(potential PotentialFunc, action ActionFunc) {
	k := newActionPotential(potential, action)
	c.AddKernel(k)
}

// OnAllBeats creates a Kernel that fires the provided action on every beat.
func (c *Clock) OnAllBeats(action ActionFunc) {
	k := newActionPotential(func(ctx Context) bool {
		return true
	}, action)
	c.AddKernel(k)
}

// OnOddBeats creates a Kernel that fires the provided action on odd beats.
func (c *Clock) OnOddBeats(action ActionFunc) {
	k := newActionPotential(func(ctx Context) bool {
		return ctx.Beat%2 != 0
	}, action)
	c.AddKernel(k)
}

// OnEvenBeats creates a Kernel that fires the provided action on even beats.
func (c *Clock) OnEvenBeats(action ActionFunc) {
	k := newActionPotential(func(ctx Context) bool {
		return ctx.Beat%2 == 0
	}, action)
	c.AddKernel(k)
}

// OnDownbeat creates a Kernel that fires the provided action on beat 0
func (c *Clock) OnDownbeat(action ActionFunc) {
	c.OnBeat(0, action)
}

// OnBeat creates a Kernel that fires the provided action on the specified beat
func (c *Clock) OnBeat(beat int, action ActionFunc) {
	k := newActionPotential(func(ctx Context) bool {
		return ctx.Beat == beat
	}, action)
	c.AddKernel(k)
}
