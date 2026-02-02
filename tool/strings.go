package tool

import "unicode"

func IsNumericByRune(s string) bool {
	if s == "" {
		return false
	}

	// 检查每个字符是否为数字
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
