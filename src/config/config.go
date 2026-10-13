package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var UserHomeDir, _ = os.UserHomeDir()

// var TURBO_STACK_DIR = UserHomeDir + "/.turbo_stack"
var TURBO_STACK_DIR =filepath.Join(UserHomeDir, ".turbo_stack")
var PROJECT_DIR = filepath.Join(UserHomeDir, ".turbo_stack", "turbo_projects")
var LIBRAIRIE_PATH = TURBO_STACK_DIR + "/librairie.json"
var PAGE_LIB_DIR = TURBO_STACK_DIR + "/librairies/page"
var COMPONENT_LIB_DIR = TURBO_STACK_DIR + "/librairies/component"
var STYLE_LIB_DIR = TURBO_STACK_DIR + "/librairies/style"

func CheckIfExist(chemin string) bool {
	info, err := os.Stat(chemin)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println(err)
			return false
		}
		return false
	}
	return info.IsDir()
}

func CheckCreateDir(path string) {
	filePath := filepath.Join(PROJECT_DIR, path)
	if !CheckIfExist(filePath) {
		os.MkdirAll(filePath, os.ModePerm)
	}
}

func CheckEmptyFile(chemin string) (bool, error) {
	info, err := os.Stat(chemin)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	if info.IsDir() {
		return false, fmt.Errorf("%s est un dossier, pas un fichier", chemin)
	}
	return info.Size() == 0, nil
}
