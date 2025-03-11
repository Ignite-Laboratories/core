package impulse

import (
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/support/atomic"
	"sync"
	"time"
)

// Clock represents a source of time.  Once a clock is started, every Kernel
// associated with it will be invoked as frequently as possible in batched waves.
// Any Kernel that runs longer than a single tick is responsible for ensuring it
// doesn't run more frequently than it's invoked by the Clock.
// This is because kernels are meant to run -once- per impulse.
type Clock struct {
	Rate    int
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
	tickCount := 0
	tickCountStart := lastNow

	for core.KeepAlive {
		tickCount++

		var ctx Context
		ctx.Now = time.Now()
		ctx.Delta = ctx.Now.Sub(lastNow)
		ctx.Beat = beat
		ctx.Clock = c
		ctx.waitGroup = &wg

		// Calculate the current clock rate
		if tickCount > 1024 {
			elapsed := ctx.Now.Sub(tickCountStart).Seconds()
			c.Rate = int(float64(tickCount) / elapsed)
			tickCount = 0
			tickCountStart = ctx.Now
		}

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

func (c *Clock) EveryNthBeat(n int, action ActionFunc) {
	count := 0
	k := newActionPotential(func(ctx Context) bool {
		count++
		if count >= n {
			count = 0
			return true
		}
		return false
	}, action)
	c.AddKernel(k)
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
