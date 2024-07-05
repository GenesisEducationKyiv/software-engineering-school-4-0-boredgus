package app

// import (
// 	"context"
// 	"time"

// 	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/daemon/scheduler/internal/config"
// )

// type (
// 	Scheduler interface {
// 		AddTask(name string, spec TaskSpec, task func())
// 		RemoveTask(taskId int)
// 		Run()
// 	}
// 	CurrencyService interface {
// 		Convert(ctx context.Context, baseCcy string, targetCcies []string) (map[string]float64, error)
// 	}
// 	DispatchService interface {
// 		GetAllDispatches(ctx context.Context) ([]dispatch_client.Dispatch, error)
// 	}
// 	Broker interface {
// 	}
// )

// type daemon struct {
// 	currencyService CurrencyService
// 	logger          config.Logger
// 	scheduler       Scheduler
// 	broker          Broker
// }

// func NewDaemon(
// 	currencyService CurrencyService,
// 	logger config.Logger,
// 	scheduler Scheduler,
// 	broker Broker,
// ) *daemon {
// 	return &daemon{
// 		currencyService: currencyService,
// 		logger:          logger,
// 		scheduler:       scheduler,
// 		broker:          broker,
// 	}
// }

// func (d *daemon) scheduleDispatchSending(ctx context.Context, id, sendAt string) {
// 	t, err := time.Parse(time.TimeOnly, sendAt)
// 	if err != nil {
// 		d.logger.Errorf("failed to parse time: %v", err)

// 		return
// 	}

// 	d.scheduler.AddTask(
// 		id,
// 		TaskSpec{Hours: t.Hour(), Mins: t.Minute()},
// 		func() {
// 			if err := d.currencyService.SendDispatch(ctx, id); err != nil {
// 				d.logger.Errorf("failed to send dispatch: %v", err)
// 			}
// 		})
// }

// func (d *daemon) Run(ctx context.Context) {
// 	defer d.shutdown()
// 	dispatches, err := d.currencyService.GetAllDispatches(ctx)
// 	if err != nil {
// 		d.logger.Errorf("failed to get dispatch: %v", err)

// 		return
// 	}

// 	for _, dispatch := range dispatches {
// 		d.scheduleDispatchSending(ctx, dispatch.Id, dispatch.SendAt)
// 	}
// 	d.scheduler.Run()
// }

// func (d *daemon) shutdown() {

// }
