package tea

import (
	"errors"
	"github.com/vqtran/tea/engines"
	"html/template"
	"os"
	"path"
	"sync"
)

type Engine interface {
	CompileFile(filepath string) (*template.Template, error)
}

type Options struct {
	FileExt string
	Recursive bool
}

var (
	engine Engine = engines.Html
	compiled map[string]*template.Template
	opt = Options{".html", true}
	mutex sync.RWMutex
)

func SetEngine(en string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if en == "html" {
		engine = engines.Html
		opt.FileExt = ".html"
	} else if en == "amber" {
		engine = engines.Amber
		opt.FileExt = ".amber"
	} else {
		return errors.New("Tea Error: Engine not supported.")
	}

	return nil
}

func GetEngine() *Engine {
	mutex.RLock()
	defer mutex.RUnlock()
	return &engine
}

func Compile(dirpath string, options Options) error {
	mutex.Lock()
	defer mutex.Unlock()

	var err error
	compiled, err = compile(dirpath, options)
	return err
}

func compile(dirpath string, opt Options) (map[string]*template.Template, error) {
	dir, err := os.Open(dirpath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	cmpl := make(map[string]*template.Template)
	for _, file := range files {
		// filename is for example "index.amber"
		filename := file.Name()
		fileext := path.Ext(filename)

		// If recursive is true and there's a subdirectory, recurse
		if opt.Recursive && file.IsDir() {
			subdirpath := path.Join(dirpath, filename)
			subcompiled, err := compile(subdirpath, opt)
			if err != nil {
				return nil, err
			}
			// Copy templates from subdirectory into parent template mapping
			for k, v := range subcompiled {
				// Concat with parent directory name for unique paths
				key := path.Join(filename, k)
				cmpl[key] = v
			}
		} else if fileext == opt.FileExt {
			// Otherwise compile the file and add to mapping
			fullpath := path.Join(dirpath, filename)
			tmpl, err := engine.CompileFile(fullpath)
			if err != nil {
				return nil, err
			}
			// Strip extension
			key := filename[0:len(filename)-len(fileext)]
			cmpl[key] = tmpl
		}
	}

	return cmpl, nil
}

func Get(key string) (*template.Template, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	t , ok := compiled[key]
	return t, ok
}

func Delete(key string) {
	mutex.Lock()
	delete(compiled, key)
   mutex.Unlock()
}

func Clear() {
	mutex.Lock()
	compiled = make(map[string]*template.Template)
	mutex.Unlock()
}
