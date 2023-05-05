package delaying

import "context"

type Logger interface {
	Debugf(c context.Context, format string, args ...interface{})
}

var log Logger

func debugf(c context.Context, format string, args ...interface{}) {
	if log != nil {
		log.Debugf(c, format, args...)
	}
}
