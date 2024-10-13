package delaying

import "time"

type Params interface {
	Queue() string
	Path() string
	Delay() time.Duration
}

var _ Params = params{}

type params struct {
	queue string
	path  string
	delay time.Duration
}

func (p params) Queue() string {
	return p.queue
}

func (p params) Path() string {
	return p.path
}

func (p params) Delay() time.Duration {
	return p.delay
}
