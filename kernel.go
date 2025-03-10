package core

import "sync"

// Kernel is a program that can be invoked at a regular interval in a similar way to a shader.
// When invoked, it's provided the current 'beat' of the clock, allowing it to intelligently
// decide if it should handle this particular tick or not.  For instance, if it's still executing
// from a prior invocation it should ignore the request and let the clock keep ticking.
type Kernel interface {
	// Tick is called by the main clock for all ticks of the main execution loop.
	// The clock-provided beat indicates which point in the loop the tick is currently executing.
	// The wait group should be decremented by one as soon as possible.
	Tick(beat int, wg *sync.WaitGroup)
}
