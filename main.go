package main

import (
	"test_api/conf"
	"test_api/exec"
	"os"
)

func main() {
	conf.Init()

	args := os.Args

	file := args[1]

	exec.ReadFile(file)
}
