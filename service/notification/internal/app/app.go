package app

import (
	"context"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/entities"
)

type (
	EventHandler interface {
		HandleEvents() error
		HandleCommands() error
	}
	DispatchFetcher interface {
		GetAll(ctx context.Context) (map[string]entities.Dispatch, error)
	}
	DispatchScheduler interface {
		Run()
		AddDispatches(ds map[string]entities.Dispatch)
		AddSubscription(entities.Subscription)
		// CancelSubscription(entities.Subscription)
		Stop()
	}
	app struct {
		eventHandler      EventHandler
		dispatchScheduler DispatchScheduler
		fetcher           DispatchFetcher
		logger            config.Logger
	}
)

func NewApp(
	eventHandler EventHandler,
	dispatchScheduler DispatchScheduler,
	logger config.Logger,
	fetcher DispatchFetcher,
) *app {
	return &app{
		eventHandler:      eventHandler,
		logger:            logger,
		dispatchScheduler: dispatchScheduler,
		fetcher:           fetcher,
	}
}

const MaxCountOfUploadAttempts int32 = 20
const IntervalBetweenUploadAttempts time.Duration = 20 * time.Minute

func (a *app) uploadOldDispatches() {
	var attemptNumber int32 = 0
	for attemptNumber < MaxCountOfUploadAttempts {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dispaches, err := a.fetcher.GetAll(ctx)
		if err != nil {
			attemptNumber++
			time.Sleep(IntervalBetweenUploadAttempts)
			return
		}

		a.dispatchScheduler.AddDispatches(dispaches)
		break
	}
}

func (a *app) Run() {
	go func() {
		if err := a.eventHandler.HandleEvents(); err != nil {
			a.logger.Error(err)
		}
	}()
	go func() {
		if err := a.eventHandler.HandleCommands(); err != nil {
			a.logger.Error(err)
		}
	}()
	go a.uploadOldDispatches()

	defer a.dispatchScheduler.Stop()

	a.dispatchScheduler.Run()
}
