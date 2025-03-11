package impulse

import (
	"github.com/ignite-laboratories/core"
	"github.com/ignite-laboratories/support/atomic"
	"sync"
	"time"
)

// Clock represents a source of time.  Once a clock is started, every Kernel
// associated with it will be tested for invocation as frequently as possible.
// If it's PotentialFunc returns true for that beat of the Clock, it's ActivationFunc
// is invoked. If the that ActivationFunc is still executing and the PotentialFunc
// returns true, the impulse is ignored - ensuring each kernel only executes serially.
type Clock struct {
	// ID is the unique entity identifier for this Clock.
	ID uint64
	// ClockRate represents the current number of beats per second - it is a read-only value.
	ClockRate int
	// Period represents the number this clock counts to before looping - it is a writable value.
	Period int
	// Activate provides a selection of helper methods for Kernel inception.
	Activate activation
	// kernels stores the currently executing kernels on this clock.
	kernels *atomic.Slice[Kernel]
}

// activation provides a selection of helper methods for Kernel inception.
type activation struct {
	// clock provides a reference back to the source Clock.
	clock *Clock
}

// NewClock creates a new Clock with a specified 'period' that determines how
// many beats the clock processes before looping back to 0.
func NewClock(period int) Clock {
	clock := &Clock{
		ID:      core.NextID(),
		Period:  period,
		kernels: atomic.NewSlice[Kernel](),
	}
	clock.Activate = activation{clock: clock}
	return *clock
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
		ctx.Moment = time.Now()
		ctx.Delta = time.Duration(0)
		ctx.Beat = beat
		ctx.Clock = c
		ctx.waitGroup = &wg

		// Calculate the current clock rate
		if tickCount > 1024 {
			elapsed := ctx.Moment.Sub(tickCountStart).Seconds()
			c.ClockRate = int(float64(tickCount) / elapsed)
			tickCount = 0
			tickCountStart = ctx.Moment
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
		if beat >= c.Period {
			beat = 0
		}
		lastNow = ctx.Moment
	}
}

// EveryNthBeat creates a Kernel that fires the provided action every n beats, regardless
// of the clock's current beat value.
func (a *activation) EveryNthBeat(n int, action ActivationFunc) {
	count := 0
	k := newActionPotential(func(ctx Context) bool {
		// Only count while the kernel is not executing
		if ctx.Kernel.IsExecuting() {
			return false
		}

		count++
		if count >= n {
			count = 0
			return true
		}
		return false
	}, action)
	a.clock.AddKernel(k)
}

// OnCondition creates a Kernel that fires the provided action whenever the invoked potential function returns true.
func (a *activation) OnCondition(potential PotentialFunc, action ActivationFunc) {
	k := newActionPotential(potential, action)
	a.clock.AddKernel(k)
}

// EveryBeat creates a Kernel that fires the provided action on every beat.
func (a *activation) EveryBeat(action ActivationFunc) {
	k := newActionPotential(func(ctx Context) bool {
		return true
	}, action)
	a.clock.AddKernel(k)
}

// OnOddBeats creates a Kernel that fires the provided action on odd beats.
func (a *activation) OnOddBeats(action ActivationFunc) {
	k := newActionPotential(func(ctx Context) bool {
		return ctx.Beat%2 != 0
	}, action)
	a.clock.AddKernel(k)
}

// OnEvenBeats creates a Kernel that fires the provided action on even beats.
func (a *activation) OnEvenBeats(action ActivationFunc) {
	k := newActionPotential(func(ctx Context) bool {
		return ctx.Beat%2 == 0
	}, action)
	a.clock.AddKernel(k)
}

// OnDownbeat creates a Kernel that fires the provided action on beat 0
func (a *activation) OnDownbeat(action ActivationFunc) {
	a.clock.Activate.OnBeatNumber(0, action)
}

// OnBeatNumber creates a Kernel that fires the provided action on the specified beat
func (a *activation) OnBeatNumber(beat int, action ActivationFunc) {
	k := newActionPotential(func(ctx Context) bool {
		return ctx.Beat == beat
	}, action)
	a.clock.AddKernel(k)
}
