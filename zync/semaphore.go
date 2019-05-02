package zync

import (
	"container/list"
	"sync"
)

// Reference
// https://godoc.org/golang.org/x/sync/semaphore

type waiter struct {
	size uint
	done chan interface{}
}

// Semaphore is a concurrency access control data structure that allows
// multiple routines to acquire access to a resource. A Semaphore has a max
// size; whenever a routine successfully acquires the lock and resource, a
// counter gets incremented. In order to acquire the lock, the routine calls
// `Acquire` or `TryAcquire` with a `delta` value; if the delta value fits
// within the Semaphores `limit` and the Semaphore has enough space to fit
// `delta` (this is calculated by checking `if delta <= sem.limit-sem.count`).
type Semaphore struct {
	lock    sync.RWMutex
	limit   uint
	count   uint
	waiters *list.List
}

// NewSemaphore will create a new empty Semaphore with the given max limit.
func NewSemaphore(limit uint) *Semaphore {
	s := new(Semaphore)
	s.limit = limit
	s.count = 0
	s.waiters = list.New()
	return s
}

// Acquire will block until it can acquire the lock and access the resource.
// If `delta > sem.limit`, the method will not block but instead return
// immediately without changing anything inside the Semaphore.
func (s *Semaphore) Acquire(delta uint) {
	// concurrency lock
	s.lock.Lock()

	// check if delta is more than cap
	if delta > s.limit {
		s.lock.Unlock()
		return
	}

	// check if delta fits and assign it
	// also check for other routines waiting
	if s.limit-s.count >= delta && s.waiters.Len() == 0 {
		s.count += delta
		s.lock.Unlock()
		return
	}

	// join waiting list
	release := make(chan interface{})
	w := waiter{size: delta, done: release}
	s.waiters.PushBack(w)
	s.lock.Unlock()

	// wait for release
	select {
	case <-release:
	}
}

// TryAcquire will attempt to acquire the lock on the semaphore and return
// `true` if successful. If `delta` is larger than the Semaphore `limit` or the
// Semaphore does not have enough space to allocate `delta`, then `false` is
// returned.
func (s *Semaphore) TryAcquire(delta uint) bool {
	// concurrency lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// check if delta is more than cap
	if delta > s.limit {
		return false
	}

	// check if delta fits and assign it
	// also check for other routines waiting
	if s.limit-s.count >= delta && s.waiters.Len() == 0 {
		s.count += delta
		return true
	}

	// fail to acquire
	return false
}

// Release will free up the given number of `delta` from the Semaphore. This
// method will return without changing anything if `delta > sem.limit`. Also
// this method will panic if too much `delta` is released from the Semaphore.
// At the end, Release will check for any waiting Acquire callers, check the
// `delta` of the first, and assign as many waiters as possible.
//
// Release/Acquire policy is first-come first-served. This means if we have
// space for one waiter that is further in the list, it will have to wait its
// turn, this method is preferred when it comes to not starving large Acquires.
func (s *Semaphore) Release(delta uint) {
	// concurrency lock
	s.lock.Lock()
	defer s.lock.Unlock()

	// check if delta is more than cap
	if delta > s.limit {
		return
	}

	// check for releasing too much
	if delta > s.count {
		panic("zync semaphore released more than acquired")
	}

	// decrement
	s.count -= delta

	// check what can be freed
	for {
		// get element from list
		next := s.waiters.Front()

		// if no more, stop
		if next == nil {
			break
		}

		// extract waiter
		w := next.Value.(waiter)

		// check if we have enough space for waiter,
		// we only release in a first-come first-served manner,
		// this could be done differently with a priority queue
		if s.limit-s.count < w.size {
			break
		}

		// acquire resources for waiter
		s.count += w.size
		s.waiters.Remove(next)
		close(w.done)
	}
}

// Len returns the current count. This refers to the current sum of unreleased
// Acquire calls delta.
func (s *Semaphore) Len() uint {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.count
}

// Cap returns the `limit` with which the Semaphore was created. this value
// does not change, expand, shrink, etc. You must create a new Semaphore if you
// want a different cap.
func (s *Semaphore) Cap() uint {
	return s.limit
}
