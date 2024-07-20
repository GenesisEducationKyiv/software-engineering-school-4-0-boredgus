package scheduler

type job struct {
	invokerF func()
}

func NewJob(f func()) *job {
	return &job{invokerF: f}
}

func (j *job) Run() {
	j.invokerF()
}
