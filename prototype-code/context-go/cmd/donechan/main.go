// Done channel — the manual fix.
//
// A plain `<-chan struct{}` is the simplest cancellation primitive in Go.
// It works, but it doesn't compose: every function invents its own channel,
// and combining timeouts or cascading cancellation gets messy fast.
// This is the pattern that `context.Context` standardizes.
package main

import (
	"fmt"
	"time"
)

func worker(done <-chan struct{}) {
	for i := 0; ; i++ {
		select {
		case <-done:
			fmt.Printf("worker: cancelled after %d ticks\n", i)
			return
		case <-time.After(100 * time.Millisecond):
			fmt.Printf("worker: tick %d\n", i)
		}
	}
}

func main() {
	done := make(chan struct{})
	go worker(done)

	time.Sleep(300 * time.Millisecond)
	close(done) // broadcast cancellation to every receiver

	time.Sleep(50 * time.Millisecond) // let the worker log its exit
}
