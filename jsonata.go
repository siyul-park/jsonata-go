//go:generate npm install jsonata
package jsonata

import (
	_ "embed"
	"reflect"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/iancoleman/strcase"
	"github.com/mitchellh/mapstructure"
)

type (
	Options struct {
		Recover bool
	}

	ExprNode struct {
		Type        string
		Value       any
		Position    int
		Arguments   []*ExprNode
		Name        string
		Procedure   *ExprNode
		Steps       []*ExprNode
		Expressions []*ExprNode
		Stages      []*ExprNode
		Lhs         []*ExprNode
		Rhs         *ExprNode
	}

	Environment struct {
		async     bool
		timestamp time.Time
		vm        *goja.Runtime
		value     *goja.Object
		mu        sync.Mutex
	}

	Focus struct {
		environment *Environment
		input       any
		vm          *goja.Runtime
		value       *goja.Object
		mu          sync.Mutex
	}

	Expression struct {
		vm    *goja.Runtime
		value *goja.Object
		mu    sync.Mutex
	}

	fieldNameMapper struct{}
)

var _ goja.FieldNameMapper = &fieldNameMapper{}

var (
	//go:embed node_modules/jsonata/jsonata.min.js
	source  string
	program = goja.MustCompile("jsonata.min.js", source, true)
)

func init() {
	source = ""
}

func MustCompile(str string, opts ...Options) *Expression {
	exp, err := Compile(str, opts...)
	if err != nil {
		panic(err)
	}
	return exp
}

func Compile(str string, opts ...Options) (*Expression, error) {
	vm := goja.New()

	vm.SetFieldNameMapper(&fieldNameMapper{})

	module := vm.NewObject()
	exports := vm.NewObject()

	if err := vm.Set("module", module); err != nil {
		return nil, err
	} else if err := vm.Set("exports", exports); err != nil {
		return nil, err
	} else if err := module.Set("exports", exports); err != nil {
		return nil, err
	} else if _, err := vm.RunProgram(program); err != nil {
		return nil, err
	}

	module = vm.Get("module").ToObject(vm)

	opt := map[string]any{}
	for _, v := range opts {
		opt["recover"] = v.Recover
	}

	compile, _ := goja.AssertFunction(module.Get("exports"))

	if exp, err := compile(goja.Undefined(), vm.ToValue(str), vm.ToValue(opt)); err != nil {
		return nil, err
	} else {
		return &Expression{
			vm:    vm,
			value: exp.ToObject(vm),
		}, nil
	}
}

func (e *Expression) Evaluate(input any, bindings ...map[string]any) (any, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	jsInput := e.vm.ToValue(input)
	jsBindings := goja.Undefined()
	if len(bindings) > 0 {
		jsBindings = e.vm.ToValue(bindings[len(bindings)-1])
	}

	evaluate, _ := goja.AssertFunction(e.value.Get("evaluate"))
	if output, err := evaluate(e.value, jsInput, jsBindings); err != nil {
		return nil, err
	} else {
		return output.Export().(*goja.Promise).Result().Export(), nil
	}
}

func (e *Expression) Assign(name string, value any) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	assign, _ := goja.AssertFunction(e.value.Get("assign"))
	_, err := assign(e.value, e.vm.ToValue(name), e.vm.ToValue(value))
	return err
}

func (e *Expression) RegisterFunction(name string, implementation func(f *Focus, args ...any) (any, error), signatures ...string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	jsName := e.vm.ToValue(name)
	jsImplementation := e.vm.ToValue(func(call goja.FunctionCall) goja.Value {
		f := &Focus{
			value: call.This.ToObject(e.vm),
			vm:    e.vm,
		}

		args := make([]any, len(call.Arguments))
		for i, arg := range call.Arguments {
			args[i] = arg.Export()
		}
		if v, err := implementation(f, args...); err != nil {
			panic(err)
		} else {
			return e.vm.ToValue(v)
		}
	})
	jsSignature := goja.Undefined()
	if len(signatures) > 0 {
		jsSignature = e.vm.ToValue(signatures[len(signatures)-1])
	}

	registerFunction, _ := goja.AssertFunction(e.value.Get("registerFunction"))
	_, err := registerFunction(e.value, jsName, jsImplementation, jsSignature)
	return err
}

func (e *Expression) Ast() (*ExprNode, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	node := &ExprNode{}
	ast, _ := goja.AssertFunction(e.value.Get("ast"))
	if v, err := ast(e.value); err != nil {
		return nil, err
	} else if err := mapstructure.Decode(v.Export(), &node); err != nil {
		return nil, err
	}
	return node, nil
}

func (f *Focus) Environment() *Environment {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.environment == nil {
		f.environment = &Environment{
			value: f.vm.Get("environment").ToObject(f.vm),
			vm:    f.vm,
		}
	}

	return f.environment
}

func (f *Focus) Input() any {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.input == nil {
		f.input = f.vm.Get("input").Export()
	}
	return f.input
}

func (e *Environment) Async() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.async {
		e.async = e.value.Get("async").ToBoolean()
	}

	return e.async
}

func (e *Environment) Timestamp() time.Time {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.timestamp == (time.Time{}) {
		e.timestamp = e.value.Get("timestamp").Export().(time.Time)
	}

	return e.timestamp
}

func (e *Environment) Bind(s string, a any) error {
	bind, _ := goja.AssertFunction(e.value.Get("bind"))
	_, err := bind(goja.Undefined(), e.vm.ToValue(s), e.vm.ToValue(a))
	return err
}

func (e *Environment) Lookup(s string) (any, error) {
	lookup, _ := goja.AssertFunction(e.value.Get("lookup"))
	if v, err := lookup(goja.Undefined(), e.vm.ToValue(s)); err != nil {
		return nil, err
	} else {
		return v.Export(), nil
	}
}

func (*fieldNameMapper) FieldName(t reflect.Type, f reflect.StructField) string {
	return strcase.ToLowerCamel(f.Name)
}

func (*fieldNameMapper) MethodName(t reflect.Type, m reflect.Method) string {
	return strcase.ToLowerCamel(m.Name)
}
