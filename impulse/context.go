package impulse

import (
	"sync"
	"time"
)

// Context represents a set of contextually relevant information for a Kernel.
type Context struct {
	// Now is the currently processed moment in time for this impulse.
	Now time.Time
	// Delta is the amount of time that has passed since the last impulse.
	Delta time.Duration
	// Beat increments up to a fixed value defined by the Clock before looping back to 0.
	Beat int
	// Period is the upper limit the Beat will increment to.
	Period int
	// Kernel is an interface back to the originating Kernel.
	Kernel    Kernel
	waitGroup *sync.WaitGroup
}
