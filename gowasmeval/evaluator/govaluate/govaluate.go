package govaluate

import (
	govaluate "github.com/Knetic/govaluate"
)

// InitGovaluate returns function that can evaulate expression
func InitGovaluate() func(string) (interface{}, error) {

	return func(expression string) (interface{}, error) {
		var evaluable_expression *govaluate.EvaluableExpression
		var result interface{}
		var err error

		evaluable_expression, err = govaluate.NewEvaluableExpression(expression)
		if err == nil {
			result, err = evaluable_expression.Evaluate(nil)
			if err == nil {
				return result, nil
			}
		}
		return nil, err
	}
}
