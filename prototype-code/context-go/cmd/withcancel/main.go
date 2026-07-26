// WithCancel — the same fix, in the shape stdlib uses everywhere.
//
// `context.WithCancel(parent)` gives you the same close-a-channel-to-cancel
// behavior as `donechan`, but with a standard interface (`context.Context`)
// that the rest of the ecosystem — net/http, database/sql, gRPC — already
// speaks. The `cancel` function must always be called to release resources.
package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context) error {
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			// ctx.Err() distinguishes Canceled vs DeadlineExceeded.
			return fmt.Errorf("worker stopping after %d ticks: %w", i, ctx.Err())
		case <-time.After(100 * time.Millisecond):
			fmt.Printf("worker: tick %d\n", i)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // safe even after the manual cancel() below

	go func() {
		if err := worker(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	time.Sleep(300 * time.Millisecond)
	cancel() // closes ctx.Done() — every receiver unblocks

	time.Sleep(50 * time.Millisecond)
}
