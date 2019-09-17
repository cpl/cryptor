// Copied and adapted from:
// https://github.com/golang/sync/blob/master/semaphore/semaphore_test.go
// under the following LICENSE:
// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// https://github.com/golang/sync/blob/master/LICENSE

package zync_test

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"cpl.li/go/cryptor/pkg/zync"
)

const maxSleep = 1 * time.Millisecond

func hammerWeighted(sem *zync.Semaphore, n uint, loops int) {
	for i := 0; i < loops; i++ {
		sem.Acquire(n)
		time.Sleep(time.Duration(
			rand.Int63n(int64(maxSleep/time.Nanosecond))) * time.Nanosecond)
		sem.Release(n)
	}
}

func TestWeighted(t *testing.T) {
	t.Parallel()
	rand.Seed(1)

	cpus := runtime.NumCPU()
	loops := 10000 / cpus

	sem := zync.NewSemaphore(uint(cpus))
	var wg sync.WaitGroup
	wg.Add(cpus)

	for i := uint(0); i < uint(cpus); i++ {
		go func() {
			defer wg.Done()
			hammerWeighted(sem, i, loops)
		}()
	}
	wg.Wait()
}

func TestReleasePanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("released un-acquired weighted semaphore without panic")
		}
	}()
	sem := zync.NewSemaphore(1)
	sem.Release(1)
}

func TestReleaseLarge(t *testing.T) {
	t.Parallel()

	sem := zync.NewSemaphore(2)
	sem.Acquire(2)
	sem.Release(3)

	assert.Equal(t, sem.Len(), uint(2), "unexpected semaphore len")
}

func TestAcquirePanic(t *testing.T) {
	t.Parallel()

	sem := zync.NewSemaphore(10)
	sem.Acquire(5)
	sem.Acquire(15)
	if ok := sem.TryAcquire(15); ok {
		t.Fatalf("semaphore succesful TryAcquire with delta too large")
	}

	assert.Equal(t, sem.Len(), uint(5), "unexpected semaphore len")
}

func TestWeightedTryAcquire(t *testing.T) {
	t.Parallel()

	sem := zync.NewSemaphore(2)
	tries := [4]bool{}

	sem.Acquire(1)
	tries[0] = sem.TryAcquire(1)
	tries[1] = sem.TryAcquire(1)
	sem.Release(2)
	tries[2] = sem.TryAcquire(1)
	sem.Acquire(1)
	tries[3] = sem.TryAcquire(1)

	want := []bool{true, false, true, false}
	for i := range tries {
		if tries[i] != want[i] {
			t.Errorf("tries[%d]: got %t, want %t", i, tries[i], want[i])
		}
	}
}

func tryAcquire(sem *zync.Semaphore, delta uint) bool {
	done := make(chan interface{})
	go func() {
		sem.Acquire(delta)
		close(done)
	}()
	select {
	case <-done:
		return true
	case <-time.After(10 * time.Millisecond):
		return false
	}
}

func TestWeightedAcquire(t *testing.T) {
	t.Parallel()

	sem := zync.NewSemaphore(2)
	tries := [4]bool{}
	sem.Acquire(1)
	tries[0] = tryAcquire(sem, 1)
	tries[1] = tryAcquire(sem, 1)

	sem.Release(2)

	tries[2] = tryAcquire(sem, 1)
	tries[3] = tryAcquire(sem, 1)

	want := []bool{true, false, true, false}
	for i := range tries {
		if tries[i] != want[i] {
			t.Errorf("tries[%d]: got %t, want %t", i, tries[i], want[i])
		}
	}
}

func TestSemaphoreCap(t *testing.T) {
	t.Parallel()

	sem := zync.NewSemaphore(4)
	assert.Equal(t, sem.Cap(), uint(4), "unexpected semaphore cap")

	sem = zync.NewSemaphore(1)
	assert.Equal(t, sem.Cap(), uint(1), "unexpected semaphore cap")

	sem = zync.NewSemaphore(1000)
	assert.Equal(t, sem.Cap(), uint(1000), "unexpected semaphore cap")
}
