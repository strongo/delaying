package delaying

func Init(f RegisterDelayedFunc) {
	if f == nil {
		panic("f is nil")
	}
	registerDelayedFunc = f
}
