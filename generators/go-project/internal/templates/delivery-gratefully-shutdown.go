package templates

const GratefullyShutdownTempl = `package delivery

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type GracefullyShutdownInput struct {
	Logger     *zerolog.Logger
	Timeout    time.Duration
	Deliveries []Delivery
}

// Gracefully shutdown all the active deliveries
func GracefullyShutdown(i *GracefullyShutdownInput) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan

	i.Logger.Info().
		Str("sig", sig.String()).
		Msg("gracefully shutdown started")

	var wg sync.WaitGroup

	for _, v := range i.Deliveries {
		wg.Add(1)
		go func() {
			defer wg.Done()

			i.Logger.Debug().
				Dur("timeout", i.Timeout).
				Msgf("start %s.Cancel", v.Name())

			v.Cancel(i.Timeout)
		}()
	}

	wg.Wait()
}
`
