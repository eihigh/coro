// Package coro provides a lightweight coroutine abstraction for Go.
// It builds upon Go's standard iter package (introduced in Go 1.22) to offer
// a more intuitive API for creating pausable functions that can yield control
// and resume execution, similar to coroutines in other languages.
//
// While Go's iter package already provides the fundamental building blocks for
// iterator-based control flow, this package simplifies the creation and management
// of coroutine-like patterns with a more familiar interface.
package coro

import "iter"

// Coro represents a coroutine that can be paused and resumed.
//
// It internally uses iter.Pull to manage the state of the coroutine.
type Coro struct {
	next func() (struct{}, bool)
	stop func()
}

// New creates a new coroutine from the given function.
//
// The function f receives a Yield function that can be called to pause
// the coroutine and return control to the caller.
func New(f func(Yield)) *Coro {
	seq := func(yield1 func(struct{}) bool) {
		f(func() bool {
			return yield1(struct{}{})
		})
	}
	next, stop := iter.Pull(seq)
	return &Coro{next, stop}
}

// Next advances the coroutine to its next yield point.
//
// Returns true if the coroutine yielded and can be resumed,
// or false if the coroutine has completed execution.
//
// When the coroutine completes, it automatically calls stop.
func (c *Coro) Next() bool {
	_, ok := c.next()
	if !ok {
		c.stop()
	}
	return ok
}

// Stop explicitly terminates the coroutine.
//
// This can be called to clean up resources when the coroutine is no longer needed.
func (c *Coro) Stop() {
	c.stop()
}

// Yield is a function that pauses the coroutine when called.
//
// Returns true if the coroutine should continue execution,
// or false if the coroutine should terminate.
type Yield func() (ok bool)

// Skip advances the coroutine by yielding n times.
//
// Stops early if the coroutine is stopped.
func (y Yield) Skip(n int) {
	for range n {
		if !y() {
			return
		}
	}
}

// Seq returns an iterator sequence of integers from 0 to n-1,
// yielding after each number is produced.
//
// The sequence stops early if either the for-range loop is
// exited (using break) or the coroutine is stopped.
func (y Yield) Seq(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range n {
			if !yield(i) || !y() {
				return
			}
		}
	}
}

// Until repeatedly yields until the condition function f returns true.
//
// If preYield is true, it yields once before checking the condition.
//
// Stops early if the coroutine is stopped.
func (y Yield) Until(f func() bool, preYield bool) {
	if preYield {
		if !y() {
			return
		}
	}
	for !f() {
		if !y() {
			return
		}
	}
}

// While repeatedly yields while the condition function f returns true.
//
// If preYield is true, it yields once before checking the condition.
//
// Stops early if the coroutine is stopped.
func (y Yield) While(f func() bool, preYield bool) {
	if preYield {
		if !y() {
			return
		}
	}
	for f() {
		if !y() {
			return
		}
	}
}
