package trans

import (
	"fmt"
	"strings"
	"test_api/cache"
	"test_api/conf"
	"test_api/fn"
	"test_api/tool"
)

func TransValueParams(params map[string]any) map[string]any {
	p := cacheParams(params)
	p = customConfigParams(p)
	p = fnParams(p)

	return p
}

func cacheParams(params map[string]any) map[string]any {
	for key, value := range params {
		if valueType, ok := value.(string); ok {
			for k, v := range cache.Cache {
				if strings.Contains(valueType, "$cache."+k) {
					if len(valueType) >= len("$cache."+k) {
						vs, err := tool.AnyToString(v)

						if err != nil {
							vs = ""
						}

						params[key] = strings.ReplaceAll(value.(string), "$cache."+k, vs)
					} else {
						params[key] = v
					}
				}
			}
		}
	}
	return params
}

func customConfigParams(params map[string]any) map[string]any {
	for key, value := range params {
		if valueType, ok := value.(string); ok {
			for k, v := range conf.Config {
				if strings.Contains(valueType, "$"+k) {
					if len(valueType) >= len("$"+k) {
						vs, err := tool.AnyToString(v)

						if err != nil {
							vs = ""
						}

						params[key] = strings.ReplaceAll(value.(string), "$"+k, vs)
					} else {
						params[key] = v
					}
				}
			}
		}
	}

	return params
}

func fnParams(params map[string]any) map[string]any {
	for key, value := range params {
		if valueType, ok := value.(string); ok {
			if strings.Contains(valueType, "fn.") && strings.Contains(valueType, ")") {
				start := strings.Index(valueType, "fn.")
				end := strings.Index(valueType, ")")
				result, err := fn.NewFuncManager().Call(valueType[start : end+1])

				if err != nil {
					params[key] = value
				}

				if len(valueType) > end+1 {
					params[key] = valueType[0:start+1] + fmt.Sprintf("%v", result) + valueType[end+1:]
				} else {
					if start > 0 {
						params[key] = valueType[0:start+1] + fmt.Sprintf("%v", result)
					} else {
						params[key] = result
					}
				}
			}
		}
	}

	return params
}
