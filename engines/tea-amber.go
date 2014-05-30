package engines

import (
	"html/template"
)

type amber struct {}

func (a amber) CompileFile(filepath string) (*template.Template, error) {
	return nil, nil
}

var Amber = amber{}


