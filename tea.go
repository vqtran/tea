package tea

/**
	Tea - Making Go templates easy.
	By Vinh Tran (vqtran)

	A package that compiles directories of templates for various templating
	engines and provides a thread-safe cache and simple usage.
**/

import (
	"errors"
	"fmt"
	"github.com/vqtran/tea/engines"
	"io"
	"os"
	"path"
	"strings"
	"sync"
)

// Different templating engines must be support via a plugin in the engines
// subdirectory.
type Engine interface {
	// Has to be able to compile a single file from its path
	CompileFile(filepath string) (interface{}, error)
	// Function to render templates (does the casting)
	Render(buf io.Writer, tmpl interface{}, data interface{}) error
}

type Options struct {
	// File extension to match when compiling directories
	FileExt string
	// Whether search is recursive or only top-level
	Recursive bool
}

var (
	engine   Engine = engines.Html
	compiled map[string]interface{}
	opt      = Options{".html", true}
	mutex    sync.RWMutex
)

// Set what engine to use during compilation with a string.
// Also sets a default extension
func SetEngine(en string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if en == "html" {
		engine = engines.Html
		opt.FileExt = ".html"
	} else if en == "amber" {
		engine = engines.Amber
		opt.FileExt = ".amber"
	} else if en == "mustache" {
		engine = engines.Mustache
		opt.FileExt = ".html.mustache"
	} else {
		return errors.New("Tea Error: Engine not supported.")
	}

	return nil
}

// Read-locked way to get the underlying tea-plugin for the engine.
// Makes it easy if user wants to compile a single file.
func GetEngine() *Engine {
	mutex.RLock()
	defer mutex.RUnlock()
	return &engine
}

// Compile an entire directory with options.
func Compile(dirpath string, options Options) error {
	// compile to temp first to avoid unecessary locking
	c, err := compile(dirpath, options)
	if err == nil {
		mutex.Lock()
		compiled = c
		mutex.Unlock()
	}
	return err
}

// Same as Compile except will panic if there is an error
func MustCompile(dirpath string, options Options) {
	err := Compile(dirpath, options)
	if err != nil {
		panic(fmt.Sprintf("Tea Error: %s", err.Error()))
	}
}

// Private function used to walk and compile directories
// Searches directory path for matching extensions and compiles them
// using the engine, and stores in map.
// The key stored in map is the path from the original dirpath to the template,
// with the extension stripped.
func compile(dirpath string, opt Options) (map[string]interface{}, error) {
	dir, err := os.Open(dirpath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	cmpl := make(map[string]interface{})
	for _, file := range files {
		// filename is for example "index.amber"
		filename := file.Name()
		fileext := ext(filename)

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
			key := filename[0 : len(filename)-len(fileext)]
			cmpl[key] = tmpl
		}
	}

	return cmpl, nil
}

// Used over path.Ext because this takes everything after the first period.
func ext(filename string) string {
	splitted := strings.SplitN(filename, ".", 2)
	return "." + splitted[len(splitted)-1]
}

// Writes the compiled template with data to the buffer.
func Render(buf io.Writer, key string, data interface{}) error {
	tmpl, ok := Get(key)
	if ok {
		// Delegate to the engine's implementation.
		return engine.Render(buf, tmpl, data)
	}
	return errors.New("Tea Error: Could not render template " + key)
}

// Thread-safe get the template. Second return variable states whether
// or not key is in the map.
func Get(key string) (interface{}, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	t, ok := compiled[key]
	return t, ok
}

// Delete a specific key/value from the map.
func Delete(key string) {
	mutex.Lock()
	delete(compiled, key)
	mutex.Unlock()
}

// Clear the entire cache of templates.
func Clear() {
	mutex.Lock()
	compiled = make(map[string]interface{})
	mutex.Unlock()
}

func GetCache() map[string]interface{} {
	mutex.RLock()
	defer mutex.RUnlock()
	return compiled
}
