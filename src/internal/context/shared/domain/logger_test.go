package domain

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestField_String(t *testing.T) {
	t.Run("String returns key: value", func(t *testing.T) {
		f := Field{Key: "foo", Value: "bar"}
		assert.Equal(t, "foo: bar", f.String())
	})
}

func TestNewField(t *testing.T) {
	t.Run("returns Field with key and value", func(t *testing.T) {
		f := NewField("k", "v")
		assert.Equal(t, "k", f.Key)
		assert.Equal(t, "v", f.Value)
	})
}

func TestDummyLogger(t *testing.T) {
	logger := DummyLogger{}
	ctx := context.Background()
	t.Run("Info does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() { logger.Info(ctx, "msg") })
	})
	t.Run("Error does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() { logger.Error(ctx, "msg") })
	})
	t.Run("Debug does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() { logger.Debug(ctx, "msg") })
	})
}
