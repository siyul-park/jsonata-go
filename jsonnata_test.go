package jsonata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	exp, err := Compile("$sum(example.value)")
	assert.NoError(t, err)
	assert.NotNil(t, exp)

	assert.NoError(t, exp.Close())
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
	defer func() { _ = exp.Close() }()

	output, err := exp.Evaluate(data, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(24), output)
}

func BenchmarkExpression_Compile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		exp, err := Compile("$sum(example.value)")
		assert.NoError(b, err)
		assert.NotNil(b, exp)

		assert.NoError(b, exp.Close())
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
	defer func() { _ = exp.Close() }()

	for i := 0; i < b.N; i++ {
		output, err := exp.Evaluate(data, nil)
		assert.NoError(b, err)
		assert.Equal(b, int64(24), output)
	}
}
