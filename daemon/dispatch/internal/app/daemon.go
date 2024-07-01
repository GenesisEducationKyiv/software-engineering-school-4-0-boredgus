package app

import (
	"context"
	"time"

	dispatch_client "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/daemon/dispatch/internal/clients/dispatch"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/daemon/dispatch/internal/config"
)

type Scheduler interface {
	AddTask(name string, spec TaskSpec, task func())
	RemoveTask(taskId int)
	Run()
}

type DispatchService interface {
	SendDispatch(ctx context.Context, dispatchId string) error
	GetAllDispatches(ctx context.Context) ([]dispatch_client.Dispatch, error)
}

type DispatchDaemon struct {
	service   DispatchService
	logger    config.Logger
	scheduler Scheduler
}

func NewDispatchDaemon(ds DispatchService, l config.Logger, sc Scheduler) *DispatchDaemon {
	return &DispatchDaemon{
		service:   ds,
		logger:    l,
		scheduler: sc,
	}
}

func (d *DispatchDaemon) scheduleDispatchSending(ctx context.Context, id, sendAt string) {
	t, err := time.Parse(time.TimeOnly, sendAt)
	if err != nil {
		d.logger.Errorf("failed to parse time: %v", err)

		return
	}

	d.scheduler.AddTask(
		id,
		TaskSpec{Hours: t.Hour(), Mins: t.Minute()},
		func() {
			if err := d.service.SendDispatch(ctx, id); err != nil {
				d.logger.Errorf("failed to send dispatch: %v", err)
			}
		})
}

func (d *DispatchDaemon) Run(ctx context.Context) {
	dispatches, err := d.service.GetAllDispatches(ctx)
	if err != nil {
		d.logger.Errorf("failed to get dispatch: %v", err)

		return
	}

	for _, dispatch := range dispatches {
		d.scheduleDispatchSending(ctx, dispatch.Id, dispatch.SendAt)
	}
	d.scheduler.Run()
}
