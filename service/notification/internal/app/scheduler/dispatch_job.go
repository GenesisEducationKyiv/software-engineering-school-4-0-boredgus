package scheduler

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-boredgus/service/notification/internal/entities"
)

type (
	Publisher interface {
		PublishAsync(subject string, payload []byte) error
	}

	sendDispatchJob struct {
		dispatch *entities.Dispatch
		invokerF func(*entities.Dispatch)
	}
)

func NewSendDispatchJob(dispatch *entities.Dispatch, f func(*entities.Dispatch)) *sendDispatchJob {
	return &sendDispatchJob{
		dispatch: dispatch,
	}
}

func (j *sendDispatchJob) Run() {
	j.invokerF(j.dispatch)
}
