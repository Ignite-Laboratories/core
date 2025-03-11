package impulse

import (
	"github.com/ignite-laboratories/core"
	"time"
)

// PotentialFunc represents a conditional test of whether to perform an action or not.
type PotentialFunc func(ctx Context) bool

// ActivationFunc represents a executable action.
type ActivationFunc func(ctx Context)

// actionPotential represents an executable entry point for a Kernel.
// While an actionPotential doesn't handle any data directly, it motivates
// data through a neural network.
type actionPotential struct {
	// ID is the unique entity identifier for this actionPotential.
	ID uint64
	// lastTrigger is the last moment in Clock-time this actionPotential was activated.
	lastTrigger    time.Time
	lastCompletion time.Time
	// executing is true if this actionPotential is currently activated.
	executing bool
	// potential is a function that determines if the ActivationFunc should be invoked.
	potential PotentialFunc
	// action is a function that is invoked if the PotentialFunc returns true.
	action ActivationFunc
}

// newActionPotential initializes a Kernel with the provided potential and action.
// The potential function checks if this instance should or should not fire on this beat
// while the action function is called asynchronously if the potential returns true.
// The neuron will not invoke another action until the last completes.
func newActionPotential(potential PotentialFunc, action ActivationFunc) Kernel {
	return &actionPotential{
		ID:        core.NextID(),
		potential: potential,
		action:    action,
	}
}

// GetID returns the ID of this actionPotential.
func (ap *actionPotential) GetID() uint64 {
	return ap.ID
}

// IsExecuting returns whether the ActivationFunc is currently executing.
func (ap *actionPotential) IsExecuting() bool {
	return ap.executing
}

// Tick is called by a Clock for every beat.  The provided WaitGroup is decremented
// once the actionPotential has finished executing (or not).  It calls the neuron's potential
// function prior to asynchronously calling the neuron's action function, if the
// potential is true.  Once an action is initiated, the neuron will not invoke another
// action until the last action completes.
func (ap *actionPotential) Tick(ctx Context) {
	if !ap.executing && ap.potential(ctx) {
		ap.executing = true
		if !ap.lastTrigger.IsZero() {
			ctx.Delta = ctx.Moment.Sub(ap.lastTrigger)
			ctx.LastExecution = ap.lastCompletion.Sub(ap.lastTrigger)
		}
		ap.lastTrigger = ctx.Moment
		go func() {
			ap.action(ctx)
			ap.lastCompletion = time.Now()
			ap.executing = false
		}()
	}
	ctx.waitGroup.Done()
}
