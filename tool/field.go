package tool

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"encoding/json"
)

// IsNotEmpty 判断 any 类型是否不为空
func IsNotEmpty(value any) bool {
	if value == nil {
		return false
	}

	v := reflect.ValueOf(value)

	// 处理指针类型
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false
		}
		// 解引用指针，检查指向的值
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		str := v.String()
		return strings.TrimSpace(str) != ""

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() != 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() != 0

	case reflect.Float32, reflect.Float64:
		return v.Float() != 0

	case reflect.Bool:
		return v.Bool()

	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() > 0

	case reflect.Struct:
		// 特殊处理 time.Time
		if t, ok := value.(time.Time); ok {
			return !t.IsZero()
		}
		// 其他结构体，检查所有字段是否都为零值
		return !isZeroStruct(v)

	case reflect.Interface:
		if v.IsNil() {
			return false
		}
		// 递归检查接口中的值
		return IsNotEmpty(v.Interface())

	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return !v.IsNil()

	default:
		return true
	}
}

// isZeroStruct 检查结构体是否所有字段都为零值
func isZeroStruct(v reflect.Value) bool {
	typ := v.Type()

	// 检查是否为 time.Time 类型
	if typ == reflect.TypeOf(time.Time{}) {
		return v.Interface().(time.Time).IsZero()
	}

	// 检查结构体的每个字段
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := typ.Field(i)

		// 跳过不可导出的字段
		if !fieldType.IsExported() {
			continue
		}

		// 递归检查字段
		if !isZeroValue(field) {
			return false
		}
	}
	return true
}

// isZeroValue 检查值是否为零值
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Struct:
		return isZeroStruct(v)
	default:
		return false
	}
}

// MapAnyToString 基础转换
func MapAnyToString(m map[string]any) (map[string]string, error) {
	if m == nil {
		return nil, nil
	}

	result := make(map[string]string, len(m))

	for key, value := range m {
		if value == nil {
			result[key] = ""
			continue
		}

		switch v := value.(type) {
		case string:
			result[key] = v
		case []byte:
			result[key] = string(v)
		case fmt.Stringer:
			result[key] = v.String()
		case error:
			result[key] = v.Error()
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64:
			result[key] = fmt.Sprintf("%d", v)
		case float32, float64:
			result[key] = fmt.Sprintf("%f", v)
		case bool:
			result[key] = strconv.FormatBool(v)
		case time.Time:
			result[key] = v.Format(time.RFC3339)
		default:
			// 使用反射处理其他类型
			str, err := AnyToString(value)
			if err != nil {
				return nil, fmt.Errorf("key %s: %v", key, err)
			}
			result[key] = str
		}
	}

	return result, nil
}

// AnyToString 通用 any 转 string
func AnyToString(value any) (string, error) {
	if value == nil {
		return "", nil
	}

	rv := reflect.ValueOf(value)

	// 处理指针
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return "", nil
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.String:
		return rv.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 64), nil
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool()), nil
	case reflect.Slice, reflect.Array:
		if rv.Len() == 0 {
			return "[]", nil
		}
		// 转换为 JSON 数组字符串
		var items []string
		for i := 0; i < rv.Len(); i++ {
			itemStr, err := AnyToString(rv.Index(i).Interface())
			if err != nil {
				return "", err
			}
			items = append(items, itemStr)
		}
		return "[" + strings.Join(items, ",") + "]", nil
	case reflect.Map:
		// 转换为 JSON 对象字符串
		if rv.Len() == 0 {
			return "{}", nil
		}
		var pairs []string
		iter := rv.MapRange()
		for iter.Next() {
			keyStr, err := AnyToString(iter.Key().Interface())
			if err != nil {
				return "", err
			}
			valStr, err := AnyToString(iter.Value().Interface())
			if err != nil {
				return "", err
			}
			pairs = append(pairs, fmt.Sprintf("%s:%s", keyStr, valStr))
		}
		return "{" + strings.Join(pairs, ",") + "}", nil
	case reflect.Struct:
		// 特殊处理 time.Time
		if t, ok := value.(time.Time); ok {
			return t.Format(time.RFC3339), nil
		}
		// 其他结构体使用 JSON 格式
		return fmt.Sprintf("%+v", value), nil
	default:
		return fmt.Sprintf("%v", value), nil
	}
}

func SliceMapToJson(value any) any {
	if value == "" {
		return value
	}

	rv := reflect.ValueOf(value)

	// 处理指针
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return value
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Struct:
		// 转换为 JSON 对象字符串
		str, err := json.Marshal(value)

		if err != nil {
			return value
		}
		return string(str)
	default:
		return value
	}
}
