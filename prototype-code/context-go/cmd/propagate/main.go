// Propagate — why ctx is worth the bookkeeping.
//
// One `cancel()` at the top stops every goroutine three layers deep,
// because each layer accepted `ctx` and passed it along. This is the
// payoff that makes `context.Context` more than just a fancy channel.
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func leaf(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("leaf %d: stopping (%v) after %d ticks\n", id, ctx.Err(), i)
			return
		case <-time.After(100 * time.Millisecond):
			fmt.Printf("leaf %d: tick %d\n", id, i)
		}
	}
}

func middle(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	var inner sync.WaitGroup
	for i := 0; i < 3; i++ {
		inner.Add(1)
		go leaf(ctx, i, &inner)
	}
	inner.Wait()
	fmt.Println("middle: all leaves done")
}

func top(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(1)
	go middle(ctx, &wg)
	wg.Wait()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go top(ctx)

	time.Sleep(300 * time.Millisecond)
	cancel() // one call, four goroutines stop

	time.Sleep(100 * time.Millisecond)
}
