package mailing

import (
	"path/filepath"
	"runtime"
)

var basePath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basePath = filepath.Dir(currentFile)
}

func PathToTemplate(filename string) string {
	return filepath.Join(basePath, "emails", filename)
}
