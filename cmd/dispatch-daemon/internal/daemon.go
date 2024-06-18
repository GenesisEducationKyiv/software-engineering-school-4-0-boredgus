package internal

import (
	"context"
	"subscription-api/config"
	"subscription-api/internal/services"
	"time"
)

type Scheduler interface {
	AddTask(name string, spec TaskSpec, task func())
	RemoveTask(taskId int)
	Run()
}

type DispatchService interface {
	SendDispatch(ctx context.Context, dispatchId string) error
	GetAllDispatches(ctx context.Context) ([]services.DispatchData, error)
}

type DispatchDaemon struct {
	ds  DispatchService
	log config.Logger
	sc  Scheduler
}

func NewDispatchDaemon(ds DispatchService, l config.Logger, sc Scheduler) *DispatchDaemon {
	return &DispatchDaemon{
		ds:  ds,
		log: l,
		sc:  sc,
	}
}

func (d *DispatchDaemon) scheduleDispatchSending(ctx context.Context, id, sendAt string) {
	t, err := time.Parse(time.TimeOnly, sendAt)
	if err != nil {
		d.log.Errorf("failed to parse time: %v", err)

		return
	}

	d.sc.AddTask(
		id,
		TaskSpec{Hours: t.Hour(), Mins: t.Minute()},
		func() {
			if err := d.ds.SendDispatch(ctx, id); err != nil {
				d.log.Errorf("failed to send dispatch: %v", err)
			}
		})
}

func (d *DispatchDaemon) Run(ctx context.Context) {
	dispatches, err := d.ds.GetAllDispatches(ctx)
	if err != nil {
		d.log.Errorf("failed to get dispatch: %v", err)

		return
	}

	for _, dispatch := range dispatches {
		d.scheduleDispatchSending(ctx, dispatch.Id, dispatch.SendAt)
	}
	d.sc.Run()
}
