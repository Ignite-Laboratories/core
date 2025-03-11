package main

import (
	"fmt"
	"github.com/ignite-laboratories/core/impulse"
	"sync"
	"time"
)

type TemporalFrame struct {
	lastNow      time.Time
	lastDuration time.Duration
	delta        time.Duration
	absDelta     time.Duration
}

func NewTemporalFrame(now time.Time, duration time.Duration) TemporalFrame {
	return TemporalFrame{
		lastNow:      now,
		lastDuration: duration,
		delta:        duration - pulseDuration,
		absDelta:     duration.Abs(),
	}
}

var pulseDuration = time.Second / 100
var deviation = 5000
var period = 100000
var clock = impulse.NewClock(period)
var timeline sync.Map

func main() {
	clock.Activate.OnDownbeat(Regulate)
	clock.Activate.OnDownbeat(ObserveAndPrint)
	clock.Activate.OnDownbeat(Observe)
	clock.Activate.OnDownbeat(Observe)
	clock.Activate.OnDownbeat(Observe)
	clock.Activate.OnDownbeat(Observe)
	clock.Activate.OnDownbeat(Observe)
	clock.Activate.OnDownbeat(Observe)
	clock.Activate.OnDownbeat(Observe)
	clock.Activate.OnDownbeat(Observe)
	clock.Activate.OnDownbeat(Observe)
	clock.Activate.OnDownbeat(Observe)
	clock.Start()
}

func Regulate(ctx impulse.Context) {
	// Setup our calculation
	avg := int64(0)
	count := 0

	// Calculate the average of all the temporal frames' deltas
	timeline.Range(func(key, value interface{}) bool {
		count++
		frame, _ := value.(TemporalFrame)
		avg += frame.delta.Nanoseconds()
		return true
	})
	if count == 0 {
		count = 1
	}
	avg /= int64(count)
	delta := time.Duration(avg)
	absDelta := delta.Abs()

	// Calculate how much the rate should change
	if absDelta > time.Second/2 {
		deviation = 50000
	} else if absDelta > time.Second/5 {
		deviation = 40000
	} else if absDelta > time.Second/10 {
		deviation = 20000
	} else if absDelta > time.Second/25 {
		deviation = 1000
	} else if absDelta > time.Second/50 {
		deviation = 100
	} else if absDelta > time.Second/100 {
		deviation = 50
	} else if absDelta > time.Second/250 {
		deviation = 25
	} else if absDelta > time.Second/500 {
		deviation = 10
	} else if absDelta > time.Second/1000 {
		deviation = 0
	}

	// Adjust the clock's period
	if delta < 0 {
		// The cycle was under the desired duration, and we must push the period out
		period += deviation
	} else {
		// The cycle was over the desired duration, and we much pull the period in
		period -= deviation
	}
	clock.Period = period
}

func ObserveAndPrint(ctx impulse.Context) {
	observe(ctx, true)
}

func Observe(ctx impulse.Context) {
	observe(ctx, false)
}

func observe(ctx impulse.Context, print bool) {
	// Get the local 'now' to observe
	now := time.Now()

	// Build the new temporal frame
	var frame TemporalFrame
	var delta time.Duration
	value, ok := timeline.Load(ctx.Kernel.GetID())
	if ok {
		// If a frame existed, calculate the delta between frames.
		lastFrame, _ := value.(TemporalFrame)
		delta = now.Sub(lastFrame.lastNow)
	}
	frame = NewTemporalFrame(now, delta)

	// Print out that this kernel did something
	if print {
		fmt.Printf("%v | %v | %v\n", frame.lastDuration, clock.Period, deviation)
	}

	// Save off the new temporal context
	frame.lastNow = now
	timeline.Store(ctx.Kernel.GetID(), frame)
}
