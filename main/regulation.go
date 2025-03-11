package main

import (
	"fmt"
	"github.com/ignite-laboratories/core/impulse"
	"time"
)

var pulseDuration = time.Second / 100
var changeRate = 5000
var period = 100000
var clock = impulse.NewClock(period)

func main() {
	clock.OnDownbeat(Observe)
	clock.Start()
}

var lastNow = time.Now()

func Observe(ctx impulse.Context) {
	now := time.Now()
	lastDuration := now.Sub(lastNow)
	delta := lastDuration - pulseDuration
	absDelta := delta.Abs()

	if absDelta > time.Second/2 {
		changeRate = 50000
	} else if absDelta > time.Second/5 {
		changeRate = 40000
	} else if absDelta > time.Second/10 {
		changeRate = 20000
	} else if absDelta > time.Second/25 {
		changeRate = 1000
	} else if absDelta > time.Second/50 {
		changeRate = 100
	} else if absDelta > time.Second/100 {
		changeRate = 50
	} else if absDelta > time.Second/250 {
		changeRate = 25
	} else if absDelta > time.Second/500 {
		changeRate = 10
	} else if absDelta > time.Second/1000 {
		changeRate = 0
	}

	if delta < 0 {
		// The cycle was under the desired duration, and we must push the period out

		period += changeRate
	} else {
		// The cycle was over the desired duration, and we much pull the period in

		period -= changeRate
	}
	clock.Period = period

	fmt.Printf("%v | %v | %v\n", lastDuration, clock.Period, changeRate)
	lastNow = now
}
