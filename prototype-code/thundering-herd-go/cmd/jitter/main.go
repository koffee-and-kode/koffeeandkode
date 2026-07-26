// Full jitter — wait a random duration in [0, base*2^attempt).
//
// Decorrelates the retry schedules of independent clients so the upstream
// sees a smeared load curve instead of a comb of synchronized spikes.
package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"
)

var calls int64

func upstream() error {
	atomic.AddInt64(&calls, 1)
	return errors.New("upstream overloaded")
}

func client(wg *sync.WaitGroup) {
	defer wg.Done()
	base := 10 * time.Millisecond
	for attempt := 0; attempt < 5; attempt++ {
		if err := upstream(); err == nil {
			return
		}
		window := base << attempt
		time.Sleep(time.Duration(rand.Int64N(int64(window))))
		// window is the ceiling, not the sleep itself
		// rand(0, base << attempt)
	}
}

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go client(&wg)
	}
	wg.Wait()
	fmt.Printf("upstream hit %d times in %v (retries de-correlated across clients)\n",
		atomic.LoadInt64(&calls), time.Since(start).Round(time.Millisecond))
}
