package impulse

// Kernel is a program that can be invoked at a regular interval in a similar way to a shader.
// When invoked, it's provided the current 'beat' of the Clock, allowing it to intelligently
// decide if it should handle this particular tick or not.  For instance, if it's still executing
// from a prior invocation it should ignore the request and let the impulse keep ticking.
type Kernel interface {
	// Tick is called by the main Clock for all beats of the main execution loop.
	Tick(ctx Context)
	// GetID returns the Kernel identifier.
	GetID() uint64
}
