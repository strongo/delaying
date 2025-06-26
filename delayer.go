package delaying

import "context"

type Delayer interface {
	ID() string
	Implementation() any
	EnqueueWork(c context.Context, params Params, args ...any) error
	EnqueueWorkMulti(c context.Context, params Params, args ...[]any) error
}

//type Function = Delayer // Deprecated: Use Delayer instead.

// NewDelayer creates a new Delayer.
func NewDelayer(
	id string,
	implementation any,
	enqueueWork func(c context.Context, params Params, args ...any) error,
	enqueueWorkMulti func(c context.Context, params Params, args ...[]any) error,
) Delayer {
	if implementation == nil {
		panic("implementation is nil")
	}
	if enqueueWork == nil {
		panic("enqueueWork is nil")
	}
	if enqueueWorkMulti == nil {
		panic("enqueueWorkMulti is nil")
	}
	return delayer{
		id:               id,
		implementation:   implementation,
		enqueueWork:      enqueueWork,
		enqueueWorkMulti: enqueueWorkMulti,
	}
}

// var NewFunction = NewDelayer // Deprecated: Use NewDelayer instead.

type delayer struct {
	id               string
	implementation   any
	enqueueWork      func(c context.Context, params Params, args ...any) error
	enqueueWorkMulti func(c context.Context, params Params, args ...[]any) error
}

func (f delayer) ID() string {
	return f.id
}

func (f delayer) Implementation() any {
	return f.implementation
}

func (f delayer) EnqueueWork(c context.Context, params Params, args ...any) error {
	return f.enqueueWork(c, params, args...)
}

func (f delayer) EnqueueWorkMulti(c context.Context, params Params, args ...[]any) error {
	return f.enqueueWorkMulti(c, params, args...)
}
