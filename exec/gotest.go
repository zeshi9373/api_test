package exec

import (
	"test_api/conf"
	"test_api/logger"
	"test_api/tool"
	"test_api/trans"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	_ "test_api/fn"

	"gopkg.in/yaml.v2"
)

type APITestCase struct {
	ApiDomain string            `yaml:"api_domain"`
	API       string            `yaml:"api"`
	Method    string            `yaml:"method"`
	Headers   map[string]string `yaml:"headers"`
	Params    map[string]any    `yaml:"params"`
	Expect    Expect            `yaml:"expect"`
	Cache     map[string]string `yaml:"cache"`
}

type Expect struct {
	AssertEquals    map[string]any    `yaml:"assertEquals,omitempty"`
	AssertNotEquals map[string]any    `yaml:"assertNotEquals,omitempty"`
	AssertContains  map[string]any    `yaml:"assertContains,omitempty"`
	AssertMatches   map[string]string `yaml:"assertMatches,omitempty"` // 正则匹配
	AssertLength    map[string]int    `yaml:"assertLength,omitempty"`  // 长度断言
	AssertType      map[string]string `yaml:"assertType,omitempty"`    // 类型断言
	AssertTrue      []string          `yaml:"assertTrue,omitempty"`    // 为真断言
	AssertFalse     []string          `yaml:"assertFalse,omitempty"`   // 为假断言
	AssertIn        map[string]any    `yaml:"assertIn,omitempty"`
}

func GetCase(filepath string) {
	byteStream, err := os.ReadFile(filepath)

	if err != nil {
		panic("read test case file error")
	}

	var testCase APITestCase
	yaml.Unmarshal(byteStream, &testCase)

	if len(testCase.API) == 0 {
		panic("api is empty")
	}

	if len(testCase.Method) == 0 {
		panic("method is empty")
	}

	domain := conf.Config["api_domain"].(string)

	if len(testCase.ApiDomain) > 0 {
		domain = trans.TransConfigValue(testCase.ApiDomain).(string)
	}

	timeout := conf.Config["timeout"].(int)
	var response []byte
	var reserr error

	fmt.Printf("接口测试开始 %s \n", testCase.API)

	testCase.Headers = trans.TransValueHeaders(testCase.Headers)
	testCase.Params = trans.TransValueParams(testCase.Params)
	// fmt.Println("Headers:", testCase.Headers)
	// fmt.Println("Params:", testCase.Params)
	// fmt.Println("Assert:", testCase.Expect)
	switch testCase.Method {
	case "GET":
		p, err := tool.MapAnyToString(testCase.Params)

		if err != nil {
			panic("params is not support")
		}

		response, reserr = tool.NewHttp(domain+testCase.API, time.Duration(timeout)*time.Second).Get(testCase.Headers, p)
	case "POST":
		d, err := json.Marshal(testCase.Params)

		if err != nil {
			panic("params is not support")
		}
		response, reserr = tool.NewHttp(domain+testCase.API, time.Duration(timeout)*time.Second).Post(testCase.Headers, d)
	default:
		panic("method is not support")
	}

	if reserr != nil {
		panic("api request " + reserr.Error())
	}

	logger.NewLogger(strings.ReplaceAll(testCase.API, "/", "_")).Info("请求接口日志", logger.Fields{
		"headers":  testCase.Headers,
		"params":   testCase.Params,
		"response": string(response),
	})

	resData := make(map[string]any)
	json.Unmarshal(response, &resData)
	NewTestRunner(&testCase).executeAssertions(resData)
	NewTestRunner(&testCase).cacheData(resData)
}
