package main

import (
	"testing"
	"time"
)

// TestThrottleNoTrailing verifies that if no extra calls occur during the throttle interval,
// only the immediate call is executed.
func TestThrottleNoTrailing(t *testing.T) {
	var counter int
	fn := func() {
		counter++
	}
	interval := 100 * time.Millisecond
	throttled := Throttle(fn, interval)
	throttled() // counter == 1
	if counter != 1 {
		t.Errorf("expected counter to be 1, got %d", counter)
	}
	time.Sleep(150 * time.Millisecond)
	if counter != 1 {
		t.Errorf("expected counter to be 1, got %d", counter)
	}
}

// TestThrottleWithTrailing verifies that if calls occur during the throttle interval,
// a trailing call is executed after the interval.
func TestThrottleWithTrailing(t *testing.T) {
	var counter int
	fn := func() {
		counter++
	}
	interval := 100 * time.Millisecond
	throttled := Throttle(fn, interval)

	throttled() // counter == 1
	// Subsequent calls during the throttle period.
	throttled()
	throttled()

	if counter != 1 {
		t.Errorf("expected counter to be 2 (immediate + trailing), got %d", counter)
	}
	time.Sleep(150 * time.Millisecond)
	if counter != 2 {
		t.Errorf("expected counter to be 2 (immediate + trailing), got %d", counter)
	}
}

// TestThrottleMultipleTrailing simulates continuous calls so that each throttle interval
// should produce a trailing call.
func TestThrottleMultipleTrailing(t *testing.T) {
	var counter int
	fn := func() {
		counter++
	}
	interval := 100 * time.Millisecond
	throttled := Throttle(fn, interval)

	throttled() // counter == 1

	start := time.Now()
	time.Sleep(50 * time.Millisecond)
	for time.Since(start) < 300*time.Millisecond {
		throttled()
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(350 * time.Millisecond)

	// Expect at least one trailing call, possibly more.
	if counter < 2 {
		t.Errorf("expected at least 2 calls, got %d", counter)
	}
}

// TestThrottleReset verifies that once the throttle interval has passed without extra calls,
// a new call executes immediately.
func TestThrottleReset(t *testing.T) {
	var counter int
	fn := func() {
		counter++
	}
	interval := 100 * time.Millisecond
	throttled := Throttle(fn, interval)

	throttled() // counter == 1
	throttled()
	time.Sleep(250 * time.Millisecond)
	throttled() // counter == 3

	if counter != 3 {
		t.Errorf("expected counter to be 3, got %d", counter)
	}
}
