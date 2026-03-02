package trans

import (
	"fmt"
	"strings"
	"test_api/cache"
	"test_api/conf"
	"test_api/fn"
	"test_api/tool"
)

func TransValueHeaders(headers map[string]string) map[string]string {
	h := cacheHeaders(headers)
	h = customConfigHeaders(h)
	h = fnHeaders(h)

	return h
}

func cacheHeaders(headers map[string]string) map[string]string {
	for key, value := range headers {
		for k, v := range cache.Cache {
			if strings.Contains(value, "$cache."+k) {
				vs, err := tool.AnyToString(v)

				if err != nil {
					vs = ""
				}

				headers[key] = strings.ReplaceAll(value, "$cache."+k, vs)
			}
		}
	}
	return headers
}

func customConfigHeaders(headers map[string]string) map[string]string {
	for key, value := range headers {
		for k, v := range conf.Config {
			if strings.Contains(value, "$"+k) {
				vs, err := tool.AnyToString(v)

				if err != nil {
					vs = ""
				}

				headers[key] = strings.ReplaceAll(value, "$"+k, vs)
			}
		}
	}
	return headers
}

func fnHeaders(headers map[string]string) map[string]string {
	for key, value := range headers {
		if strings.Contains(value, "fn.") && strings.Contains(value, ")") {
			start := strings.Index(value, "fn.")
			end := strings.Index(value, ")")
			result, err := fn.NewFuncManager().Call(value[start : end+1])

			if err != nil {
				headers[key] = value
			}

			if len(value) > end+1 {
				headers[key] = value[0:start+1] + fmt.Sprintf("%v", result) + value[end+1:]
			} else {
				if start > 0 {
					headers[key] = value[0:start+1] + fmt.Sprintf("%v", result)
				} else {
					headers[key] = fmt.Sprintf("%v", result)
				}
			}

		}
	}

	return headers
}
