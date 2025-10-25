package domain

import (
	"context"
	"log/slog"
	"os"

	"github.com/aperezgdev/api-snipme/src/pkg"
)

type Field struct {
	Key   string
	Value any
}

func (f Field) String() string {
	return f.Key + ": " + f.Value.(string)
}

func NewField(key string, value any) Field {
	return Field{Key: key, Value: value}
}

type Logger interface {
	Info(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
	Debug(ctx context.Context, msg string, fields ...Field)
}

type DummyLogger struct{}

func (DummyLogger) Info(ctx context.Context, msg string, fields ...Field)  {}
func (DummyLogger) Error(ctx context.Context, msg string, fields ...Field) {}
func (DummyLogger) Debug(ctx context.Context, msg string, fields ...Field) {}

type ConsoleLogger struct {
	logger slog.Logger
}

func NewConsoleLogger() ConsoleLogger {
	return ConsoleLogger{
		logger: *slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
}

func (c ConsoleLogger) Info(ctx context.Context, msg string, fields ...Field) {
	c.logger.Info(msg, fieldsToLogAttr(withRequestID(ctx, fields...)...)...)
}

func (c ConsoleLogger) Error(ctx context.Context, msg string, fields ...Field) {
	c.logger.Error(msg, fieldsToLogAttr(withRequestID(ctx, fields...)...)...)
}

func (c ConsoleLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	c.logger.Debug(msg, fieldsToLogAttr(withRequestID(ctx, fields...)...)...)
}

func withRequestID(ctx context.Context, fields ...Field) []Field {
	const requestIDKey = "request_id"
	reqID, _ := ctx.Value(requestIDKey).(string)
	if reqID == "" {
		return fields
	}
	for _, f := range fields {
		if f.Key == "request_id" {
			return fields
		}
	}
	return append(fields, NewField("request_id", reqID))
}

type CompositeLogger struct {
	loggers []Logger
}

func NewCompositeLogger(loggers ...Logger) *CompositeLogger {
	return &CompositeLogger{loggers: loggers}
}

func (c *CompositeLogger) Info(ctx context.Context, msg string, fields ...Field) {
	for _, logger := range c.loggers {
		logger.Info(ctx, msg, fields...)
	}
}

func (c *CompositeLogger) Error(ctx context.Context, msg string, fields ...Field) {
	for _, logger := range c.loggers {
		logger.Error(ctx, msg, fields...)
	}
}

func (c *CompositeLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	for _, logger := range c.loggers {
		logger.Debug(ctx, msg, fields...)
	}
}

func fieldsToLogAttr(fields ...Field) []any {
	return pkg.Map(fields, func(t Field) any {
		return slog.Any(t.Key, t.Value)
	})
}
