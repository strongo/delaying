package delaying

import (
	"context"
	"reflect"
	"time"
)

// sleepFunc is a function type for sleeping
type sleepFunc func(delay time.Duration)

// sleep is a package-level variable that holds the function for sleeping
// It can be replaced during tests to avoid actual sleeping
var sleep sleepFunc = func(delay time.Duration) {
	if delay == 0 {
		delay = 100 * time.Millisecond
	}

	// Convert time based on magnitude
	switch {
	case delay > 24*time.Hour: // If the delay is > 1 day, convert days to seconds
		delay = time.Duration(float64(delay) / float64(24*time.Hour) * float64(time.Second))
	case delay > time.Hour: // If the delay is > 1 hour, convert hours to seconds
		delay = time.Duration(float64(delay) / float64(time.Hour) * float64(time.Second))
	case delay > time.Minute: // If the delay is > 1 minute, convert minutes to seconds
		delay = time.Duration(float64(delay) / float64(time.Minute) * float64(time.Second))
	case delay > 10*time.Second: // If the delay is > 10 seconds, divide by 10
		delay = time.Duration(float64(delay) / 10)
	case delay > 2*time.Second: // If the delay is > 2 seconds, divide by 2
		delay = time.Duration(float64(delay) / 2)
	default: // For smaller delays, maintain original behavior
		delay = time.Duration(float64(delay) / float64(time.Minute) * float64(time.Second))
	}

	// Ensure minimum delay is 1 second
	if delay < time.Second {
		delay = time.Second
	}
	timeSleep(delay)
}

var timeSleep = time.Sleep

// executeWorkerWithDelay handles the common logic for executing a worker function:
// 1. Applies delay if specified in params
// 2. Converts arguments to reflect.Value
// 3. Calls the worker function using reflection
func executeWorkerWithDelay(worker any, params Params, args []any) {
	// Check if params is not nil and delay is specified
	if params != nil {
		if delay := params.Delay(); delay > 0 {
			sleep(delay)
		}
	}

	// Execute the worker with reflection
	workerValue := reflect.ValueOf(worker)
	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		reflectArgs[i] = reflect.ValueOf(arg)
	}
	workerValue.Call(reflectArgs)
}

func GoRoutineWithLog(key string, worker any) Delayer {
	workers := make(map[string]any)
	workers[key] = worker
	return NewDelayer(key, worker,
		func(c context.Context, params Params, args ...any) error {
			debugf(c, "%s.EnqueueWork(%+v): %+v", key, args, params)
			// Call the worker with args in a goroutine
			go func() {
				executeWorkerWithDelay(worker, params, args)
			}()
			return nil
		},
		func(c context.Context, params Params, arguments ...[]any) error {
			debugf(c, "%s.EnqueueWorkMulti(%+v): %+v", key, arguments, params)
			for _, args := range arguments {
				args2 := args
				go func() {
					executeWorkerWithDelay(worker, params, args2)
				}()
			}
			return nil
		},
	)
}
