// Leak — the bug, not the fix.
//
// A goroutine that has no way to be told to stop. main() returns, but the
// goroutine is still running until the process exits. In a long-lived
// server this is how you accumulate background workers that won't die.
package main

import (
	"fmt"
	"time"
)

func worker() {
	for i := 0; ; i++ {
		// Pretend this is real work: a poll, a tick, a tail.
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("worker: tick %d\n", i)
	}
}

func main() {
	go worker()
	// "Caller" only needs the worker for 300ms of real work.
	time.Sleep(300 * time.Millisecond)
	fmt.Println("caller: done, but worker has no idea")
	// The goroutine is still running. In a server, this leaks forever.
	time.Sleep(200 * time.Millisecond)
}
