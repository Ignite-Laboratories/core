package beat

import "github.com/ignite-laboratories/core"

// All returns true for every beat.
func All(ctx core.Context) bool {
	return true
}

// Downbeat returns true for beat 0 only.
func Downbeat(ctx core.Context) bool {
	return ctx.Beat == 0
}

// Even returns true for any even beat number.
func Even(ctx core.Context) bool {
	return ctx.Beat%2 == 0
}

// Odd returns true for any odd beat number.
func Odd(ctx core.Context) bool {
	return ctx.Beat%2 != 0
}
