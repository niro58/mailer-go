package service

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

var (
	ErrVariableLength   = errors.New("variable length mismatch")
	ErrVariableNotFound = errors.New("variable not found")
)

type Template struct {
	Name      string `json:"name"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	Variables []string
}

func (t Template) Validate(vars map[string]string) error {
	if len(vars) != len(t.Variables) {
		return ErrVariableLength
	}

	for _, v := range t.Variables {
		found := false
		for varr := range vars {
			if v == varr {
				found = true
				break
			}
		}
		if !found {
			return ErrVariableNotFound
		}
	}
	return nil
}
func (t Template) FindVariables() []string {
	vars := make(map[string]bool)
	for _, varsRaw := range append(strings.Split(t.Body, "{{"), strings.Split(t.Subject, "{{")...) {
		parts := strings.Split(varsRaw, "}}")
		if len(parts) == 2 {
			vars[parts[0]] = true
		}
	}

	var result []string
	for k := range vars {
		result = append(result, k)
	}
	return result
}
func getTemplates() (map[string]Template, error) {
	dirname, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	folder := path.Join(dirname, "/templates")
	folderFiles, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	result := make(map[string]Template)
	for _, file := range folderFiles {
		jsonFile, err := os.Open(path.Join(folder, file.Name()))
		if err != nil {
			return nil, err
		}
		defer jsonFile.Close()
		var template Template
		byteValue, _ := io.ReadAll(jsonFile)
		json.Unmarshal([]byte(byteValue), &template)
		template.Variables = template.FindVariables()

		result[template.Name] = template
	}

	return result, nil
}
