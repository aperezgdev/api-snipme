package pkg

import (
	"strconv"
	"strings"
	"testing"
)

func TestMapNumberToString(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		expect []string
	}{
		{
			name:   "empty slice",
			input:  []int{},
			expect: []string{},
		},
		{
			name:   "single element",
			input:  []int{1},
			expect: []string{"1"},
		},
		{
			name:   "multiple elements",
			input:  []int{1, 2, 3},
			expect: []string{"1", "2", "3"},
		},
		{
			name:   "negative numbers",
			input:  []int{-1, -2, -3},
			expect: []string{"-1", "-2", "-3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map(tt.input, func(i int) string {
				return strconv.FormatInt(int64(i), 10)
			})
			if len(result) != len(tt.expect) {
				t.Errorf("expected length %d, got %d", len(tt.expect), len(result))
			}
			for i, v := range result {
				if v != tt.expect[i] {
					t.Errorf("expected %s at index %d, got %s", tt.expect[i], i, v)
				}
			}
		})
	}
}

func TestMapStringToUpper(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		expect []string
	}{
		{
			name:   "empty slice",
			input:  []string{},
			expect: []string{},
		},
		{
			name:   "single element",
			input:  []string{"hello"},
			expect: []string{"HELLO"},
		},
		{
			name:   "multiple elements",
			input:  []string{"hello", "world"},
			expect: []string{"HELLO", "WORLD"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map(tt.input, func(s string) string {
				return strings.ToUpper(s)
			})
			if len(result) != len(tt.expect) {
				t.Errorf("expected length %d, got %d", len(tt.expect), len(result))
			}
			for i, v := range result {
				if v != tt.expect[i] {
					t.Errorf("expected %s at index %d, got %s", tt.expect[i], i, v)
				}
			}
		})
	}
}

func TestMapToDouble(t *testing.T) {
	tests := []struct {
		name   string
		input  []float64
		expect []float64
	}{
		{
			name:   "empty slice",
			input:  []float64{},
			expect: []float64{},
		},
		{
			name:   "single element",
			input:  []float64{1.5},
			expect: []float64{3.0},
		},
		{
			name:   "multiple elements",
			input:  []float64{1.5, 2.5, 3.5},
			expect: []float64{3.0, 5.0, 7.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map(tt.input, func(f float64) float64 {
				return f * 2
			})
			if len(result) != len(tt.expect) {
				t.Errorf("expected length %d, got %d", len(tt.expect), len(result))
			}
			for i, v := range result {
				if v != tt.expect[i] {
					t.Errorf("expected %f at index %d, got %f", tt.expect[i], i, v)
				}
			}
		})
	}
}
