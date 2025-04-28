// Package coro makes Go coroutines simple and useful.
//
// Go 1.23 introduced coroutines as pull-style transformations of the standard push-style
// iterators. Go coroutines are known as stackful coroutines, which are concurrent like
// goroutines but not parallel. For more details, please read Russ Cox's article
// [Coroutines for Go].
//
// [Coroutines for Go]: https://research.swtch.com/coro
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
// Returns true if the coroutine should continue execution,
// or false if the coroutine is stopped.
func (y Yield) Skip(n int) bool {
	for range n {
		if !y() {
			return false
		}
	}
	return true
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
// Returns true if the coroutine should continue execution,
// or false if the coroutine is stopped.
func (y Yield) Until(f func() bool, preYield bool) bool {
	if preYield {
		if !y() {
			return false
		}
	}
	for !f() {
		if !y() {
			return false
		}
	}
	return true
}

// While repeatedly yields while the condition function f returns true.
//
// If preYield is true, it yields once before checking the condition.
//
// Returns true if the coroutine should continue execution,
// or false if the coroutine is stopped.
func (y Yield) While(f func() bool, preYield bool) bool {
	if preYield {
		if !y() {
			return false
		}
	}
	for f() {
		if !y() {
			return false
		}
	}
	return true
}
