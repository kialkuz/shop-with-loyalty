package scheduler

import (
	"context"
)

type Job interface {
	Run(ctx context.Context)
}

type Runner struct {
	jobs []Job
}

func NewRunner(jobs []Job) *Runner {
	return &Runner{jobs: jobs}
}

func (r *Runner) Run(ctx context.Context) {
	for _, job := range r.jobs {
		go job.Run(ctx)
	}
}
