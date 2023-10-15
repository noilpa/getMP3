package shutdown

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Closer interface {
	Add(f func() error)
	Wait(ctx context.Context) error
	CloseAll()
}

type Close struct {
	sync.Mutex
	once    sync.Once
	funcs   []func() error
	osSig   SignalTrap
	timeout time.Duration
}

func New(timeout time.Duration) *Close {
	c := &Close{
		osSig:   TermSignalTrap(),
		timeout: timeout,
	}
	return c
}
func (c *Close) Add(f func() error) {
	c.Lock()
	defer c.Unlock()
	c.funcs = append(c.funcs, f)
}

func (c *Close) Wait(ctx context.Context) error {
	return c.osSig.Wait(ctx)
}

func (c *Close) CloseAll() {
	c.once.Do(func() {
		c.Lock()
		defer c.Unlock()

		wg := sync.WaitGroup{}
		wg.Add(len(c.funcs))
		doneCh := make(chan struct{})

		for i := range c.funcs {
			go func(f func() error) {
				if err := f(); err != nil {
					//log.Errorf("close func error: %v", err)
				}

				wg.Done()
			}(c.funcs[i])
		}

		go func() {
			wg.Wait()
			close(doneCh)
		}()

		select {
		case <-time.After(c.timeout):
			fmt.Println("shutdown timeout")
		case <-doneCh:
			fmt.Println("shutdown success")
		}
	})
}
