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

func NewRunner(jobs ...[]Job) *Runner {
	var all []Job

	for _, job := range jobs {
		all = append(all, job...)
	}

	return &Runner{jobs: all}
}

func (r *Runner) Run(ctx context.Context) {
	for _, job := range r.jobs {
		go job.Run(ctx)
	}
}
