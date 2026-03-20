package fn

import "strings"

func RandString(flag string, length int) string {
	var baseString, resultString string

	if strings.Contains(flag, "N") {
		baseString += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	if strings.Contains(flag, "n") {
		baseString += "abcdefghijklmnopqrstuvwxyz"
	}

	if strings.Contains(flag, "1") {
		baseString += "0123456789"
	}

	if strings.Contains(flag, "S") {
		baseString += "!@#$%^&*()_+-=[]{};':|,.<>/?"
	}

	for i := 0; i < length; i++ {
		resultString += string(baseString[RandInt(0, len(baseString)-1)])
	}

	return resultString
}

func ReplaceString(str, old, new string) string {
	return strings.ReplaceAll(str, old, new)
}
