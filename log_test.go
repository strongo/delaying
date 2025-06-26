package delaying

import (
	"context"
	"testing"
)

// mockLogger is a mock implementation of the logger interface for testing
type mockLogger struct {
	called   bool
	lastCtx  context.Context
	lastFmt  string
	lastArgs []any
}

// Debugf implements the logger interface
func (m *mockLogger) Debugf(c context.Context, format string, args ...any) {
	m.called = true
	m.lastCtx = c
	m.lastFmt = format
	m.lastArgs = args
}

func TestInitLogger(t *testing.T) {
	// Store the original logger to restore after tests
	originalLogger := log

	// Restore the original logger after tests
	t.Cleanup(func() {
		log = originalLogger
	})

	t.Run("sets logger and uses it", func(t *testing.T) {
		// Create a mock logger
		mockLog := &mockLogger{}

		// Initialize the logger with the mock
		InitLogger(mockLog)

		// Verify the logger was set
		if log != mockLog {
			t.Error("Expected log to be set to mockLog")
		}

		// Create a test context
		ctx := context.Background()

		// Call debugf with some test data
		testFormat := "test format %s %d"
		testArgs := []any{"string", 123}
		debugf(ctx, testFormat, testArgs...)

		// Verify the mock logger was called
		if !mockLog.called {
			t.Error("Expected mockLog.Debugf to be called")
		}

		// Verify the context was passed correctly
		if mockLog.lastCtx != ctx {
			t.Errorf("Expected context to be %v, got %v", ctx, mockLog.lastCtx)
		}

		// Verify the format string was passed correctly
		if mockLog.lastFmt != testFormat {
			t.Errorf("Expected format to be %q, got %q", testFormat, mockLog.lastFmt)
		}

		// Verify the arguments were passed correctly
		if len(mockLog.lastArgs) != len(testArgs) {
			t.Errorf("Expected %d args, got %d", len(testArgs), len(mockLog.lastArgs))
		} else {
			for i, arg := range testArgs {
				if mockLog.lastArgs[i] != arg {
					t.Errorf("Expected arg %d to be %v, got %v", i, arg, mockLog.lastArgs[i])
				}
			}
		}
	})

	t.Run("sets nil logger", func(t *testing.T) {
		// Initialize with a non-nil logger first
		mockLog := &mockLogger{}
		InitLogger(mockLog)

		// Then set to nil
		InitLogger(nil)

		// Verify the logger was set to nil
		if log != nil {
			t.Error("Expected log to be nil")
		}

		// Call debugf - should not panic
		ctx := context.Background()
		debugf(ctx, "test format")

		// Verify the previous mock logger was not called
		if mockLog.called {
			t.Error("Expected mockLog.Debugf not to be called")
		}
	})

	t.Run("replaces existing logger", func(t *testing.T) {
		// Initialize with first mock logger
		firstMock := &mockLogger{}
		InitLogger(firstMock)

		// Initialize with second mock logger
		secondMock := &mockLogger{}
		InitLogger(secondMock)

		// Verify the logger was set to the second mock
		if log != secondMock {
			t.Error("Expected log to be set to secondMock")
		}

		// Call debugf
		ctx := context.Background()
		debugf(ctx, "test format")

		// Verify the second mock was called
		if !secondMock.called {
			t.Error("Expected secondMock.Debugf to be called")
		}

		// Verify the first mock was not called
		if firstMock.called {
			t.Error("Expected firstMock.Debugf not to be called")
		}
	})
}
