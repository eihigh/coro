package coro

import "iter"

type Coro struct {
	next func() (struct{}, bool)
	stop func()
}

func New(f func(Yield)) *Coro {
	seq := func(yield1 func(struct{}) bool) {
		f(func() {
			yield1(struct{}{})
		})
	}
	next, stop := iter.Pull(seq)
	return &Coro{next, stop}
}

func (c *Coro) Next() bool {
	_, ok := c.next()
	if !ok {
		c.stop()
	}
	return ok
}

type Yield func()

func (y Yield) Skip(n int) {
	for range n {
		y()
	}
}

func (y Yield) Forever() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; ; i++ {
			if !yield(i) {
				return
			}
			y()
		}
	}
}

func (y Yield) Loop(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range n {
			if !yield(i) {
				return
			}
			y()
		}
	}
}

func (y Yield) Until(f func() bool) {
	for !f() {
		y()
	}
}
