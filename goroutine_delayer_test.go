package delaying

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestSleep(t *testing.T) {
	// Store the original timeSleep function
	originalTimeSleep := timeSleep

	// Restore the original timeSleep function after the test
	t.Cleanup(func() {
		timeSleep = originalTimeSleep
	})

	// Define test cases
	tests := []struct {
		name          string
		inputDelay    time.Duration
		expectedDelay time.Duration
	}{
		{
			name:          "with normal delay",
			inputDelay:    5 * time.Second,
			expectedDelay: 2*time.Second + 500*time.Millisecond, // 5 seconds divided by 2
		},
		{
			name:          "with zero delay",
			inputDelay:    0,
			expectedDelay: time.Second, // For zero delay, use default of 100 ms, but the minimum is 1 second
		},
		{
			name:          "with delay > 2 seconds (exact)",
			inputDelay:    4 * time.Second,
			expectedDelay: 2 * time.Second, // 4 seconds divided by 2
		},
		{
			name:          "with delay > 2 seconds (odd)",
			inputDelay:    3 * time.Second,
			expectedDelay: 1*time.Second + 500*time.Millisecond, // 3 seconds divided by 2
		},
		{
			name:          "with delay > 10 seconds",
			inputDelay:    20 * time.Second,
			expectedDelay: 2 * time.Second, // 20 seconds divided by 10
		},
		{
			name:          "with delay > 10 seconds (exact)",
			inputDelay:    15 * time.Second,
			expectedDelay: 1*time.Second + 500*time.Millisecond, // 15 seconds divided by 10
		},
		{
			name:          "with small delay",
			inputDelay:    30 * time.Second,
			expectedDelay: 3 * time.Second, // 30 seconds divided by 10
		},
		{
			name:          "with minute delay",
			inputDelay:    5 * time.Minute,
			expectedDelay: 5 * time.Second, // 5 minutes converted to 5 seconds
		},
		{
			name:          "with hour delay",
			inputDelay:    3 * time.Hour,
			expectedDelay: 3 * time.Second, // 3 hours converted to 3 seconds
		},
		{
			name:          "with day delay",
			inputDelay:    2 * 24 * time.Hour,
			expectedDelay: 2 * time.Second, // 2 days converted to 2 seconds
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock implementation of time.Sleep
			var mockCalled bool
			var mockDelay time.Duration
			mockTimeSleep := func(delay time.Duration) {
				mockCalled = true
				mockDelay = delay
			}

			// Replace timeSleep with the mock to test the sleep function
			timeSleep = mockTimeSleep

			// Call the sleep function with the test delay
			sleep(tt.inputDelay)

			// Verify that the mock was called
			if !mockCalled {
				t.Error("Expected sleep to call timeSleep, but it was not called")
			}

			// Allow a small margin of error for floating point calculations
			lowerBound := time.Duration(float64(tt.expectedDelay) * 0.99)
			upperBound := time.Duration(float64(tt.expectedDelay) * 1.01)
			if mockDelay < lowerBound || mockDelay > upperBound {
				t.Errorf("Expected sleep to call timeSleep with delay %v, but got %v", tt.expectedDelay, mockDelay)
			}
		})
	}
}

func TestGoRoutineWithLog(t *testing.T) {
	t.Run("calls worker with correct arguments", func(t *testing.T) {
		// Create a channel to signal when the worker has been called
		done := make(chan bool)
		var receivedArgs []any

		// Create a worker function that captures the arguments it was called with
		worker := func(arg1 string, arg2 int, arg3 bool) {
			receivedArgs = []any{arg1, arg2, arg3}
			done <- true
		}

		// Create a delayer using GoRoutineWithLog
		delayer := GoRoutineWithLog("test-key", worker)

		// Call EnqueueWork with some arguments
		expectedArgs := []any{"test", 42, true}
		err := delayer.EnqueueWork(context.Background(), nil, expectedArgs...)
		if err != nil {
			t.Fatalf("EnqueueWork returned an error: %v", err)
		}

		// Wait for the worker to be called or timeout
		select {
		case <-done:
			// Check that the worker was called with the correct arguments
			if len(receivedArgs) != len(expectedArgs) {
				t.Fatalf("Expected %d arguments, got %d", len(expectedArgs), len(receivedArgs))
			}
			for i, expected := range expectedArgs {
				if receivedArgs[i] != expected {
					t.Errorf("Expected argument %d to be %v, got %v", i, expected, receivedArgs[i])
				}
			}
		case <-time.After(time.Second):
			t.Fatal("Timed out waiting for worker to be called")
		}
	})

	t.Run("works with different worker function signatures", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(2)

		// Test with a worker that takes no arguments
		noArgsWorker := func() {
			wg.Done()
		}
		noArgsDelayer := GoRoutineWithLog("no-args", noArgsWorker)
		err := noArgsDelayer.EnqueueWork(context.Background(), nil)
		if err != nil {
			t.Fatalf("EnqueueWork returned an error: %v", err)
		}

		// Test with a worker that takes a single argument
		singleArgWorker := func(arg string) {
			if arg != "hello" {
				t.Errorf("Expected argument to be 'hello', got '%s'", arg)
			}
			wg.Done()
		}
		singleArgDelayer := GoRoutineWithLog("single-arg", singleArgWorker)
		err = singleArgDelayer.EnqueueWork(context.Background(), nil, "hello")
		if err != nil {
			t.Fatalf("EnqueueWork returned an error: %v", err)
		}

		// Wait for all workers to be called or timeout
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// All workers were called successfully
		case <-time.After(time.Second):
			t.Fatal("Timed out waiting for workers to be called")
		}
	})

	t.Run("handles EnqueueWorkMulti", func(t *testing.T) {
		// Create a worker that takes a string argument
		worker := func(arg string) {}
		delayer := GoRoutineWithLog("multi-test", worker)

		// EnqueueWorkMulti should not panic
		// Pass a slice of slices, where each inner slice contains arguments for one call
		args := [][]any{
			{"test"},
		}
		err := delayer.EnqueueWorkMulti(context.Background(), nil, args...)
		if err != nil {
			t.Fatalf("EnqueueWorkMulti returned an error: %v", err)
		}

		// Wait a bit for the goroutine to complete
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("EnqueueWorkMulti calls worker correct number of times with single argument", func(t *testing.T) {
		// Create a channel to signal when the worker has been called
		done := make(chan string)

		// Create a worker function that takes a string argument
		worker := func(arg string) {
			done <- arg
		}

		delayer := GoRoutineWithLog("multi-count-test-single", worker)

		// Create arguments for EnqueueWorkMulti
		// Each element in the outer slice is a set of arguments for one call
		args := [][]any{
			{"arg1"},
			{"arg2"},
			{"arg3"},
		}

		// Call EnqueueWorkMulti with the arguments
		err := delayer.EnqueueWorkMulti(context.Background(), nil, args...)
		if err != nil {
			t.Fatalf("EnqueueWorkMulti returned an error: %v", err)
		}

		// Collect results from the worker calls
		receivedArgs := make([]string, 0, len(args))
		timeout := time.After(time.Second)

		// Wait for all worker calls or timeout
		for i := 0; i < len(args); i++ {
			select {
			case arg := <-done:
				receivedArgs = append(receivedArgs, arg)
			case <-timeout:
				t.Fatalf("Timed out waiting for worker to be called %d times, got %d calls", len(args), i)
			}
		}

		// Check that the worker was called the correct number of times
		if len(receivedArgs) != len(args) {
			t.Errorf("Expected worker to be called %d times, but was called %d times", len(args), len(receivedArgs))
		}

		// Verify that all arguments were received
		expectedArgs := []string{"arg1", "arg2", "arg3"}
		// Check that all expected arguments were received (order may vary due to goroutines)
		for _, expected := range expectedArgs {
			found := false
			for _, received := range receivedArgs {
				if received == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected argument %q was not received", expected)
			}
		}
	})

	t.Run("EnqueueWorkMulti calls worker correct number of times with three arguments", func(t *testing.T) {
		// Create a channel to signal when the worker has been called
		type workerArgs struct {
			s string
			i int
			b bool
		}
		done := make(chan workerArgs)

		// Create a worker function that takes three arguments (string, int, bool)
		worker := func(s string, i int, b bool) {
			done <- workerArgs{s, i, b}
		}

		delayer := GoRoutineWithLog("multi-count-test-three", worker)

		// Create arguments for EnqueueWorkMulti
		// Each element in the outer slice is a set of arguments for one call
		args := [][]any{
			{"arg1", 1, true},
			{"arg2", 2, false},
			{"arg3", 3, true},
		}

		// Call EnqueueWorkMulti with the arguments
		err := delayer.EnqueueWorkMulti(context.Background(), nil, args...)
		if err != nil {
			t.Fatalf("EnqueueWorkMulti returned an error: %v", err)
		}

		// Collect results from the worker calls
		receivedArgs := make([]workerArgs, 0, len(args))
		timeout := time.After(time.Second)

		// Wait for all worker calls or timeout
		for i := 0; i < len(args); i++ {
			select {
			case arg := <-done:
				receivedArgs = append(receivedArgs, arg)
			case <-timeout:
				t.Fatalf("Timed out waiting for worker to be called %d times, got %d calls", len(args), i)
			}
		}

		// Check that the worker was called the correct number of times
		if len(receivedArgs) != len(args) {
			t.Errorf("Expected worker to be called %d times, but was called %d times", len(args), len(receivedArgs))
		}

		// Verify that all arguments were received
		expectedArgs := []workerArgs{
			{"arg1", 1, true},
			{"arg2", 2, false},
			{"arg3", 3, true},
		}

		// Check that all expected arguments were received (order may vary due to goroutines)
		for _, expected := range expectedArgs {
			found := false
			for _, received := range receivedArgs {
				if received.s == expected.s && received.i == expected.i && received.b == expected.b {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected argument {%q, %d, %t} was not received", expected.s, expected.i, expected.b)
			}
		}
	})

	t.Run("returns correct ID and Implementation", func(t *testing.T) {
		worker := func() {}
		key := "test-key"
		delayer := GoRoutineWithLog(key, worker)

		if delayer.ID() != key {
			t.Errorf("Expected ID to be '%s', got '%s'", key, delayer.ID())
		}

		// Can't directly compare function values in Go
		// Instead, check that Implementation() returns a non-nil value
		if delayer.Implementation() == nil {
			t.Errorf("Expected Implementation to be non-nil")
		}
	})
}
