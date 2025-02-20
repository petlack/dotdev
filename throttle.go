package main

import (
	"sync"
	"time"
)

// Throttle returns a version of fn that executes immediately on the first call,
// then at most once every interval. If additional calls occur during the throttle
// period, one trailing invocation will be executed when the interval expires.
func Throttle(fn func(), interval time.Duration) func() {
	var mu sync.Mutex
	var timer *time.Timer
	var trailing bool

	return func() {
		mu.Lock()
		defer mu.Unlock()
		if timer == nil {
			// Not throttled: execute immediately.
			fn()
			// Start a timer for the throttle interval.
			timer = time.AfterFunc(interval, func() {
				mu.Lock()
				if trailing {
					// A call happened during the throttle period: clear the flag and unlock.
					trailing = false
					mu.Unlock()
					// Execute the trailing call.
					fn()
					// Restart the timer so that if more calls come in during the next interval,
					// we can schedule another trailing invocation.
					mu.Lock()
					// It's safe to call Reset because timer is still active.
					timer.Reset(interval)
					mu.Unlock()
				} else {
					// No trailing call was requested; clear the timer so that the next call
					// will execute immediately.
					timer = nil
					mu.Unlock()
				}
			})
		} else {
			// Timer is active: mark that a trailing execution is needed.
			trailing = true
		}
	}
}
