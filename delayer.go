package delaying

import "context"

type Delayer interface {
	ID() string
	Implementation() any
	EnqueueWork(c context.Context, params Params, args ...interface{}) error
	EnqueueWorkMulti(c context.Context, params Params, args ...[]interface{}) error
}

// Deprecated: Use Delayer instead.
type Function = Delayer

// NewDelayer creates a new Delayer.
func NewDelayer(
	id string,
	implementation any,
	enqueueWork func(c context.Context, params Params, args ...interface{}) error,
	enqueueWorkMulti func(c context.Context, params Params, args ...[]interface{}) error,
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

// Deprecated: Use NewDelayer instead.
var NewFunction = NewDelayer

type delayer struct {
	id               string
	implementation   any
	enqueueWork      func(c context.Context, params Params, args ...interface{}) error
	enqueueWorkMulti func(c context.Context, params Params, args ...[]interface{}) error
}

func (f delayer) ID() string {
	return f.id
}

func (f delayer) Implementation() any {
	return f.implementation
}

func (f delayer) EnqueueWork(c context.Context, params Params, args ...interface{}) error {
	return f.enqueueWork(c, params, args...)
}

func (f delayer) EnqueueWorkMulti(c context.Context, params Params, args ...[]interface{}) error {
	return f.enqueueWorkMulti(c, params, args...)
}
