package core

import "sync/atomic"

// Alive globally keeps any long-running routine alive until it is set to false.
var Alive = true

// masterId holds the last provided entity identifier value.
var masterId uint64

// NextID provides a thread-safe unique entity identifier to every caller.
func NextID() uint64 {
	return atomic.AddUint64(&masterId, 1)
}
