package jsonata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	exp, err := Compile("$sum(example.value)")
	assert.NoError(t, err)
	assert.NotNil(t, exp)
}

func TestExpression_Evaluate(t *testing.T) {
	data := map[string]any{
		"example": []map[string]any{
			{"value": 4},
			{"value": 7},
			{"value": 13},
		},
	}

	exp := MustCompile("$sum(example.value)")

	output, err := exp.Evaluate(data, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(24), output)
}

func BenchmarkExpression_Compile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exp, err := Compile("$sum(example.value)")
		assert.NoError(b, err)
		assert.NotNil(b, exp)
	}
}

func BenchmarkExpression_Evaluate(b *testing.B) {
	data := map[string]any{
		"example": []map[string]any{
			{"value": 4},
			{"value": 7},
			{"value": 13},
		},
	}

	exp := MustCompile("$sum(example.value)")

	for i := 0; i < b.N; i++ {
		output, err := exp.Evaluate(data, nil)
		assert.NoError(b, err)
		assert.Equal(b, int64(24), output)
	}
}
