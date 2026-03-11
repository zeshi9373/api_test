package exec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"test_api/cache"
	"test_api/tool"
	"test_api/trans"
)

// TestRunner 测试运行器
type TestRunner struct {
	TestCase *APITestCase
	BaseURL  string
	Cache    map[string]any
}

// NewTestRunner 创建测试运行器
func NewTestRunner(testCase *APITestCase) *TestRunner {
	return &TestRunner{
		TestCase: testCase,
		Cache:    make(map[string]any),
	}
}

// executeAssertions 执行断言
func (r *TestRunner) executeAssertions(respData map[string]any) bool {
	errors := make([]string, 0)
	// assertEquals 断言
	for path, expected := range r.TestCase.Expect.AssertEquals {
		if err := r.assertPathEquals(respData, path, expected); err != nil {
			fmt.Printf("❌断言失败: %s \n", err.Error())
			errors = append(errors, fmt.Sprintf("❌断言失败: %s \n", err.Error()))
		}
	}

	// assertNotEquals 断言
	for path, notExpected := range r.TestCase.Expect.AssertNotEquals {
		if err := r.assertPathNotEquals(respData, path, notExpected); err != nil {
			fmt.Printf("❌断言失败: %s \n", err.Error())
			errors = append(errors, fmt.Sprintf("❌断言失败: %s \n", err.Error()))
		}
	}

	// assertContains 断言
	for path, expected := range r.TestCase.Expect.AssertContains {
		if err := r.assertPathContains(respData, path, expected); err != nil {
			fmt.Printf("❌断言失败: %s \n", err.Error())
			errors = append(errors, fmt.Sprintf("❌断言失败: %s \n", err.Error()))
		}
	}

	// assertMatches 断言
	for path, pattern := range r.TestCase.Expect.AssertMatches {
		if err := r.assertPathMatches(respData, path, pattern); err != nil {
			fmt.Printf("❌断言失败: %s \n", err.Error())
			errors = append(errors, fmt.Sprintf("❌断言失败: %s \n", err.Error()))
		}
	}

	// assertType 断言
	for path, typeName := range r.TestCase.Expect.AssertType {
		if err := r.assertPathType(respData, path, typeName); err != nil {
			fmt.Printf("❌断言失败: %s \n", err.Error())
			errors = append(errors, fmt.Sprintf("❌断言失败: %s \n", err.Error()))
		}
	}

	// assertLength 断言
	for path, length := range r.TestCase.Expect.AssertLength {
		if err := r.assertPathLength(respData, path, length); err != nil {
			fmt.Printf("❌断言失败: %s \n", err.Error())
			errors = append(errors, fmt.Sprintf("❌断言失败: %s \n", err.Error()))
		}
	}

	if len(errors) > 0 {
		TestTotal.Fail++
		p, _ := json.Marshal(r.TestCase.Params)
		TestTotal.FailDetal = append(TestTotal.FailDetal, fmt.Sprintf("测试文件：%s\n测试接口：%s (%s)\n请求参数：%s\n%s", r.TestCase.FilePath, r.TestCase.API, r.TestCase.Description, string(p), strings.Join(errors, "")))
		return false
	}

	TestTotal.Pass++
	return true
}

// getValueByPath 通过路径获取值
func (r *TestRunner) getValueByPath(data map[string]any, path string) (any, error) {
	parts := strings.Split(path, ".")
	var current any = data

	for k := 0; k < len(parts); k++ {
		if tool.IsNumericByRune(parts[k]) {
			i, _ := strconv.Atoi(parts[k])
			if _, ok := current.([]any); ok {
				if i >= len(current.([]any)) {
					return nil, fmt.Errorf("路径 %s 数据不存在", path)
				}

				current = current.([]any)[i]
				continue
			} else {
				return nil, fmt.Errorf("路径 %s 数据不存在", path)
			}
		}

		m, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("路径 %s 无法访问: 不是对象类型", parts[k])
		}

		val, exists := m[parts[k]]
		if !exists {
			return nil, fmt.Errorf("路径不存在: %s", path)
		}

		current = val
	}

	return current, nil
}

// assertPathEquals 断言路径值相等
func (r *TestRunner) assertPathEquals(data map[string]any, path string, expected any) error {
	expected = r.getDataValueByPath(data, expected)
	actual, err := r.getValueByPath(data, path)

	if err != nil {
		return fmt.Errorf("assertEquals %s: %v", path, err)
	}

	// 特殊处理浮点数（JSON 中的数字都是 float64）
	if expectedInt, ok := expected.(int); ok {
		if actualFloat, ok := actual.(float64); ok {
			if float64(expectedInt) == actualFloat {
				fmt.Printf("✅ assertEquals %s: %v \n", path, expected)
				return nil
			} else {
				return fmt.Errorf("assertEquals %s: 期望 %v, 实际 %v", path, expected, actual)
			}
		}
	}

	if !deepEqual(actual, expected) {
		return fmt.Errorf("assertEquals %s: 期望 %v, 实际 %v", path, expected, actual)
	}

	fmt.Printf("✅ assertEquals %s: %v \n", path, expected)
	return nil
}

// assertPathNotEquals 断言路径值不相等
func (r *TestRunner) assertPathNotEquals(data map[string]any, path string, notExpected any) error {
	notExpected = r.getDataValueByPath(data, notExpected)
	actual, err := r.getValueByPath(data, path)
	if err != nil {
		return fmt.Errorf("assertNotEquals %s: %v", path, err)
	}

	if deepEqual(actual, notExpected) {
		return fmt.Errorf("assertNotEquals %s: 值不应为 %v, 实际 %v", path, notExpected, actual)
	}

	if notExpected == "" {
		notExpected = "空"
	}

	fmt.Printf("✅ assertNotEquals %s: 不为 %v \n", path, notExpected)
	return nil
}

// assertPathContains 断言路径值包含
func (r *TestRunner) assertPathContains(data map[string]any, path string, expected any) error {
	actual, err := r.getValueByPath(data, path)
	if err != nil {
		return fmt.Errorf("assertContains %s: %v", path, err)
	}

	switch actualVal := actual.(type) {
	case string:
		expectedStr, ok := expected.(string)
		if !ok {
			return fmt.Errorf("assertContains %s: 期望值应为字符串", path)
		}

		if !strings.Contains(expectedStr, "@@") && !strings.Contains(expectedStr, "||") {
			expected = r.getDataValueByPath(data, expectedStr)
			// fmt.Println("expected:", expected)
			if !strings.Contains(actualVal, fmt.Sprintf("%v", expected)) {
				return fmt.Errorf("assertContains %s: 期望包含 '%s', 实际 '%s'",
					path, expectedStr, actualVal)
			}
		} else {
			// 是否判断多个值
			// @@ 多个值都包含
			// || 任意一个值包含
			if strings.Contains(expectedStr, "@@") {
				expectedSlice := strings.Split(expectedStr, "@@")

				for _, v := range expectedSlice {
					e := r.getDataValueByPath(data, v)

					if !strings.Contains(actualVal, fmt.Sprintf("%v", e)) {
						return fmt.Errorf("assertContains @@ %s: 期望包含 '%s', 实际 '%s'",
							path, expectedStr, actualVal)
					}
				}
			}

			if strings.Contains(expectedStr, "||") {
				var isContains bool
				expectedSlice := strings.Split(expectedStr, "||")

				for _, v := range expectedSlice {
					e := r.getDataValueByPath(data, v)

					if strings.Contains(actualVal, fmt.Sprintf("%v", e)) {
						isContains = true
						break
					}
				}

				if !isContains {
					return fmt.Errorf("assertContains || %s: 期望包含 '%s', 实际 '%s'",
						path, expectedStr, actualVal)
				}
			}
		}
	default:
		return fmt.Errorf("assertContains %s: 仅支持字符串类型", path)
	}

	fmt.Printf("✅ assertContains %s: 包含 %v \n", path, expected)
	return nil
}

// assertPathMatches 断言路径值匹配正则
func (r *TestRunner) assertPathMatches(data map[string]any, path, pattern string) error {
	actual, err := r.getValueByPath(data, path)
	if err != nil {
		return fmt.Errorf("assertMatches %s: %v", path, err)
	}

	actualStr, ok := actual.(string)
	if !ok {
		return fmt.Errorf("assertMatches %s: 值应为字符串类型", path)
	}

	pattern = fmt.Sprintf("%v", r.getDataValueByPath(data, pattern))
	matched, err := regexp.MatchString(pattern, actualStr)
	if err != nil {
		return fmt.Errorf("assertMatches %s: 正则表达式错误: %v", path, err)
	}

	if !matched {
		return fmt.Errorf("assertMatches %s: 期望匹配 '%s', 实际 '%s'",
			path, pattern, actualStr)
	}

	fmt.Printf("✅ assertMatches %s: 匹配 %s \n", path, pattern)
	return nil
}

// assertPathType 断言路径值类型
func (r *TestRunner) assertPathType(data map[string]any, path, typeName string) error {
	actual, err := r.getValueByPath(data, path)
	if err != nil {
		return fmt.Errorf("assertType %s: %v", path, err)
	}

	var actualType string
	switch actual.(type) {
	case string:
		actualType = "string"
	case float64:
		actualType = "number"
	case bool:
		actualType = "boolean"
	case map[string]any:
		actualType = "object"
	case []any:
		actualType = "array"
	case nil:
		actualType = "null"
	default:
		actualType = "unknown"
	}

	if actualType != typeName {
		return fmt.Errorf("assertType %s: 期望类型 %s, 实际类型 %s",
			path, typeName, actualType)
	}

	fmt.Printf("✅ assertType %s: %s \n", path, typeName)
	return nil
}

// assertPathLength 断言路径值长度
func (r *TestRunner) assertPathLength(data map[string]any, path string, expectedLength int) error {
	actual, err := r.getValueByPath(data, path)
	if err != nil {
		return fmt.Errorf("assertLength %s: %v", path, err)
	}

	var actualLength int
	switch v := actual.(type) {
	case string:
		actualLength = len(v)
	case []any:
		actualLength = len(v)
	case map[string]any:
		actualLength = len(v)
	default:
		return fmt.Errorf("assertLength %s: 不支持的类型 %T", path, actual)
	}

	if actualLength != expectedLength {
		return fmt.Errorf("assertLength %s: 期望长度 %d, 实际长度 %d",
			path, expectedLength, actualLength)
	}

	fmt.Printf("✅ assertLength %s: %d \n", path, expectedLength)
	return nil
}

// AssertTrue 断言真
func (r *TestRunner) AssertTrue(data map[string]any, path string) error {
	actual, err := r.getValueByPath(data, path)
	if err != nil {
		return fmt.Errorf("assertTrue %s: %v", path, err)
	}

	if actual != true {
		return fmt.Errorf("assertTrue %s: 值应为 true, 实际 %v", path, actual)
	}

	fmt.Printf("✅ assertTrue %s: true \n", path)
	return nil
}

// AssertFalse 断言假
func (r *TestRunner) AssertFalse(data map[string]any, path string) error {
	actual, err := r.getValueByPath(data, path)
	if err != nil {
		return fmt.Errorf("assertFalse %s: %v", path, err)
	}

	if actual != false {
		return fmt.Errorf("assertFalse %s: 值应为 false, 实际 %v", path, actual)
	}

	fmt.Printf("✅ assertFalse %s: false \n", path)
	return nil
}

// AssertIn 断言路径值在列表中
func (r *TestRunner) AssertIn(data map[string]any, path string, expected []any) error {
	actual, err := r.getValueByPath(data, path)
	if err != nil {
		return fmt.Errorf("assertIn %s: %v", path, err)
	}

	for _, v := range expected {
		v = r.getDataValueByPath(data, v)
		if v == actual {
			fmt.Printf("✅ assertIn %s: %v \n", path, v)
			return nil
		}
	}

	return fmt.Errorf("assertIn %s: 值不在列表中, 列表: %v", path, expected)
}

// deepEqual 深度比较
func deepEqual(a, b any) bool {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return string(aJSON) == string(bJSON)
}

// cacheData 缓存数据
func (r *TestRunner) cacheData(data map[string]any) {
	for cacheKey, dataPath := range r.TestCase.Cache {
		value, err := r.getValueByPath(data, dataPath)
		if err == nil && value != nil {
			cache.Cache[cacheKey] = tool.SliceMapToJson(value)
		}
	}
}
func (r *TestRunner) getDataValueByPath(data map[string]any, expected any) any {
	for {
		if expectedStr, ok := expected.(string); ok {
			if strings.Contains(expectedStr, "dataValue(") {
				if strings.Index(expectedStr, "dataValue(") == 0 && strings.Index(expectedStr, ")") == len(expectedStr)-1 {
					expectedStrPath := expectedStr[strings.Index(expectedStr, "dataValue(")+len("dataValue(") : strings.Index(expectedStr, ")")]
					expectedValue, _ := r.getValueByPath(data, expectedStrPath)
					expected = trans.TransValue(expectedValue)
				} else {
					var expectedStrLeft, expectedStrRight string

					if strings.Index(expectedStr, "dataValue(") > 0 {
						expectedStrLeft = expectedStr[0:strings.Index(expectedStr, "dataValue(")]
					}

					expectedStrReMain := expectedStr[strings.Index(expectedStr, "dataValue("):]
					// fmt.Println("expectedStrReMain:", expectedStrReMain)
					// fmt.Println("expectedStrReMainIndex1:", strings.Index(expectedStrReMain, "dataValue(")+len("dataValue("))
					// fmt.Println("expectedStrReMainIndex2:", strings.Index(expectedStrReMain, ")"))
					expectedStrPath := expectedStrReMain[strings.Index(expectedStrReMain, "dataValue(")+len("dataValue(") : strings.Index(expectedStrReMain, ")")-strings.Index(expectedStrReMain, "dataValue(")]
					// fmt.Println("expectedStrPath:", expectedStrPath)

					if strings.Index(expectedStrReMain, ")") < len(expectedStrReMain)-1 {
						expectedStrRight = expectedStrReMain[strings.Index(expectedStrReMain, ")")+1:]
					}

					value, _ := r.getValueByPath(data, expectedStrPath)
					expected = trans.TransValue(expectedStrLeft + fmt.Sprintf("%v", value) + expectedStrRight)
				}
			} else {
				break
			}
		} else {
			break
		}
	}

	return trans.TransValue(expected)
}
