// Exponential backoff — each retry waits 2^attempt longer.
//
// Spreads retries across time, but every client follows the same schedule:
// the herd just shifts in time. Necessary, not sufficient. See cmd/jitter.
package main

import (
	"errors"
	"fmt"
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
		time.Sleep(base << attempt) // 10ms, 20ms, 40ms, 80ms, 160ms
		// left bit-shift
		// base * 2^attempt -> single CPU instruction
		// attempt=0:  base << 0  = base * 1   =  10ms
		// attempt=4:  base << 4  = base * 16  = 160ms
		// idiomatic Go for powers-of-two scaling

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
	fmt.Printf("upstream hit %d times in %v (clients still synchronized at each tier)\n",
		atomic.LoadInt64(&calls), time.Since(start).Round(time.Millisecond))
}
