package evaluator

import (
	govaluate "gowasmeval/evaluator/govaluate"
	"syscall/js"
)

type GenericEvaluator interface {
	Bind() js.Func
}

type NormalEvaluator struct {
	EvaluateFn func(string) (interface{}, error)
}

// NewEvaluator returns the evaluator object
func NewEvaluator() GenericEvaluator {

	evaluator := &NormalEvaluator{
		EvaluateFn: govaluate.InitGovaluate(),
	}

	return evaluator
}

// Bind creates the wrapper go function exposed to js
func (e *NormalEvaluator) Bind() js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		expression := args[0].String()
		result, err := e.EvaluateFn(expression)
		if isErr(err) {
			return js.ValueOf(err.Error())
		}
		return js.ValueOf(result)
	})
}

func isErr(err error) bool {
	if err != nil {
		//log.Printf("error encountered: %s\n", err.Error())
		return true
	}
	return false
}
