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

func PathToTemplate(filename string) string {
	return filepath.Join(basePath, "templates", filename)
}

var TemplateParseErr = errors.New("template error")

func ParseHTMLTemplate(templateName string, data any) ([]byte, error) {
	templateFile := PathToTemplate(templateName + ".html")
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {

		return nil, errors.Join(TemplateParseErr, err)
	}
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {

		return nil, errors.Join(TemplateParseErr, err)
	}

	return buffer.Bytes(), nil
}
