package core

import "sync/atomic"

// KeepAlive is used by all clocks to keep ticking - set it to false in order to terminate the beat.
var KeepAlive = true
var masterId uint64

// NextID provides a unique identifier to every function that calls it.
func NextID() uint64 {
	return atomic.AddUint64(&masterId, 1)
}
