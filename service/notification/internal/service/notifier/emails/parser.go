package emails

import (
	"bytes"
	"errors"
	"html/template"
	"path/filepath"
	"runtime"
)

var basePath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basePath = filepath.Dir(currentFile)
}

func pathToTemplate(filename string) string {
	return filepath.Join(basePath, "templates", filename)
}

var HTMLTemplateErr = errors.New("html template error")

func ParseHTMLTemplate(templateName string, data any) ([]byte, error) {
	templateFile := pathToTemplate(templateName + ".html")
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return nil, errors.Join(HTMLTemplateErr, err)
	}
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		return nil, errors.Join(HTMLTemplateErr, err)
	}

	return buffer.Bytes(), nil
}
