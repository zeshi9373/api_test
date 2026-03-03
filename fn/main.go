package fn

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// 函数管理器
type FuncManager struct {
	funcs map[string]interface{}
}

func NewFuncManager() *FuncManager {
	return &FuncManager{
		funcs: map[string]interface{}{
			"fn.RandInt":         RandInt,
			"fn.Time":            Time,
			"fn.TimeNanos":       TimeNanos,
			"fn.TimeMillis":      TimeMillis,
			"fn.TimeMicros":      TimeMicros,
			"fn.Date":            Date,
			"fn.DateTime":        DateTime,
			"fn.RandString":      RandString,
			"fn.SearchListValue": SearchListValue,
		},
	}
}

func (fm *FuncManager) Call(funcExpr string) (interface{}, error) {
	// 解析函数表达式：函数名(参数1, 参数2, ...)
	re := regexp.MustCompile(`^(\w+\.\w+)\((.*)\)$`)
	matches := re.FindStringSubmatch(funcExpr)

	if matches == nil {
		return funcExpr, nil
	}

	funcName := matches[1]
	argsStr := matches[2]

	fn, exists := fm.funcs[funcName]
	if !exists {
		return nil, fmt.Errorf("function %s not found", funcName)
	}

	// 解析参数
	var args []reflect.Value
	if argsStr != "" {
		argStrs := splitArgs(argsStr)
		fnType := reflect.TypeOf(fn)

		for i, argStr := range argStrs {
			if i >= fnType.NumIn() {
				break
			}

			paramType := fnType.In(i)
			arg, err := parseArg(argStr, paramType)
			if err != nil {
				return nil, err
			}
			args = append(args, reflect.ValueOf(arg))
		}
	}

	// 调用函数
	results := reflect.ValueOf(fn).Call(args)
	if len(results) > 0 {
		return results[0].Interface(), nil
	}

	return nil, nil
}

// 辅助函数
func splitArgs(s string) []string {
	// 简单分割，不支持嵌套括号
	var args []string
	var current strings.Builder
	parenDepth := 0

	for _, r := range s {
		switch r {
		case '(':
			parenDepth++
			current.WriteRune(r)
		case ')':
			parenDepth--
			current.WriteRune(r)
		case '[':
			parenDepth++
			current.WriteRune(r)
		case ']':
			parenDepth--
			current.WriteRune(r)
		case '{':
			parenDepth++
			current.WriteRune(r)
		case '}':
			parenDepth--
			current.WriteRune(r)
		case ',':
			if parenDepth == 0 {
				args = append(args, strings.TrimSpace(current.String()))
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		args = append(args, strings.TrimSpace(current.String()))
	}

	return args
}

func parseArg(argStr string, typ reflect.Type) (interface{}, error) {
	switch typ.Kind() {
	case reflect.Int:
		return strconv.Atoi(argStr)
	case reflect.Int64:
		return strconv.ParseInt(argStr, 10, 64)
	case reflect.String:
		return argStr, nil
	case reflect.Float64:
		return strconv.ParseFloat(argStr, 64)
	default:
		return argStr, nil
	}
}
