package internal

import (
	"fmt"
	"subscription-api/config"
	"time"

	"github.com/robfig/cron/v3"
)

type TaskSpec struct {
	Mins  int
	Hours int
}

func (s TaskSpec) Parse() string {
	return fmt.Sprintf("%v %v * * *", s.Mins, s.Hours)
}

type scheduler struct {
	cron           *cron.Cron
	scheduledTasks map[string]cron.EntryID
	logger         config.Logger
}

func NewScheduler(logger config.Logger) *scheduler {
	return &scheduler{
		logger:         logger,
		cron:           cron.New(cron.WithLocation(time.UTC)),
		scheduledTasks: make(map[string]cron.EntryID),
	}
}

func (s *scheduler) AddTask(name string, spec TaskSpec, task func()) {
	id, err := s.cron.AddFunc(spec.Parse(), task)
	if err != nil {
		s.logger.Errorf("failed to add cron func: %v", err)

		return
	}
	s.scheduledTasks[name] = id
}

func (s *scheduler) RemoveTask(taskId int) {
	s.cron.Remove(cron.EntryID(taskId))
}

func (s *scheduler) Run() {
	s.cron.Run()
}
