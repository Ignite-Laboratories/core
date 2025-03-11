package impulse

import (
	"fmt"
	"sync"
	"time"
)

// Context represents a set of contextually relevant information for a Kernel.
type Context struct {
	// Moment is the moment in time the Clock itself observed for this beat.
	Moment time.Time
	// Delta is the amount of time that has passed since the Kernel's last impulse.
	Delta time.Duration
	// LastDuration is the amount of time that the Kernel last took to finish execution.
	LastDuration time.Duration
	// Beat is the current beat number the clock is on.
	Beat int
	// Clock is a reference back to the source Clock.
	Clock *Clock
	// Kernel is an interface back to the originating Kernel.
	Kernel Kernel
	// waitGroup is decremented immediately after the PotentialFunc is called regardless of activation.
	waitGroup *sync.WaitGroup
}

func (ctx Context) String() string {
	if ctx.Delta == 0 {
		return fmt.Sprintf("[Pulse %d.%d] Kernel %d activated", ctx.Clock.ID, ctx.Beat, ctx.Kernel.GetID())
	}
	return fmt.Sprintf("[Pulse %d.%d] Kernel %d activated (Δt %s) (Execution time %s)", ctx.Clock.ID, ctx.Beat, ctx.Kernel.GetID(), ctx.Delta, ctx.LastDuration)
}
