package exec

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"test_api/conf"
	"test_api/tool"
	"time"

	yaml "gopkg.in/yaml.v2"
)

var (
	TestTotal TestCaseTotal
)

type TestCaseTotal struct {
	Total     int      `yaml:"total"`
	Pass      int      `yaml:"pass"`
	Fail      int      `yaml:"fail"`
	FailDetal []string `yaml:"fail_detal"`
}

type TestFile struct {
	Include []string `yaml:"include"`
}

func ReadFile(filepath string) {
	byteStream, err := os.ReadFile(filepath)

	if err != nil {
		panic("read test file error")
	}

	testFile := TestFile{}
	yaml.Unmarshal(byteStream, &testFile)

	TestTotal.Total = len(testFile.Include)

	for _, file := range testFile.Include {
		GetCase(file)
	}

	// fmt.Println("TestTotal:", TestTotal)
	if robotUrl, ok := conf.Config["feishu_robot"].(string); ok {
		if len(robotUrl) > 0 {
			message := fmt.Sprintf("执行文件：%s \n\n测试用例总数：%d\n通过数：%d\n失败数：%d\n通过率：%s\n失败明细：\n%s", filepath, TestTotal.Total, TestTotal.Pass, TestTotal.Fail, fmt.Sprintf("%.2f", float32(TestTotal.Pass)/float32(TestTotal.Total)*100)+"%", strings.Join(TestTotal.FailDetal, "\n"))
			data := map[string]interface{}{
				"msg_type": "text",
				"content":  map[string]string{"text": message},
			}

			str, _ := json.Marshal(data)
			tool.NewHttp(robotUrl, 3*time.Second).Post(nil, str)
		}
	}

	time.Sleep(1500 * time.Millisecond)
}
