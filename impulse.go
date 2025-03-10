package core

import "sync"

// Impulse represents an executable entry point for a Kernel.
// While an Impulse doesn't handle any data directly, it motivates
// data through a neural network.
type Impulse struct {
	waiting   bool
	potential func(beat int) bool
	action    func(beat int)
}

// NewImpulse initializes an Impulse struct with the provided potential and action.
// The potential function checks if this instance should or should not fire on this beat
// while the action function is called asynchronously if the potential returns true.
// The neuron will not invoke another action until the last completes.
func NewImpulse(potential func(beat int) bool, action func(beat int)) *Impulse {
	return &Impulse{
		potential: potential,
		action:    action,
	}
}

// Tick is called by a Clock for every beat.  The provided WaitGroup is decremented
// once the Impulse has finished executing (or not).  It calls the neuron's potential
// function prior to asynchronously calling the neuron's action function, if the
// potential is true.  Once an action is initiated, the neuron will not invoke another
// action until the last action completes.
func (n *Impulse) Tick(beat int, wg *sync.WaitGroup) {
	if !n.waiting && n.potential(beat) {
		n.waiting = true
		go func() {
			n.action(beat)
			n.waiting = false
		}()
	}
	wg.Done()
}
