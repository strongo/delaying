package delaying

import (
	"context"
)

func VoidWithLog(key string, _ any) Delayer {
	doNothing := func() {}
	return NewDelayer(key, doNothing,
		func(c context.Context, params Params, args ...any) error {
			debugf(c, "%s.EnqueueWork(%+v): %+v", key, args, params)
			return nil
		},
		func(c context.Context, params Params, args ...[]any) error {
			debugf(c, "%s.EnqueueWorkMulti(%+v): %+v", key, args, params)
			return nil
		},
	)
}
