// Naive retry — the bug, not the fix.
//
// Every client retries immediately on failure, so the load spike lands
// exactly when the upstream is least healthy. This is the thundering herd.
package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var calls int64

func upstream() error {
	atomic.AddInt64(&calls, 1)
	return errors.New("upstream overloaded")
}

func client(wg *sync.WaitGroup) {
	defer wg.Done()
	for attempt := 0; attempt < 5; attempt++ {
		if err := upstream(); err == nil {
			return
		}
		// No backoff. No jitter. Just retry.
	}
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go client(&wg)
	}
	wg.Wait()
	fmt.Printf("upstream hit %d times by 100 clients x 5 retries (synchronized burst)\n",
		atomic.LoadInt64(&calls))
}
