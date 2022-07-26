package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func Build(configPath string, sourceDir string, relativeDir string, targetDir string) {
	entries, err := os.ReadDir(sourceDir)
	Check(err)
	for _, entry := range entries {
		info, err := entry.Info()
		name := info.Name()
		path := filepath.Join(sourceDir, name)
		targetFolder := filepath.Join(targetDir, relativeDir)
		targetPath := filepath.Join(targetFolder, name)
		Check(err)
		if !entry.IsDir() {
			if filepath.Ext(info.Name()) == ".xml" {
				Check(err)
				content := replaceVariable(configPath, path)
				os.MkdirAll(targetFolder, os.ModePerm)
				log.Printf("Compiled %s", targetPath)
				err := os.WriteFile(targetPath, content, os.ModePerm)
				Check(err)
			}
		} else {
			_relativeDir := filepath.Join(relativeDir, name)
			Build(configPath, path, _relativeDir, targetDir)
		}
	}
}

func replaceVariable(configPath string, path string) []byte {
	content, err := os.ReadFile(path)
	policy := string(content)
	Check(err)

	variables := GetRequestedVariables(policy)
	for _, _var := range variables {
		val, err := GetVariable(_var, configPath)
		Check(err)
		if val == "" || val == "null" {
			log.Fatalf("Variable %s is not provided in the config file. File: %s\n", _var, path)
		}
		regexpstring := fmt.Sprintf("(?s)({{\\s*)%s?(\\s*}})", _var)
		re := regexp.MustCompile(regexpstring)
		policy = re.ReplaceAllString(policy, val)
	}
	return []byte(policy)
}
