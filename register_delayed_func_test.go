package delaying

import "testing"

func TestMustRegisterFunc(t *testing.T) {
	registerDelayedFunc = func(key string, i any) Delayer {
		return delayer{}
	}
	doSomething := func() {
	}
	MustRegisterFunc("key", doSomething)
}
