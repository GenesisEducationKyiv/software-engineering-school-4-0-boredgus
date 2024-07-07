package app

import (
	"context"
	"fmt"
	"time"

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
		AddDispatches(ds map[string]entities.Dispatch)
		Stop()
	}
	app struct {
		handler           MessageHandler
		dispatchScheduler DispatchScheduler
		fetcher           DispatchFetcher
		logger            config.Logger
	}
)

func NewApp(
	handler MessageHandler,
	dispatchScheduler DispatchScheduler,
	logger config.Logger,
	fetcher DispatchFetcher,
) *app {
	return &app{
		handler:           handler,
		logger:            logger,
		dispatchScheduler: dispatchScheduler,
		fetcher:           fetcher,
	}
}

const MaxCountOfUploadAttempts int32 = 3
const IntervalBetweenUploadAttempts time.Duration = 1 * time.Minute

func (a *app) uploadOldDispatches() error {
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutOfProcessing)
	defer cancel()

	dispatches, err := a.fetcher.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch scheduled dispatches: %v", err)
	}

	a.dispatchScheduler.AddDispatches(dispatches)
	a.logger.Infof("successfully scheduled %v dispatches", len(dispatches))

	return nil
}

func (a *app) Run() {
	if err := a.uploadOldDispatches(); err != nil {
		a.logger.Error(err)

		return
	}

	if err := a.handler.HandleMessages(); err != nil {
		a.logger.Error(err)
	}

	defer a.dispatchScheduler.Stop()
	a.dispatchScheduler.Run()
}
