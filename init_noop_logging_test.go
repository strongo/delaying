package delaying

import "testing"

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
}
