package delaying

import "testing"

func TestMustRegisterFunc(t *testing.T) {
	orig := registerDelayedFunc
	defer func() { registerDelayedFunc = orig }()

	registerDelayedFunc = func(key string, i any) Delayer {
		return delayer{}
	}
	doSomething := func() {
	}
	MustRegisterFunc("key", doSomething)

	t.Run("nil_panics", func(t *testing.T) {
		registerDelayedFunc = nil
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic when registerDelayedFunc is nil")
			}
		}()
		MustRegisterFunc("key", doSomething)
	})
}
