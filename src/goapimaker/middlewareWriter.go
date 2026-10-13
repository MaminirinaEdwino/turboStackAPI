package goapimaker

import (
	"fmt"
	"os"
	"strings"

	"github.com/MaminirinaEdwino/turbostackapi/src/config"
)

func MiddlewareWriter(projectName string) {
	filePath := fmt.Sprintf("%s/%s/api/src/middlewares/cors.go", config.PROJECT_DIR, projectName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating middleware file %s : %v\n", filePath, err)
		return
	}
	defer file.Close()

	var sb strings.Builder
	sb.WriteString("package middlewares\n\n")
	sb.WriteString("import \"net/http\"\n\n")
	sb.WriteString("// CorsMiddleware gère les autorisations Cross-Origin (CORS) pour autoriser toutes les sources\n")
	sb.WriteString("func CorsMiddleware(next http.Handler) http.Handler {\n")
	sb.WriteString("\treturn http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {\n")
	sb.WriteString("\t\tw.Header().Set(\"Access-Control-Allow-Origin\", \"*\")\n")
	sb.WriteString("\t\tw.Header().Set(\"Access-Control-Allow-Methods\", \"GET, POST, PUT, DELETE, OPTIONS\")\n")
	sb.WriteString("\t\tw.Header().Set(\"Access-Control-Allow-Headers\", \"Content-Type, Authorization\")\n\n")
	sb.WriteString("\t\tif r.Method == \"OPTIONS\" {\n")
	sb.WriteString("\t\t\tw.WriteHeader(http.StatusOK)\n")
	sb.WriteString("\t\t\treturn\n")
	sb.WriteString("\t\t}\n\n")
	sb.WriteString("\t\tnext.ServeHTTP(w, r)\n")
	sb.WriteString("\t})\n")
	sb.WriteString("}\n")

	file.WriteString(sb.String())
}
