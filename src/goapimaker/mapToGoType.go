package goapimaker

import "strings"

func MapToGoType(uiType string) string {
	switch strings.ToLower(uiType) {
	case "int":
		return "int"
	case "string", "text":
		return "string"
	default:
		return "string"
	}
}