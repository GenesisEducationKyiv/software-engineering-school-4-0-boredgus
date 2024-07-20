package scheduler

import (
	"context"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/app"
	messages "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/broker/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/entities"
	"github.com/robfig/cron/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	Publisher interface {
		PublishAsync(subject string, payload []byte) error
	}
	Converter interface {
		Convert(ctx context.Context, baseCcy string, targetCcies []string) (map[string]float64, error)
	}
	DispatchFetcher interface {
		GetAll(ctx context.Context) (map[string]entities.Dispatch, error)
	}

	dispatchScheduler struct {
		cron      *cron.Cron
		fetcher   DispatchFetcher
		broker    Publisher
		converter Converter
		logger    config.Logger
	}
)

const (
	Every15MinJobSpec string        = "*/15 * * * *"
	Duration15Min     time.Duration = 15 * time.Minute
	InvokeTimeout     time.Duration = 10 * time.Second
)

func NewDispatchScheduler(
	fetcher DispatchFetcher,
	broker Publisher,
	converter Converter,
	logger config.Logger,
) *dispatchScheduler {
	return &dispatchScheduler{
		fetcher:   fetcher,
		broker:    broker,
		converter: converter,
		logger:    logger,
		cron:      cron.New(cron.WithLocation(time.UTC)),
	}
}

func (s *dispatchScheduler) Run() {
	s.cron.AddJob(Every15MinJobSpec, NewJob(s.invokeSendingOfDispatches))

	s.cron.Run()
}

func (s *dispatchScheduler) Stop() {
	s.cron.Stop()
}

func (s *dispatchScheduler) invokeSendingOfDispatches() {
	ctx, cancel := context.WithTimeout(context.Background(), InvokeTimeout)
	defer cancel()

	dispatches, err := s.fetcher.GetAll(ctx)
	if err != nil {
		s.logger.Errorf("failed to fetch schedulled dispatches: %v", err)

		return
	}

	now := time.Now()

	for _, dispatch := range dispatches {
		if dispatch.SendAt.Before(now) && dispatch.SendAt.Add(Duration15Min).After(now) {
			s.processDispatch(ctx, dispatch)
		}
	}
}

func (s *dispatchScheduler) processDispatch(ctx context.Context, d entities.Dispatch) {
	msg := messages.SendDispatchCommand{
		EventType: messages.EventType_SEND_DISPATCH,
		Timestamp: timestamppb.New(time.Now().UTC()),
		Payload: &messages.Dispatch{
			Emails:  d.Emails,
			BaseCcy: d.BaseCcy,
		},
	}

	rates, err := s.converter.Convert(ctx, d.BaseCcy, d.TargetCcies)
	if err != nil {
		s.logger.Errorf("failed to get rates: %v", err)

		return
	}
	msg.Payload.Rates = rates

	marshalled, err := proto.Marshal(&msg)
	if err != nil {
		s.logger.Errorf("failed to marshal SendDispatchCommand: %v", err)

		return
	}

	if err := s.broker.PublishAsync(app.SendDispatchCommand, marshalled); err != nil {
		s.logger.Errorf("failed to emit SendDispatch command: %v", err)
	}
}
