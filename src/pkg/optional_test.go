
package pkg

import "testing"

func TestOptional(t *testing.T) {
	t.Run("Some and IsPresent", func(t *testing.T) {
		t.Parallel()
		opt := Some(42)
		if !opt.IsPresent() {
			t.Errorf("Expected IsPresent to be true")
		}
		if v := opt.Get(); v != 42 {
			t.Errorf("Expected value 42, got %v", v)
		}
	})

	t.Run("EmptyOptional", func(t *testing.T) {
		t.Parallel()
		opt := EmptyOptional[int]()
		if opt.IsPresent() {
			t.Errorf("Expected IsPresent to be false")
		}
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic on Get for empty optional")
			}
		}()
		_ = opt.Get()
	})

	t.Run("OrElse with present value", func(t *testing.T) {
		t.Parallel()
		opt := Some("hello")
		if v := opt.OrElse("default"); v != "hello" {
			t.Errorf("Expected 'hello', got %v", v)
		}
	})

	t.Run("OrElse with empty optional", func(t *testing.T) {
		t.Parallel()
		empty := EmptyOptional[string]()
		if v := empty.OrElse("default"); v != "default" {
			t.Errorf("Expected 'default', got %v", v)
		}
	})
}
