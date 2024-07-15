package scheduler

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/app"
	messages "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/broker/gen"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/entities"
	"github.com/robfig/cron/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type JobSpec struct {
	Mins  int
	Hours int
}

func NewJobSpec(t time.Time) JobSpec {
	return JobSpec{Mins: t.Minute(), Hours: t.Hour()}
}
func (s JobSpec) String() string {
	return fmt.Sprintf("%v %v * * *", s.Mins, s.Hours)
}

type (
	ScheduledDispatch struct {
		Data    *entities.Dispatch
		EntryID cron.EntryID
		Spec    JobSpec
	}

	Publisher interface {
		PublishAsync(subject string, payload []byte) error
	}
	Converter interface {
		Convert(ctx context.Context, baseCcy string, targetCcies []string) (map[string]float64, error)
	}

	dispatchScheduler struct {
		cron                *cron.Cron
		mu                  *sync.Mutex
		scheduledDispatches map[string]*ScheduledDispatch
		broker              Publisher
		converter           Converter
		logger              config.Logger
	}
)

func NewDispatchScheduler(
	broker Publisher,
	converter Converter,
	logger config.Logger,
) *dispatchScheduler {
	return &dispatchScheduler{
		mu:                  &sync.Mutex{},
		cron:                cron.New(cron.WithLocation(time.UTC)),
		scheduledDispatches: make(map[string]*ScheduledDispatch),
	}
}

func (s *dispatchScheduler) invokeSendingOfDispatch(d *entities.Dispatch) {
	msg := messages.SendDispatchCommand{
		EventType: messages.EventType_SEND_DISPATCH,
		Timestamp: timestamppb.New(time.Now().UTC()),
		Payload: &messages.Dispatch{
			Emails:  d.Emails,
			BaseCcy: d.BaseCcy,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), app.TimeoutOfProcessing)
	defer cancel()

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
		s.logger.Errorf("failed to emit SendDispatch commands: %v", err)
	}
}

func (s *dispatchScheduler) Run() {
	s.cron.Run()
}

func (s *dispatchScheduler) Stop() {
	s.cron.Stop()
}

func (s *dispatchScheduler) AddDispatches(ds map[string]entities.Dispatch) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, dispatch := range ds {
		scheduled := ScheduledDispatch{
			Data: &dispatch,
			Spec: NewJobSpec(dispatch.SendAt),
		}

		entryID, err := s.cron.AddJob(scheduled.Spec.String(), NewSendDispatchJob(scheduled.Data, s.invokeSendingOfDispatch))
		if err != nil {
			return
		}

		scheduled.EntryID = entryID

		s.scheduledDispatches[dispatch.ID] = &scheduled
	}
}

func (s *dispatchScheduler) AddSubscriberToDispatch(sub *entities.Subscription) {
	s.mu.Lock()
	defer s.mu.Unlock()

	dispatch, ok := s.scheduledDispatches[sub.DispatchID]
	if !ok {
		dispatch = &ScheduledDispatch{
			Data: sub.ToDispatch(),
			Spec: NewJobSpec(sub.SendAt),
		}

		s.scheduledDispatches[sub.DispatchID] = dispatch
	} else {
		dispatch.Data.Emails = append(dispatch.Data.Emails, sub.Email)
	}

	entryID, err := s.cron.AddJob(dispatch.Spec.String(), NewSendDispatchJob(dispatch.Data, s.invokeSendingOfDispatch))
	if err != nil {
		return
	}

	dispatch.EntryID = entryID
	s.scheduledDispatches[sub.DispatchID] = dispatch
}

func (s *dispatchScheduler) RemoveSubscriberFromDispatch(email, dispatchID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	schdledDispatch, ok := s.scheduledDispatches[dispatchID]
	if !ok {
		return
	}
	emailIndex := slices.Index(schdledDispatch.Data.Emails, email)
	if emailIndex < 0 {
		return
	}

	s.cron.Remove(schdledDispatch.EntryID)
	countOfSubscribers := len(schdledDispatch.Data.Emails)
	if countOfSubscribers == 1 {
		delete(s.scheduledDispatches, dispatchID)

		return
	}

	schdledDispatch.Data.Emails[emailIndex] = schdledDispatch.Data.Emails[countOfSubscribers-1]
	schdledDispatch.Data.Emails = schdledDispatch.Data.Emails[:countOfSubscribers-1]

	s.cron.AddJob(schdledDispatch.Spec.String(), NewSendDispatchJob(schdledDispatch.Data, s.invokeSendingOfDispatch))
}
