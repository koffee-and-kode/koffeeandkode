---
title: "Thundering Herd 0 to 1"
date: 2026-05-07
tags: [reliability, go, retry, backoff, distributed-systems]
rwp_categories: [coordination]
excerpt: "What the thundering herd actually looks like, and four bare-bones Go programs that walk from the bug to the fix — naive retry, exponential backoff, jitter, singleflight."
read_time: 7
repo_url: https://github.com/koffee-and-kode/koffeeandkode/tree/main/prototype-code/thundering-herd-go
---

A *thundering herd* is what happens when many clients react to the same event at the same time — a cache key expiring, an upstream blipping, a service restarting — and slam a downstream system in unison. The system was already under pressure. Now everyone retried at once and made it worse.

This post walks the problem from zero: a tiny Go program that exhibits the bug, then three small fixes that each address one part of why it happens. The runnable code lives at [`prototype-code/thundering-herd-go/`](https://github.com/koffee-and-kode/koffeeandkode/tree/main/prototype-code/thundering-herd-go), one `cmd/<technique>/main.go` per idea.

Each program is intentionally bare-bones — 100 goroutines, a fake upstream, an atomic counter. Just enough wiring to see the shape of the problem.

## The bug — naive retry

```go
func client(wg *sync.WaitGroup) {
    defer wg.Done()
    for attempt := 0; attempt < 5; attempt++ {
        if err := upstream(); err == nil {
            return
        }
        // No backoff. No jitter. Just retry.
    }
}
```

100 clients, 5 retries each. Run it:

```
upstream hit 500 times by 100 clients x 5 retries (synchronized burst)
```

500 calls, all synchronized within microseconds of each other. The upstream sees a wall, not a curve. If it was failing because it was overloaded, it now stays overloaded — the retries *are* the load.

This is the failure mode in its purest form. Two things make it bad:

1. **No spacing** — the retries arrive immediately, while the upstream is still recovering.
2. **Synchronization** — every client follows the same schedule, so the retries arrive *together*.

The next three fixes peel those apart.

## Fix 1 — exponential backoff

The first instinct is right: stop retrying instantly. Wait, and wait longer each time.

```go
base := 10 * time.Millisecond
for attempt := 0; attempt < 5; attempt++ {
    if err := upstream(); err == nil {
        return
    }
    time.Sleep(base << attempt) // 10ms, 20ms, 40ms, 80ms, 160ms
}
```

The doubling matters. A linear delay (10ms, 20ms, 30ms…) keeps adding load at a near-constant rate. Doubling means each tier carries half the pressure of the previous one, so a persistently failing upstream gets exponentially less traffic the longer it stays down — which is exactly what you want a struggling system to see.

But run the program and the count is still 500:

```
upstream hit 500 times in 316ms (clients still synchronized at each tier)
```

Backoff spread the calls across time, but every client waits *the same* 10ms, then the same 20ms, then 40ms… The herd is still a herd — it's just marching in step. Necessary, not sufficient.

## Fix 2 — jitter

The fix for synchronization is randomness. Replace each fixed wait with a random one in the same range:

```go
window := base << attempt
time.Sleep(time.Duration(rand.Int64N(int64(window))))
```

This is *full jitter* (the AWS Architecture Blog name). Each client picks a uniformly random delay inside the backoff window, so two clients almost never retry at the same instant.

```
upstream hit 500 times in 267ms (retries de-correlated across clients)
```

Same total calls, but smeared across the window instead of stacked at its edges. The upstream now sees a curve — load it can absorb instead of a spike that knocks it over.

There are other variants (*equal jitter*, *decorrelated jitter*) that trade off how aggressive vs. how spread-out retries get. Full jitter is the safe default and the one to reach for first.

Backoff + jitter is the answer when independent clients are retrying transient failures. But it doesn't help with the other shape of thundering herd: many clients all wanting the *same* thing at the *same* time, on the first try.

## Fix 3 — singleflight (request coalescing)

Picture a hot cache key — the homepage feed, a popular user profile — and the moment its TTL expires. 100 requests in flight all miss the cache simultaneously. Each one fires off to the upstream to repopulate. The upstream gets 100 identical queries; 99 of them are wasted work.

Backoff doesn't help here — nothing has *failed* yet. The fix is to notice that 100 in-flight requests for the same key only need *one* upstream call. The rest can wait and share the result.

```go
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
    if g.m == nil {
        g.m = map[string]*call{}
    }
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
    c.wg.Done()

    g.mu.Lock()
    delete(g.m, key)
    g.mu.Unlock()
    return c.val, c.err
}
```

The first caller for a key creates a `call` entry, runs `fn`, and signals its `WaitGroup` when done. Every later caller arriving while that's in flight finds the existing entry, waits on the same group, and reads the same result. After the call completes, the entry is deleted so the *next* miss does its own fetch.

Run it:

```
upstream hit 1 time(s) for 100 concurrent requests of the same key
```

100 requests, 1 upstream call. The herd collapses into a single fetch.

This pattern is what [`golang.org/x/sync/singleflight`](https://pkg.go.dev/golang.org/x/sync/singleflight) implements — it's worth using directly in real code (it handles panics, `Forget`, and the shared/unshared distinction) — but the 25-line version above shows why it works.

## Picking the right one

Each technique addresses a different failure shape:

- **Exponential backoff** — clients are retrying a transient error and need to back off *fast* when it persists.
- **Jitter** — multiple clients are independently retrying and would otherwise sync up.
- **Singleflight** — multiple clients want the same answer at the same time and you want exactly one upstream call.

Backoff and jitter compose — you almost always want both together for any retry loop. Singleflight is orthogonal: it sits in front of an expensive read, not around a retryable write.

A few things deliberately left out for the MVP — all of which are covered in the follow-up post (publishing soon):

- **Circuit breakers** — close the call entirely when an upstream is clearly down, instead of retrying with backoff forever.
- **Token-bucket rate limiting** — cap the *rate* at which retries leave a client, not just space them out.
- **Probabilistic early expiration** — refresh hot cache keys *before* they expire so no synchronized miss ever happens.
- **Server-side load shedding** — `Retry-After` and 429s give the upstream a way to push the herd back even if clients aren't well-behaved.

The four programs in this post are the floor: if you don't have these, the others won't save you.

## Running the examples

```bash
cd prototype-code/thundering-herd-go
go run ./cmd/naive
go run ./cmd/backoff
go run ./cmd/jitter
go run ./cmd/singleflight
```

Every program is under 70 lines. The numbers you'll see are deterministic in shape — same call counts, similar timings — and the difference between *naive* and *jitter* is the difference between a wall and a curve.
