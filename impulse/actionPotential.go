package impulse

import "github.com/ignite-laboratories/core"

// PotentialFunc represents a conditional test to perform an action.
type PotentialFunc func(ctx Context) bool

// ActionFunc represents a performable action.
type ActionFunc func(ctx Context)

// actionPotential represents an executable entry point for a Kernel.
// While an actionPotential doesn't handle any data directly, it motivates
// data through a neural network.
type actionPotential struct {
	ID        uint64
	waiting   bool
	potential PotentialFunc
	action    ActionFunc
}

// newActionPotential initializes a Kernel with the provided potential and action.
// The potential function checks if this instance should or should not fire on this beat
// while the action function is called asynchronously if the potential returns true.
// The neuron will not invoke another action until the last completes.
func newActionPotential(potential PotentialFunc, action ActionFunc) Kernel {
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

// Tick is called by a Clock for every beat.  The provided WaitGroup is decremented
// once the actionPotential has finished executing (or not).  It calls the neuron's potential
// function prior to asynchronously calling the neuron's action function, if the
// potential is true.  Once an action is initiated, the neuron will not invoke another
// action until the last action completes.
func (ap *actionPotential) Tick(ctx Context) {
	if !ap.waiting && ap.potential(ctx) {
		ap.waiting = true
		go func() {
			ap.action(ctx)
			ap.waiting = false
		}()
	}
	ctx.waitGroup.Done()
}
