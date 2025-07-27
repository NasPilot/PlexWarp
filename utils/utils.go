package utils

import (
	"strings"
)

// Center 居中对齐文本
func Center(text string, width int, fillChar string) string {
	textLen := len(text)
	if textLen >= width {
		return text
	}

	padding := width - textLen
	leftPadding := padding / 2
	rightPadding := padding - leftPadding

	return strings.Repeat(fillChar, leftPadding) + text + strings.Repeat(fillChar, rightPadding)
}

// Contains 检查字符串切片是否包含指定字符串
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveEmpty 移除字符串切片中的空字符串
func RemoveEmpty(slice []string) []string {
	var result []string
	for _, s := range slice {
		if strings.TrimSpace(s) != "" {
			result = append(result, s)
		}
	}
	return result
}