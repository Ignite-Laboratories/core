package impulse

import (
	"fmt"
	"sync"
	"time"
)

// Context represents a set of contextually relevant information for a Kernel.
type Context struct {
	// Now is the currently processed moment in time for this impulse.
	Now time.Time
	// Delta is the amount of time that has passed since the Kernel's last impulse.
	Delta time.Duration
	// LastExecution is the amount of time that the Kernel last took to finish execution.
	LastExecution time.Duration
	// Beat increments up to a fixed value defined by the Clock before looping back to 0.
	Beat int
	// Clock is a reference back to the source Clock.
	Clock *Clock
	// Kernel is an interface back to the originating Kernel.
	Kernel    Kernel
	waitGroup *sync.WaitGroup
}

func (ctx Context) String() string {
	if ctx.Delta == 0 {
		return fmt.Sprintf("[Pulse %d.%d] Kernel %d activated", ctx.Clock.ID, ctx.Beat, ctx.Kernel.GetID())
	}
	return fmt.Sprintf("[Pulse %d.%d] Kernel %d activated (Δt %s) (Execution time %s)", ctx.Clock.ID, ctx.Beat, ctx.Kernel.GetID(), ctx.Delta, ctx.LastExecution)
}
