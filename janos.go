package core

import (
	"sync"
	"sync/atomic"
	"time"
)

// KeepAlive is used by all clocks to keep ticking - set it to false in order to terminate the beat.
var KeepAlive = true
var masterId uint64

// NextID provides a unique identifier to every function that calls it.
func NextID() uint64 {
	return atomic.AddUint64(&masterId, 1)
}

// PotentialFunc represents a conditional test to perform an action.
type PotentialFunc func(ctx Context) bool

// ActionFunc represents a performable action.
type ActionFunc func(ctx Context)

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
