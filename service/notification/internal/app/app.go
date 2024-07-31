package app

import (
	"context"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/entities"
)

type (
	MessageHandler interface {
		HandleMessages() error
	}
	DispatchFetcher interface {
		GetAll(ctx context.Context) (map[string]entities.Dispatch, error)
	}
	DispatchScheduler interface {
		Run()
		Stop()
	}
	app struct {
		handler   MessageHandler
		scheduler DispatchScheduler
		logger    config.Logger
	}
)

func NewApp(
	handler MessageHandler,
	dispatchScheduler DispatchScheduler,
	logger config.Logger,
) *app {
	return &app{
		handler:   handler,
		logger:    logger,
		scheduler: dispatchScheduler,
	}
}

func (a *app) Run() {
	if err := a.handler.HandleMessages(); err != nil {
		a.logger.Error(err)
	}

	a.logger.Infof("notification service started...")

	defer a.scheduler.Stop()
	a.scheduler.Run()
}
