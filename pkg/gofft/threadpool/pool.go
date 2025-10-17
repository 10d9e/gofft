package threadpool

import (
	"runtime"
	"sync"
)

// Minimal pool that can run funcs concurrently.
type Pool struct {
	wg sync.WaitGroup
	ch chan func()
}

// New creates a pool with n workers; if n<=0 uses GOMAXPROCS.
func New(n int) *Pool {
	if n <= 0 { n = runtime.GOMAXPROCS(0) }
	p := &Pool{ ch: make(chan func(), 1024) }
	for i := 0; i < n; i++ {
		go func() {
			for fn := range p.ch {
				fn()
				p.wg.Done()
			}
		}()
	}
	return p
}

func (p *Pool) Submit(fn func()) {
	p.wg.Add(1)
	p.ch <- fn
}

func (p *Pool) Wait() { p.wg.Wait() }

func (p *Pool) Close() { close(p.ch) }
