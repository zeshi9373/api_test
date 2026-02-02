package conf

import (
	"test_api/tool"
	"os"

	"gopkg.in/yaml.v2"
)

var Config = make(map[string]any)

func Init() {
	// 初始化变量
	byteStream, err := os.ReadFile("./config.yaml")

	if err != nil {
		panic("read config.yaml file error")
	}

	if err := yaml.Unmarshal(byteStream, &Config); err != nil {
		panic("unmarshal config.yaml file error")
	}

	if !tool.IsNotEmpty(Config["api_domain"]) {
		panic("api_domain is empty")
	}
}
