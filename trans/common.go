package trans

import (
	"fmt"
	"strings"
	"test_api/cache"
	"test_api/conf"
	"test_api/fn"
	"test_api/tool"
)

func TransValue(value any) any {
	value = TransConfigValue(value)
	value = TransCacheValue(value)
	value = TransFnValue(value)
	value = TransExprValue(value)

	return value
}

func TransConfigValue(value any) any {
	if valueType, ok := value.(string); ok {
		count := strings.Count(valueType, "$")

		if count > 0 {
			for i := 0; i < count; i++ {
				if valueType, ok = value.(string); ok {
					if strings.Contains(valueType, "$") {
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
					} else {
						break
					}
				}
			}
		}
	}

	return value
}

func TransCacheValue(value any) any {
	if valueType, ok := value.(string); ok {
		count := strings.Count(valueType, "$cache")
		if count > 0 {
			for i := 0; i < count; i++ {
				if valueType, ok = value.(string); ok {
					if strings.Contains(valueType, "$cache") {
						for k, v := range cache.Cache {
							if strings.Contains(valueType, "$cache."+k) {
								if len(valueType) >= len("$cache."+k) {
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
					} else {
						break
					}
				}
			}
		}

	}

	return value
}

func TransFnValue(value any) any {
	if valueType, ok := value.(string); ok {
		count := strings.Count(valueType, "fn.")
		if count > 0 {
			for i := 0; i < count; i++ {
				if valueType, ok = value.(string); ok {
					if strings.Contains(valueType, "fn.") && strings.Contains(valueType, ")") {
						start := strings.Index(valueType, "fn.")
						end := strings.Index(valueType, ")")
						result, err := fn.NewFuncManager().Call(valueType[start : end+1])
						// fmt.Println("result:", result, " err:", err)
						if err != nil {
							return value
						}

						if len(valueType) > end+1 {
							value = valueType[0:start] + fmt.Sprintf("%v", result) + valueType[end+1:]
						} else {
							if start > 0 {
								value = valueType[0:start] + fmt.Sprintf("%v", result)
							} else {
								value = result
							}
						}
					} else {
						break
					}
				}
			}
		}
	}

	return value
}

func TransExprValue(value any) any {
	if valueType, ok := value.(string); ok {
		count := strings.Count(valueType, "expr(")
		if count > 0 {
			for i := 0; i < count; i++ {
				if valueType, ok = value.(string); ok {
					if strings.Contains(valueType, "expr(") && strings.Contains(valueType, ")") {
						result, err := tool.EvaluateWithGoval(valueType[strings.Index(valueType, "expr(")+5 : strings.Index(valueType, ")")-1])

						if err == nil {
							value = result
						}
					} else {
						break
					}
				}
			}
		}
	}

	return value
}
