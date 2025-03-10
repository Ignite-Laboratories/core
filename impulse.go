package core

// Impulse represents an executable entry point for a Kernel.
// While an Impulse doesn't handle any data directly, it motivates
// data through a neural network.
type Impulse struct {
	ID        uint64
	waiting   bool
	potential PotentialFunc
	action    ActionFunc
}

// NewImpulse initializes an Impulse struct with the provided potential and action.
// The potential function checks if this instance should or should not fire on this beat
// while the action function is called asynchronously if the potential returns true.
// The neuron will not invoke another action until the last completes.
func NewImpulse(potential PotentialFunc, action ActionFunc) *Impulse {
	return &Impulse{
		ID:        NextID(),
		potential: potential,
		action:    action,
	}
}

// GetID returns the ID of this Impulse.
func (imp *Impulse) GetID() uint64 {
	return imp.ID
}

// Tick is called by a Clock for every beat.  The provided WaitGroup is decremented
// once the Impulse has finished executing (or not).  It calls the neuron's potential
// function prior to asynchronously calling the neuron's action function, if the
// potential is true.  Once an action is initiated, the neuron will not invoke another
// action until the last action completes.
func (imp *Impulse) Tick(ctx Context) {
	if !imp.waiting && imp.potential(ctx) {
		imp.waiting = true
		go func() {
			imp.action(ctx)
			imp.waiting = false
		}()
	}
	ctx.waitGroup.Done()
}
