package goapimaker

import "strings"

func ControllerImportWriter(sb *strings.Builder, projectName string, method string) {
	sb.WriteString("package controllers\n\n")
	sb.WriteString("import (\n")
	switch method {
	case "PUT":
		sb.WriteString("\t\"log\"\n")
	case "POST":
		sb.WriteString("\t\"fmt\"\n")
	}
	sb.WriteString("\t\"encoding/json\"\n")
	sb.WriteString("\t\"net/http\"\n")
	if method != "DELETE" {
		sb.WriteString("\t\"" + strings.ReplaceAll(projectName, " ", "_") + "/src/models\"\n")
	}
	sb.WriteString("\t\"" + strings.ReplaceAll(projectName, " ", "_") + "/src/config\"\n")
	sb.WriteString("\t_ \"github.com/lib/pq\"\n")
	sb.WriteString(")\n\n")
}
