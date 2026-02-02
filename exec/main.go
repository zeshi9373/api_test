package exec

import (
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

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
	for _, file := range testFile.Include {
		GetCase(file)
	}

	time.Sleep(1500 * time.Millisecond)
}
