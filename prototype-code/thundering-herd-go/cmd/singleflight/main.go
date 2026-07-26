// Singleflight — coalesce N concurrent requests for the same key into 1.
//
// Cache-stampede defense: when a hot key expires and 100 goroutines all
// race to repopulate it, only one calls the upstream and the rest share
// the result. The herd becomes a single request.
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var calls int64

type call struct {
	wg  sync.WaitGroup
	val any
	err error
}

type group struct {
	mu sync.Mutex
	m  map[string]*call
}

func (g *group) Do(key string, fn func() (any, error)) (any, error) {
	g.mu.Lock()
	// create empty
	//map miss → "I'm the winner, do the work"
	if g.m == nil {
		g.m = map[string]*call{}
	}
	//map hit → "Someone's already doing it, wait"
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := &call{}
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	//Done() then delete
	c.wg.Done()

	//the winner must remove its entry from the map
	//
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
	return c.val, c.err
}

func upstream() (any, error) {
	atomic.AddInt64(&calls, 1)
	time.Sleep(50 * time.Millisecond) // simulate slow fetch
	return "user:42", nil
}

func main() {
	var g group
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = g.Do("user:42", upstream)
		}()
	}
	wg.Wait()
	fmt.Printf("upstream hit %d time(s) for 100 concurrent requests of the same key\n",
		atomic.LoadInt64(&calls))
}
