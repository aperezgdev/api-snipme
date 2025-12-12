package domain

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
)

type testLogger struct {
	infos  []string
	errors []string
	debugs []string
}

func (t *testLogger) Info(ctx context.Context, msg string, fields ...Field)  { t.infos = append(t.infos, msg) }
func (t *testLogger) Error(ctx context.Context, msg string, fields ...Field) { t.errors = append(t.errors, msg) }
func (t *testLogger) Debug(ctx context.Context, msg string, fields ...Field) { t.debugs = append(t.debugs, msg) }

func TestCompositeLogger(t *testing.T) {
	t.Run("delegates Info, Error, Debug to all loggers", func(t *testing.T) {
		l1 := &testLogger{}
		l2 := &testLogger{}
		cl := NewCompositeLogger(l1, l2)
		ctx := context.Background()
		cl.Info(ctx, "info")
		cl.Error(ctx, "err")
		cl.Debug(ctx, "dbg")
		assert.Equal(t, []string{"info"}, l1.infos)
		assert.Equal(t, []string{"info"}, l2.infos)
		assert.Equal(t, []string{"err"}, l1.errors)
		assert.Equal(t, []string{"err"}, l2.errors)
		assert.Equal(t, []string{"dbg"}, l1.debugs)
		assert.Equal(t, []string{"dbg"}, l2.debugs)
	})
}

func TestWithRequestID(t *testing.T) {
	t.Run("adds request_id if missing and present in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "request_id", "abc")
		fields := []Field{NewField("foo", "bar")}
		result := withRequestID(ctx, fields...)
		found := false
		for _, f := range result {
			if f.Key == "request_id" && f.Value == "abc" {
				found = true
			}
		}
		assert.True(t, found)
	})

	t.Run("does not add request_id if already present", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "request_id", "abc")
		fields := []Field{NewField("request_id", "abc")}
		result := withRequestID(ctx, fields...)
		count := 0
		for _, f := range result {
			if f.Key == "request_id" {
				count++
			}
		}
		assert.Equal(t, 1, count)
	})

	t.Run("does nothing if no request_id in context", func(t *testing.T) {
		ctx := context.Background()
		fields := []Field{NewField("foo", "bar")}
		result := withRequestID(ctx, fields...)
		assert.Equal(t, fields, result)
	})
}
