package tool

import (
	"fmt"

	expression "github.com/expr-lang/expr"
)

func EvaluateWithGoval(expr string) (any, error) {
	fmt.Println("expr:", expr)
	program, err := expression.Compile(expr)
	fmt.Println("program:", program)
	if err != nil {
		return nil, err
	}

	return expression.Run(program, nil)
}
