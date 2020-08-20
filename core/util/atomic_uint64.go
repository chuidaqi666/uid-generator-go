package util

import "sync/atomic"

func Uint64UpdateAndGet(v *int64, f func(int64) int64) int64 {
	var old, next int64
	for {
		old = atomic.LoadInt64(v)
		next = f(old)
		if atomic.CompareAndSwapInt64(v, old, next) {
			break
		}
	}
	return next
}
