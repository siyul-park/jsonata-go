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

	output, err := exp.Evaluate(data)
	assert.NoError(t, err)
	assert.Equal(t, int64(24), output)
}

func TestExpression_Assign(t *testing.T) {
	exp := MustCompile("$greet")

	err := exp.Assign("greet", "Hello world")
	assert.NoError(t, err)

	output, err := exp.Evaluate(nil)
	assert.NoError(t, err)
	assert.Equal(t, "Hello world", output)
}

func TestExpression_RegisterFunction(t *testing.T) {
	exp := MustCompile("$greet()")

	err := exp.RegisterFunction("greet", func(f *Focus, args ...any) (any, error) { return "Hello world", nil })
	assert.NoError(t, err)

	output, err := exp.Evaluate(nil)
	assert.NoError(t, err)
	assert.Equal(t, "Hello world", output)
}

func TestExpression_Ast(t *testing.T) {
	exp := MustCompile("$sum(example.value)")

	ast, err := exp.Ast()
	assert.NoError(t, err)
	assert.NotNil(t, ast)
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
		output, err := exp.Evaluate(data)
		assert.NoError(b, err)
		assert.Equal(b, int64(24), output)
	}
}

func BenchmarkExpression_Assign(b *testing.B) {
	exp := MustCompile("$greet")

	for i := 0; i < b.N; i++ {
		err := exp.Assign("greet", "Hello world")
		assert.NoError(b, err)
	}
}

func BenchmarkExpression_RegisterFunction(b *testing.B) {
	exp := MustCompile("$greet()")

	for i := 0; i < b.N; i++ {
		err := exp.RegisterFunction("greet", func(f *Focus, args ...any) (any, error) { return "Hello world", nil })
		assert.NoError(b, err)
	}
}

func BenchmarkExpression_Ast(b *testing.B) {
	exp := MustCompile("$sum(example.value)")

	for i := 0; i < b.N; i++ {
		ast, err := exp.Ast()
		assert.NoError(b, err)
		assert.NotNil(b, ast)
	}
}
