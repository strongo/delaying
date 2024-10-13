package delaying

import (
	"testing"
)

func TestInit(t *testing.T) {
	if registerDelayedFunc != nil {
		t.Fatal("registerDelayedFunc is NOT nil")
	}
	defer func() {
		registerDelayedFunc = nil
	}()
	var called bool
	f := func(key string, i any) Delayer {
		called = true
		return delayer{}
	}
	Init(f)

	if registerDelayedFunc == nil {
		t.Fatal("registerDelayedFunc is nil")
	}
	registerDelayedFunc("key", func() {})
	if !called {
		t.Fatal("called is false")
	}
	registerDelayedFunc = nil
}
