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
type actionPotential struct {
	// ID is the unique entity identifier for this actionPotential.
	ID uint64
	// lastBeatMoment is the last beat's moment from which type of Kernel was activated.
	lastBeatMoment time.Time
	// lastCompletion is the last moment in time type of Kernel finished execution.
	lastCompletion time.Time
	// executing is true if this actionPotential is currently activated.
	executing bool
	// action is a function that is invoked if the PotentialFunc returns true.
	action ActivationFunc
	// potential is a function that determines if the ActivationFunc should be invoked.
	potential PotentialFunc
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

// Execute is called by a Clock for every beat.  It calls the kernel's PotentialFunc
// prior to asynchronously calling the kernel's ActivationFunc, if the potential
// function returns true.  Once the ActivationFunc starts executing, the system
// will not invoke it again until the current action completes. The provided WaitGroup
// is decremented once the PotentialFunc returns, regardless of activation.
func (ap *actionPotential) Execute(ctx Context) {
	if !ap.executing && ap.potential(ctx) {
		ap.executing = true
		if !ap.lastBeatMoment.IsZero() {
			ctx.Delta = ctx.Moment.Sub(ap.lastBeatMoment)
			ctx.LastDuration = ap.lastCompletion.Sub(ap.lastBeatMoment)
		}
		go func() {
			ap.action(ctx)
			ap.lastCompletion = time.Now()
			ap.executing = false
		}()
		ap.lastBeatMoment = ctx.Moment
	}
	ctx.waitGroup.Done()
}
