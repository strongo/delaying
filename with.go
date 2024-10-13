package delaying

import (
	"strings"
	"time"
)

func With(queue, path string, delay time.Duration) Params {
	if strings.TrimSpace(queue) == "" {
		panic("queue is empty")
	}
	if strings.TrimSpace(path) == "" {
		panic("path is empty")
	}
	if delay < 0 {
		panic("delay is negative")
	}
	return params{queue: queue, path: path, delay: delay}
}
