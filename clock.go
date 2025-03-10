package core

import "sync"

// Clock represents a source of time.  Once a clock is started, all neurons
// associated with it will be invoked as frequently as possible in batched waves.
// Neurons that run longer than a single tick are responsible for ensuring they
// don't run more frequently than they are invoked by the clock.
type Clock struct {
	Neurons []*Neuron
}

// Start is the entry point to begin ticking. The provided period defines how
// high to increment the current beat before looping back to 0.  If your clock
// should loop between 0-31, provide a period of 32 beats.
func (c Clock) Start(period int) {
	beat := 0
	var wg sync.WaitGroup

	for KeepAlive {
		wg.Add(len(c.Neurons))
		for _, n := range c.Neurons {
			go (*n).Tick(beat, &wg)
		}
		wg.Wait()

		beat++
		if beat == period {
			beat = 0
		}
	}
}
