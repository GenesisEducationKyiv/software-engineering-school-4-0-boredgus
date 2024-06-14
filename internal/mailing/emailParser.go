package mailing

import (
	"bytes"
	"html/template"
	"path/filepath"
	"runtime"
	"subscription-api/config"
)

var basePath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basePath = filepath.Dir(currentFile)
}

func path(filename string) string {
	return filepath.Join(basePath, "emails", filename)
}

type HTMLTemplateParser interface {
	Parse(templateName string, data any) ([]byte, error)
}

type htmlTemplateParser struct {
	l config.Logger
}

func NewHTMLTemplateParser(l config.Logger) *htmlTemplateParser {
	return &htmlTemplateParser{l: l}
}
func (p htmlTemplateParser) Parse(templateName string, data any) ([]byte, error) {
	var buffer bytes.Buffer
	filename := path(templateName + ".html")
	if err := template.
		Must(template.ParseFiles(filename)).
		Execute(&buffer, data); err != nil {
		p.l.Errorf("failed to execute html template %s: %v", filename, err)

		return []byte{}, err
	}

	return buffer.Bytes(), nil
}
