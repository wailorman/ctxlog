package ctxlog

import (
	"context"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

type logKey string

// LoggerContextKey _
const LoggerContextKey logKey = "ctxlog.logger"

// LoggerFieldsContextKey _
const LoggerFieldsContextKey logKey = "ctxlog.fields"

// FromContext _
func FromContext(ctx context.Context, prefix string) logrus.FieldLogger {
	logger, ok := ctx.Value(LoggerContextKey).(logrus.FieldLogger)

	if !ok {
		return nil
	}

	return logger.WithField("prefix", prefix)
}

func FromContextOr(ctx context.Context, prefix string, fallbackLoggers ...logrus.FieldLogger) logrus.FieldLogger {
	logger := FromContext(ctx, prefix)

	if logger == nil {
		for _, fallbackLogger := range fallbackLoggers {
			if fallbackLogger != nil {
				logger = fallbackLogger
				break
			}
		}
	}

	logger = logger.WithField("prefix", prefix)

	return logger
}

func ContextWithLogger(ctx context.Context, logger *logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, LoggerContextKey, logger)
}

func WithPrefix(ctx context.Context, prefix string) context.Context {
	return WithField(ctx, "prefix", prefix)
}

func WithField(ctx context.Context, key string, value interface{}) context.Context {
	return WithFields(ctx, logrus.Fields{key: value})
}

func WithFields(ctx context.Context, lFields logrus.Fields) context.Context {
	fields, ok := ctx.Value(LoggerFieldsContextKey).(logrus.Fields)

	if !ok {
		return ctx
	}

	for key, value := range lFields {
		fields[key] = value
	}

	return context.WithValue(ctx, LoggerFieldsContextKey, fields)
}

func EmptyLogger() logrus.FieldLogger {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	return logger
}
