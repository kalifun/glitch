package utils

import (
	"regexp"
	"strings"
)

func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func EscapeString(str string) string {
	// 定义需要转义的特殊字符
	specialChars := `"` + "`" + `\` + "\n" + "\r" + "\t"

	// 创建正则表达式
	re := regexp.MustCompile(`[` + regexp.QuoteMeta(specialChars) + `]`)

	// 使用 ReplaceAllStringFunc 替换特殊字符
	return re.ReplaceAllStringFunc(str, func(s string) string {
		return `\` + s
	})
}

func ToCamelCase(s string) string {
	if s == "" {
		return ""
	}

	// Split by underscore and convert to camel case
	parts := strings.Split(s, "_")
	result := ""

	for _, part := range parts {
		if part != "" {
			result += FirstUpper(part)
		}
	}

	return result
}
