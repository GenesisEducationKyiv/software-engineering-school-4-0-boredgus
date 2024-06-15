package internal

import (
	"context"
	"subscription-api/config"
	pb_ds "subscription-api/pkg/grpc/dispatch_service"
	"time"

	"google.golang.org/grpc"
)

type Scheduler interface {
	AddTask(name string, spec TaskSpec, task func())
	RemoveTask(taskId int)
	Run()
}

type DispatchService interface {
	SendDispatch(ctx context.Context, in *pb_ds.SendDispatchRequest, opts ...grpc.CallOption) (*pb_ds.SendDispatchResponse, error)
	GetAllDispatches(ctx context.Context, in *pb_ds.GetAllDispatchesRequest, opts ...grpc.CallOption) (*pb_ds.GetAllDispatchesResponse, error)
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
		d.log.Errorf("faild to parse time: %v", err)

		return
	}

	d.sc.AddTask(
		id,
		TaskSpec{Hours: t.Hour(), Mins: t.Minute()},
		func() {
			_, err = d.ds.SendDispatch(ctx, &pb_ds.SendDispatchRequest{DispatchId: id})
			if err != nil {
				d.log.Errorf("faild to send dispatch: %v", err)
			}
		})
}

func (d *DispatchDaemon) Run(ctx context.Context) {
	resp, err := d.ds.GetAllDispatches(ctx, &pb_ds.GetAllDispatchesRequest{})
	if err != nil {
		d.log.Errorf("faild to get dispatch: %v", err)

		return
	}

	for _, dispatch := range resp.Dispatches {
		d.scheduleDispatchSending(ctx, dispatch.Id, dispatch.SendAt)
	}
	d.sc.Run()
}
