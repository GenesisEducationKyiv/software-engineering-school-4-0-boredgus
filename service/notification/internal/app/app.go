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
		ctx, cancel := context.WithTimeout(context.Background(), TimeoutOfProcessing)
		defer cancel()

		dispaches, err := a.fetcher.GetAll(ctx)
		if err != nil {
			a.logger.Errorf("failed to fetch scheduled dispatches: %v", err)

			attemptNumber++
			time.Sleep(IntervalBetweenUploadAttempts)

			return
		}

		a.dispatchScheduler.AddDispatches(dispaches)
		a.logger.Infof("successfully scheduled %v dispatches", len(dispaches))

		break
	}
}

func (a *app) Run() {
	if err := a.eventHandler.HandleEvents(); err != nil {
		a.logger.Error(err)
	}
	if err := a.eventHandler.HandleCommands(); err != nil {
		a.logger.Error(err)
	}
	go a.uploadOldDispatches()

	defer a.dispatchScheduler.Stop()
	a.dispatchScheduler.Run()
}
