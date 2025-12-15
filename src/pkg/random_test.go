package pkg

import (
	"errors"
	"io"
	"math/big"
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		code, err := GenerateShortCode(8)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(code) != 8 {
			t.Errorf("Expected code length 8, got %d", len(code))
		}
		for _, c := range code {
			if !containsChar(charset, c) {
				t.Errorf("Character %c not in charset", c)
			}
		}
	})

	t.Run("error from randFunc", func(t *testing.T) {
		t.Parallel()
		_, err := GenerateShortCodeWithRand(4, func(_ io.Reader, _ *big.Int) (*big.Int, error) {
			return nil, errors.New("forced error")
		})
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func containsChar(set string, c rune) bool {
	for _, s := range set {
		if s == c {
			return true
		}
	}
	return false
}
