package delaying

import "context"

type logger interface {
	Debugf(c context.Context, format string, args ...any)
}

var log logger

func debugf(c context.Context, format string, args ...any) {
	if log != nil {
		log.Debugf(c, format, args...)
	}
}

func InitLogger(logger interface {
	Debugf(c context.Context, format string, args ...any)
}) {
	log = logger
}
