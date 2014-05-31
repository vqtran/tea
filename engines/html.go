package engines

import (
	"html/template"
)

type html struct{}

func (h html) CompileFile(filepath string) (*template.Template, error) {
	return nil, nil
}

var Html = html{}
