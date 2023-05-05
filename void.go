package delaying

import (
	"context"
)

func VoidWithLog(key string, i any) Function {
	doNothing := func() {}
	return NewFunction(key, doNothing,
		func(c context.Context, params Params, args ...interface{}) error {
			debugf(c, "%s.EnqueueWork(%+v): %+v", key, args, params)
			return nil
		},
		func(c context.Context, params Params, args ...[]interface{}) error {
			debugf(c, "%s.EnqueueWorkMulti(%+v): %+v", key, args, params)
			return nil
		},
	)
}
