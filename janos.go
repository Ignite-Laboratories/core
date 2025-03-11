package core

import (
	"sync/atomic"
)

// KeepAlive globally keeps JanOS's persistent loops alive until it is set to false.
var KeepAlive = true

// masterId is the currently unique entity identifier value.
var masterId uint64

// NextID provides a unique entity identifier to every function that calls it.
func NextID() uint64 {
	return atomic.AddUint64(&masterId, 1)
}
