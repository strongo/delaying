package delaying

import "context"

func InitNoopLogging() {
	Init(func(key string, i any) Delayer {
		return noOpFunction{id: key}
	})
}

type noOpFunction struct {
	id string
}

func (n noOpFunction) ID() string {
	return n.id
}

func (n noOpFunction) Implementation() any {
	panic("implement me") //TODO implement me
}

func (n noOpFunction) EnqueueWork(_ context.Context, _ Params, _ ...interface{}) error {
	panic("implement me") //TODO implement me
}

func (n noOpFunction) EnqueueWorkMulti(_ context.Context, _ Params, _ ...[]interface{}) error {
	panic("implement me") //TODO implement me
}
