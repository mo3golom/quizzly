package cmd

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"quizzly/pkg/logger"
	"sync"
	"syscall"
	"time"
)

type Service interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
}

type ServiceRunner struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	log logger.Logger
}

func NewRunner(log logger.Logger) *ServiceRunner {
	ctx, cancel := context.WithCancel(context.Background())
	return &ServiceRunner{
		ctx:    ctx,
		cancel: cancel,
		log:    log,
	}
}

func (r *ServiceRunner) Start(services ...Service) {
	if len(services) == 0 {
		return
	}

	for _, service := range services {
		r.wg.Add(1)
		go func(s Service) {
			defer r.wg.Done()
			s.Start(r.ctx)
			// Call Stop when context is done
			s.Stop(r.ctx)
		}(service)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	r.cancel()

	done := make(chan struct{})
	go func() {
		r.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		r.log.Info("All services stopped successfully")
	case <-time.After(10 * time.Second):
		r.log.Error("Shutdown timed out", errors.New("shutdown timed out"))
	}
}
