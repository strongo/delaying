package delaying

import (
	"context"
	"testing"
)

func TestInitNoopLogging(t *testing.T) {
	if registerDelayedFunc != nil {
		t.Fatal("registerDelayedFunc is NOT nil")
	}
	defer func() {
		registerDelayedFunc = nil
	}()

	InitNoopLogging()

	if registerDelayedFunc == nil {
		t.Fatal("registerDelayedFunc is nil")
	}

	d := registerDelayedFunc("test-id", nil)
	if d.ID() != "test-id" {
		t.Fatalf("expected test-id, got: %s", d.ID())
	}

	assertPanic := func(fn func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()
		fn()
	}

	assertPanic(func() { d.Implementation() })
	assertPanic(func() { _ = d.EnqueueWork(context.Background(), nil) })
	assertPanic(func() { _ = d.EnqueueWorkMulti(context.Background(), nil) })
}
