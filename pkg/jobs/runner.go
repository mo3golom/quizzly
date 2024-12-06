package jobs

import (
	"context"
	"errors"
	"fmt"
	"quizzly/pkg/logger"
	"runtime/debug"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	fallbackInterval = time.Second * 5
)

type (
	DefaultRunner struct {
		log logger.Logger

		jobs      map[string]Job
		enrichers []ContextEnricher
		isActive  bool
	}
)

func NewDefaultRunner(log logger.Logger, enrichers ...ContextEnricher) *DefaultRunner {
	return &DefaultRunner{jobs: map[string]Job{}, log: log, enrichers: enrichers}
}

func (r *DefaultRunner) RegisterAll(jobs ...Job) error {
	for _, job := range jobs {
		if err := r.Register(job); err != nil {
			return err
		}
	}

	return nil
}

func (r *DefaultRunner) Register(job Job) error {
	if _, alreadyExists := r.jobs[job.Name()]; alreadyExists {
		return fmt.Errorf("job '%s' already registered", job.Name())
	}

	r.jobs[job.Name()] = job

	return nil
}

func (r *DefaultRunner) Start(ctx context.Context) {
	group, ctx := errgroup.WithContext(ctx)

	for _, job := range r.jobs {
		job := job

		group.Go(func() error {
			return r.run(ctx, job)
		})
	}
	r.isActive = true

	fmt.Println("Jobs are running")
	if err := group.Wait(); err != nil {
		r.log.Error("jobs error", err)
		return
	}
}

func (r *DefaultRunner) Stop(_ context.Context) {
	r.log.Info("Jobs were stopped")
}

func (r *DefaultRunner) run(ctx context.Context, job Job) error {
	for {
		iterationCtx := ctx
		for _, enricher := range r.enrichers {
			iterationCtx = enricher.Enrich(iterationCtx)
		}

		r.log.Info(
			"run job",
			logger.Field{
				Key:   "name",
				Value: job.Name(),
			},
		)
		interval, err := r.execute(iterationCtx, job)
		if errors.Is(err, context.Canceled) {
			r.log.Info("job canceled")
			return nil
		}

		if err = r.await(iterationCtx, interval); err != nil {
			if errors.Is(err, context.Canceled) {
				r.log.Info("job await canceled")
				return nil
			}
			r.log.Error("await failed", err)

			return err
		}
	}
}

func (r *DefaultRunner) execute(ctx context.Context, job Job) (interval time.Duration, err error) {
	defer func() {
		recovered := recover()
		if recovered != nil {
			r.log.Error(
				"panic happened in job",
				errors.New("panic happened in job"),
				logger.Field{
					Key:   "panic",
					Value: recovered,
				},
				logger.Field{
					Key:   "stack_trace",
					Value: string(debug.Stack()),
				},
			)

			interval = fallbackInterval
		}
	}()

	if err = job.Perform(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			r.log.Info("job canceled")
			return fallbackInterval, err
		}

		r.log.Error("job returned error", err)
	}

	determined, err := job.DetermineInterval(ctx)
	if err != nil {
		r.log.Error("job failed to determine interval", err)

		return fallbackInterval, nil
	}
	if determined == nil {
		r.log.Error("job returned nil interval", errors.New("job returned nil interval"))

		return fallbackInterval, nil
	}

	return *determined, nil
}

func (r *DefaultRunner) await(ctx context.Context, interval time.Duration) error {
	var ticker *time.Ticker
	err := safe(func() error {
		// тут внутри паника может случиться
		ticker = time.NewTicker(interval)
		return nil
	})
	if err != nil {
		r.log.Error(
			"failed to create ticker for next interval, using fallback",
			err,
			logger.Field{
				Key:   "ticker_interval",
				Value: interval,
			},
		)

		ticker = time.NewTicker(fallbackInterval)
	}
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ticker.C:
		return nil
	}
}

func safe(block func() error) (err error) {
	defer func() {
		recovered := recover()
		if recovered != nil {
			err = fmt.Errorf("%s", recovered)
		}
	}()

	return block()
}
