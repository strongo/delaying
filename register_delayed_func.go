package delaying

type RegisterDelayedFunc func(key string, delayedFunc any) Delayer

var registerDelayedFunc RegisterDelayedFunc

func MustRegisterFunc(key string, delayedFunc any) Delayer {
	if registerDelayedFunc == nil {
		panic("No implementation has been registered. Application should call delaying.Init(...) before using delaying.MustRegisterFunc()")
	}
	return registerDelayedFunc(key, delayedFunc)
}

var _ RegisterDelayedFunc = MustRegisterFunc
