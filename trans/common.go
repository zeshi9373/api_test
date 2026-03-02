package trans

import (
	"fmt"
	"strings"
	"test_api/conf"
	"test_api/fn"
	"test_api/tool"
)

func TransValue(value any) any {
	value = TransConfigValue(value)
	value = TransCacheValue(value)
	value = TransFnValue(value)

	if valueType, ok := value.(string); ok {
		if strings.Contains(valueType, "expr(") && strings.Contains(valueType, ")") {
			result, err := tool.EvaluateWithGoval(valueType[strings.Index(valueType, "expr(")+5 : strings.Index(valueType, ")")-1])

			if err != nil {
				return value
			}

			return result
		}
	}

	return value
}

func TransConfigValue(value any) any {
	if valueType, ok := value.(string); ok {
		for k, v := range conf.Config {
			if strings.Contains(valueType, "$"+k) {
				if len(valueType) == len("$"+k) {
					vs, err := tool.AnyToString(v)

					if err != nil {
						vs = ""
					}

					value = strings.ReplaceAll(value.(string), "$"+k, vs)
				} else {
					value = v
				}
			}
		}
	}

	return value
}

func TransCacheValue(value any) any {
	if valueType, ok := value.(string); ok {
		for k, v := range conf.Config {
			if strings.Contains(valueType, "$cache."+k) {
				if len(valueType) == len("$cache."+k) {
					vs, err := tool.AnyToString(v)

					if err != nil {
						vs = ""
					}

					value = strings.ReplaceAll(value.(string), "$cache."+k, vs)
				} else {
					value = v
				}
			}
		}
	}

	return value
}

func TransFnValue(value any) any {
	if valueType, ok := value.(string); ok {
		if strings.Contains(valueType, "fn.") && strings.Contains(valueType, ")") {
			start := strings.Index(valueType, "fn.")
			end := strings.Index(valueType, ")")
			result, err := fn.NewFuncManager().Call(valueType[start : end+1])

			if err != nil {
				return value
			}

			if len(valueType) > end+1 {
				value = valueType[0:start+1] + fmt.Sprintf("%v", result) + valueType[end+1:]
			} else {
				if start > 0 {
					value = valueType[0:start+1] + fmt.Sprintf("%v", result)
				} else {
					value = result
				}
			}
		}
	}

	return value
}
