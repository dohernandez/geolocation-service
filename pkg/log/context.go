package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ctxKey struct{}

// FromContext returns logger from context
func FromContext(ctx context.Context) logrus.FieldLogger {
	log, ok := ctx.Value(ctxKey{}).(logrus.FieldLogger)
	if !ok {
		return nil
	}

	return log
}

// ToContext returns context with logger
func ToContext(ctx context.Context, log logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, ctxKey{}, log)
}
