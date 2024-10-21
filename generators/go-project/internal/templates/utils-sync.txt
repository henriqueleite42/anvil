package utils

import (
	"errors"
	"sync"
	"time"
)

// Waits a WaitGroup with a timeout to be sure that it doesn't waits forever
// Returns an error if timeout occurs
func WaitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return nil
	case <-time.After(timeout):
		return errors.New("timeout")
	}
}
