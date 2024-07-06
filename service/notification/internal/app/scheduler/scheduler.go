package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/entities"
	"github.com/robfig/cron/v3"
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

	dispatchScheduler struct {
		cron                *cron.Cron
		mu                  *sync.Mutex
		scheduledDispatches map[string]*ScheduledDispatch
		dispatchInvokerF    func(*entities.Dispatch)
	}
)

func NewDispatchScheduler() *dispatchScheduler {
	return &dispatchScheduler{
		mu:                  &sync.Mutex{},
		cron:                cron.New(cron.WithLocation(time.UTC)),
		scheduledDispatches: make(map[string]*ScheduledDispatch),
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

		entryID, err := s.cron.AddJob(scheduled.Spec.String(), NewSendDispatchJob(scheduled.Data, s.dispatchInvokerF))
		if err != nil {
			return
		}

		scheduled.EntryID = entryID

		s.scheduledDispatches[dispatch.ID] = &scheduled
	}
}

func (s *dispatchScheduler) AddSubscription(sub entities.Subscription, invokerF func(*entities.Dispatch)) {
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

	entryID, err := s.cron.AddJob(dispatch.Spec.String(), NewSendDispatchJob(dispatch.Data, invokerF))
	if err != nil {
		return
	}

	dispatch.EntryID = entryID
	s.scheduledDispatches[sub.DispatchID] = dispatch
}

// func (s *dispatchScheduler) CancelSubscription(sub entities.Subscription) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	schdledDispatch, ok := s.scheduledDispatches[sub.DispatchID]
// 	if !ok {
// 		return
// 	}
// 	emailIndex := slices.Index(schdledDispatch.Data.Emails, sub.Email)
// 	if emailIndex < 0 {
// 		return
// 	}
// 	countOfSubscribers := len(schdledDispatch.Data.Emails)
// 	schdledDispatch.Data.Emails[emailIndex] = schdledDispatch.Data.Emails[countOfSubscribers-1]
// 	schdledDispatch.Data.Emails = schdledDispatch.Data.Emails[:countOfSubscribers-1]
// 	s.cron.Remove(schdledDispatch.entryID)
// 	s.cron.AddJob(schdledDispatch.spec.String(), NewDispatchJob(schdledDispatch.Data, s.broker))
// }
