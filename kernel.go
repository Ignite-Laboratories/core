package core

// Kernel is a program that can be invoked at a regular interval in a similar way to a shader.
// When invoked, it's provided the current 'beat' of the Clock, allowing it to intelligently
// decide if it should handle this particular tick or not.  For instance, if it's still executing
// from a prior invocation it should ignore the request and let the clock keep ticking.
type Kernel interface {
	// Tick is called by the main Clock for all beats of the main execution loop.
	Tick(ctx Context)
	// GetID returns the Kernel identifier.
	GetID() uint64
}

// actionPotential represents an executable entry point for a Kernel.
// While an actionPotential doesn't handle any data directly, it motivates
// data through a neural network.
type actionPotential struct {
	ID        uint64
	waiting   bool
	potential PotentialFunc
	action    ActionFunc
}

// NewKernel initializes a Kernel with the provided potential and action.
// The potential function checks if this instance should or should not fire on this beat
// while the action function is called asynchronously if the potential returns true.
// The neuron will not invoke another action until the last completes.
func NewKernel(potential PotentialFunc, action ActionFunc) Kernel {
	return &actionPotential{
		ID:        NextID(),
		potential: potential,
		action:    action,
	}
}

// GetID returns the ID of this actionPotential.
func (imp *actionPotential) GetID() uint64 {
	return imp.ID
}

// Tick is called by a Clock for every beat.  The provided WaitGroup is decremented
// once the actionPotential has finished executing (or not).  It calls the neuron's potential
// function prior to asynchronously calling the neuron's action function, if the
// potential is true.  Once an action is initiated, the neuron will not invoke another
// action until the last action completes.
func (imp *actionPotential) Tick(ctx Context) {
	if !imp.waiting && imp.potential(ctx) {
		imp.waiting = true
		go func() {
			imp.action(ctx)
			imp.waiting = false
		}()
	}
	ctx.waitGroup.Done()
}
