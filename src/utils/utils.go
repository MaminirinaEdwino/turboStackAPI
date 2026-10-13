package utils

import (
	"fmt"
	"strings"
)

func WriteErrorCheker(message string) string {
	return fmt.Sprintf("if err != nil {\nfmt.Println(\"%s\", err)\n}", message)
}

func ErrorChecker(err error) {
	if err != nil {
		panic(err)
	}
}

func ToUpperFirstLetter(value string) string {
	tmpTab := []string{
		strings.ToUpper(strings.Split(value, "")[0]),
		
	}
	for _, val := range strings.Split(value, "")[1:]{
		tmpTab = append(tmpTab, val)
	}
	return strings.Join(tmpTab, "")
}