package util

import (
	"bytes"
	"fmt"
	"text/template"
)

func LoadTemplate(templateContent, templateName string, dataMapContainer interface{}) (bool, []byte) {
	tmpl, err := template.New(templateName).Parse(templateContent)
	var outputBytes bytes.Buffer
	err = tmpl.Execute(&outputBytes, dataMapContainer)
	if err != nil {
		fmt.Printf("Error in loading template %s %v\n", templateName, err)
		return false, nil
	}
	return true, outputBytes.Bytes()
}
