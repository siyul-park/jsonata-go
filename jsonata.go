//go:generate npm install jsonata
package jsonata

import (
	_ "embed"
	"reflect"
	"sync"

	"github.com/dop251/goja"
	"github.com/iancoleman/strcase"
)

type (
	Options struct {
		Recover bool
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

	opt := map[string]any{}
	for _, v := range opts {
		if v.Recover {
			opt["recover"] = true
		}
	}

	parse, _ := goja.AssertFunction(exports)

	if exp, err := parse(goja.Undefined(), vm.ToValue(str), vm.ToValue(opt)); err != nil {
		return nil, err
	} else {
		return &Expression{
			vm:    vm,
			value: exp.ToObject(vm),
		}, nil
	}
}

func (e *Expression) Evaluate(input any, bindings map[string]any) (any, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	evaluate, _ := goja.AssertFunction(e.value.Get("evaluate"))
	if output, err := evaluate(e.value, e.vm.ToValue(input), e.vm.ToValue(bindings)); err != nil {
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

func (*fieldNameMapper) FieldName(t reflect.Type, f reflect.StructField) string {
	return strcase.ToLowerCamel(f.Name)
}

func (*fieldNameMapper) MethodName(t reflect.Type, m reflect.Method) string {
	return strcase.ToLowerCamel(m.Name)
}
