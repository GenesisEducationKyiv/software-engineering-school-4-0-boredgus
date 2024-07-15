package stubs

import "github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/entities"

type schedulerMock struct {
}

func NewSchedulerMock() *schedulerMock {
	return &schedulerMock{}
}

func (m *schedulerMock) AddSubscriberToDispatch(*entities.Subscription)        {}
func (m *schedulerMock) RemoveSubscriberFromDispatch(email, dispatchID string) {}
