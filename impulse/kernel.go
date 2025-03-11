package impulse

// Kernel is a program that can be invoked at a regular interval by an execution Clock.
// When invoked it's provided some temporal context, allowing it to intelligently decide
// if it should execute or not.
//
// This is the foundation for action potential driven execution.
type Kernel interface {
	// Tick is called by the main Clock for all beats of the main execution loop.
	Tick(ctx Context)
	// GetID returns the Kernel identifier.
	GetID() uint64
	// IsExecuting returns whether the Kernel is currently executing.
	IsExecuting() bool
}
