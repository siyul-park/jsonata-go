package util

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNumeric(t *testing.T) {
	testCases := []struct {
		when   any
		expect bool
	}{
		{
			when:   nil,
			expect: false,
		},
		{
			when:   "",
			expect: false,
		},
		{
			when:   false,
			expect: false,
		},
		{
			when:   []any{},
			expect: false,
		},
		{
			when:   map[string]any{},
			expect: false,
		},
		{
			when:   0,
			expect: true,
		},
		{
			when:   0.0,
			expect: true,
		},
		{
			when:   math.NaN(),
			expect: false,
		},
		{
			when:   math.Inf(0),
			expect: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.when), func(t *testing.T) {
			ok := IsNumeric(tc.when)
			assert.Equal(t, tc.expect, ok)
		})
	}
}

func TestIsArrayOfStrings(t *testing.T) {
	testCases := []struct {
		when   any
		expect bool
	}{
		{
			when:   nil,
			expect: true,
		},
		{
			when:   "",
			expect: false,
		},
		{
			when:   []string{"a", "b", "c"},
			expect: true,
		},
		{
			when:   []any{"a", "b", "c"},
			expect: true,
		},
		{
			when:   []any{"a", 0, "c"},
			expect: false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.when), func(t *testing.T) {
			ok := IsArrayOfStrings(tc.when)
			assert.Equal(t, tc.expect, ok)
		})
	}
}
