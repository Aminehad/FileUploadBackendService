package app

import (
	"context"
	"os"
	"os/signal"
)

type Service interface {
	Start(context.Context) error
	Stop()
}

type App struct {
	services []Service
}

func New(services ...Service) *App {
	return &App{services: services}
}

func (a *App) Start(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	ch := make(chan error)
	for _, s := range a.services {
		go func(s Service) {
			ch <- s.Start(ctx)
		}(s)
	}

	select {
	case <-ctx.Done():
		return nil
	case err := <-ch:
		return err
	}
}

func (a *App) Stop() {
	for _, s := range a.services {
		s.Stop()
	}
}
