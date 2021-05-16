package evaluator

import (
	govaluate "gowasmeval/evaluator/govaluate"
	"log"
	"syscall/js"
)

type GenericEvaluator interface {
	EvaluateFn(string) (interface{}, error)
}

type NormalEvaluator struct {
	EvaluateFn func(string) (interface{}, error)
}

func InitAndBindEvaluator() js.Func {

	evaluator := &NormalEvaluator{
		EvaluateFn: govaluate.InitGovaluate(),
	}

	return evaluator.Bind()
}

func (e *NormalEvaluator) Bind() js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		expression := args[0].String()
		result, err := e.EvaluateFn(expression)
		if isErr(err) {
			return false
		}
		return js.ValueOf(result)
	})
}

func isErr(err error) bool {
	if err != nil {
		log.Printf("error encountered: %s\n", err.Error())
		return true
	}
	return false
}
