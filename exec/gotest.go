package exec

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"test_api/conf"
	"test_api/logger"
	"test_api/tool"
	"test_api/trans"
	"time"

	_ "test_api/fn"

	"gopkg.in/yaml.v2"
)

type APITestCase struct {
	FilePath    string            `yaml:"file_path"`
	ApiDomain   string            `yaml:"api_domain"`
	API         string            `yaml:"api"`
	Description string            `yaml:"description"`
	Method      string            `yaml:"method"`
	Headers     map[string]string `yaml:"headers"`
	Params      map[string]any    `yaml:"params"`
	Expect      Expect            `yaml:"expect"`
	Cache       map[string]string `yaml:"cache"`
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
		TestTotal.Fail++
		TestTotal.FailDetal = append(TestTotal.FailDetal, fmt.Sprintf("%s", "测试文件不存在 ："+filepath))
		fmt.Printf("测试文件不存在 ：%s", filepath)
		return
	}

	var testCase APITestCase
	err = yaml.Unmarshal(byteStream, &testCase)

	if err != nil {
		fmt.Println("unmarshal err", err)
		TestTotal.Fail++
		TestTotal.FailDetal = append(TestTotal.FailDetal, fmt.Sprintf("%s", "测试文件格式不支持 ："+filepath))
		fmt.Printf("测试文件格式不支持 ：%s", filepath)
		return
	}

	if len(testCase.API) == 0 {
		TestTotal.Fail++
		TestTotal.FailDetal = append(TestTotal.FailDetal, fmt.Sprintf("%s", "测试文件api为空 ："+filepath))
		fmt.Printf("%s", "测试文件api为空 ："+filepath)
		return
	}

	if len(testCase.Method) == 0 {
		TestTotal.Fail++
		TestTotal.FailDetal = append(TestTotal.FailDetal, fmt.Sprintf("%s", "测试文件method为空 ："+filepath))
		fmt.Printf("%s", "测试文件method为空 ："+filepath)
		return
	}

	domain := conf.Config["api_domain"].(string)

	if len(testCase.ApiDomain) > 0 {
		domain = trans.TransConfigValue(testCase.ApiDomain).(string)
	}

	timeout := conf.Config["timeout"].(int)
	var response []byte
	var reserr error

	p, _ := json.Marshal(testCase.Params)
	testCase.FilePath = filepath
	fmt.Printf("接口测试开始 %s (%s)\n执行文件：%s\n请求参数：%s\n", testCase.API, testCase.Description, filepath, string(p))

	testCase.Headers = trans.TransValueHeaders(testCase.Headers)
	testCase.Params = trans.TransValueParams(testCase.Params)
	// fmt.Println("Headers:", testCase.Headers)
	// fmt.Println("Params:", testCase.Params)
	// fmt.Println("Assert:", testCase.Expect)
	switch testCase.Method {
	case "GET":
		p, err := tool.MapAnyToString(testCase.Params)

		if err != nil {
			TestTotal.Fail++
			TestTotal.FailDetal = append(TestTotal.FailDetal, fmt.Sprintf("%s", "测试文件params格式不支持 ："+filepath))
			return
		}

		response, reserr = tool.NewHttp(domain+testCase.API, time.Duration(timeout)*time.Second).Get(testCase.Headers, p)
	case "POST":
		d, err := json.Marshal(testCase.Params)

		if err != nil {
			TestTotal.Fail++
			TestTotal.FailDetal = append(TestTotal.FailDetal, fmt.Sprintf("%s", "测试文件params格式不支持 ："+filepath))
			return
		}
		response, reserr = tool.NewHttp(domain+testCase.API, time.Duration(timeout)*time.Second).Post(testCase.Headers, d)
	default:
		TestTotal.Fail++
		TestTotal.FailDetal = append(TestTotal.FailDetal, fmt.Sprintf("%s", "测试文件method格式不支持 ："+filepath))
		return
	}

	if reserr != nil {
		TestTotal.Fail++
		TestTotal.FailDetal = append(TestTotal.FailDetal, fmt.Sprintf("测试文件：%s，请求报错：%s", filepath, reserr.Error()))
		return
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
