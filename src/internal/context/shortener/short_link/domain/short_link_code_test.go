package domain

import (
	"testing"
)

func TestNewCode(t *testing.T) {
	t.Run("Success generating code", func(t *testing.T) {
		t.Parallel()
		code, err := NewCode()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if string(code) == "" {
			t.Error("Expected code to be generated")
		}
		if len(string(code)) != 6 {
			t.Errorf("Expected code length of 6, got %d", len(string(code)))
		}
	})

	t.Run("Generates different codes", func(t *testing.T) {
		t.Parallel()
		code1, err1 := NewCode()
		code2, err2 := NewCode()
		if err1 != nil || err2 != nil {
			t.Fatal("Expected no errors generating codes")
		}
		if code1 == code2 {
			t.Error("Expected different codes to be generated")
		}
	})
}
