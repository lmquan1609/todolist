package asyncjob

import (
	"context"
	"time"
)

// Job requirements:
// 1. Job can do something (handler)
// 2. Job can retry
//   2.1. Config retry times and duration
// 3. Should be stateful
// 4. We should have job manager to manage jobs

type Job interface {
	Execute(ctx context.Context) error
	Retry(ctx context.Context) error
	State() JobState
	SetRetryDurations(times []time.Duration)
}

const (
	defaultMaxTimeout = time.Second * 10
)

var (
	defaultRetryTime = []time.Duration{time.Second, time.Second * 2, time.Second * 4}
)

type JobHandler func(context.Context) error

type JobState int

const (
	StateInit JobState = iota
	StateRunning
	StateFailed
	StateTimeout
	StateCompleted
	StateRetryFailed
)

func (js JobState) String() string {
	return [6]string{"Init", "Running", "Failed", "Timeout", "Completed", "RetryFailed"}[js]
}

type jobConfig struct {
	Name       string
	MaxTimeout time.Duration
	Retries    []time.Duration
}

type job struct {
	config     jobConfig
	handler    JobHandler
	state      JobState
	retryIndex int
	stopChan   chan bool
}

func NewJob(handler JobHandler, options ...OptionHdl) *job {
	j := job{
		config: jobConfig{
			MaxTimeout: defaultMaxTimeout,
			Retries:    defaultRetryTime,
		},
		handler:    handler,
		state:      StateInit,
		retryIndex: -1,
		stopChan:   make(chan bool),
	}

	for i := range options {
		options[i](&j.config)
	}
	return &j
}
