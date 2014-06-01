package engines

/**
	Tea engine plugin for html/template.
	Supports usual templates but also the "include" function.
	Essentially a macro in the form of {{ include "file.html" }}
	that replaces itself with the contentsof file.html (non-recursively)
**/

import (
	"bufio"
	"html/template"
	"os"
	"path"
	"regexp"
	"strings"
)

type html struct{}

func (h html) CompileFile(filepath string) (*template.Template, error) {
	// Read the original file into string
	file, err := readFile(filepath)
	if err != nil {
		return nil, err
	}
	// Replace the include statements with their file contents
	included, err := insertIncludes(file, filepath)
	if err != nil {
		return nil, err
	}
	// Parse the string into an html/template.
	return template.New(filepath).Parse(included)
}

var Html = html{}

// Reads a file into a string in memory
func readFile(path string) (string, error) {
  file, err := os.Open(path)
  if err != nil {
    return "", err
  }
  defer file.Close()

  contents := ""
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
	  text := scanner.Text()
	  if text != "" {
		  contents += text + "\n"
	  }
  }
  return contents, scanner.Err()
}

// Finds all the includes and replaces their include statement with the
// contents of the referenced file.
func insertIncludes(s string, basepath string) (string, error) {
	// Gets all the filepaths from the include statements
	includes := getIncludes(s)
	// Split on include statements
	parsed := regSplit(s, `{{[\s]*include(.+?)}}`)

	inserted := ""
	// Current index of include to insert
	includeCount := 0
	for _, fragment := range parsed {
		fragment = strings.TrimSpace(fragment)
		if fragment != "" {
			// Add original fragment first
			inserted += fragment + "\n"
		}
		// If there are includes left insert the next include's contents
		if includeCount < len(includes) {
			fullpath := path.Join(path.Dir(basepath), includes[includeCount])
			included, err := readFile(fullpath)
			if err != nil {
				return "", err
			}
			inserted += strings.TrimSpace(included) + "\n"
			includeCount++
		}
	}
	return inserted, nil
}

// Looks through an html file string and returns a list of filepaths
// that needs to be included
func getIncludes(s string) []string {
	// Regex to match include statements
	match := regexp.MustCompile(`{{[\s]*include(.+?)}}`)
	// Regex matching brackets
	removeBrackets := regexp.MustCompile(`{|}`)
	// Regex matching quotes
	removeQuotes := regexp.MustCompile(`"|'`)

	cleaned := make([]string, 0, 5)
	// For all matches include their filepaths
	for _, str := range match.FindAllString(s, -1) {
		str = removeBrackets.ReplaceAllString(str, "")
		filepaths := strings.Fields(str)
		if filepaths[0] == "include" {
			for i, filepath := range filepaths {
				// Skipping "include"
				if i != 0 {
					filepath = removeQuotes.ReplaceAllString(filepath, "")
					cleaned = append(cleaned, filepath)
				}
			}
		}
	}
	return cleaned
}

// Credit to Stack Overflow:
// http://stackoverflow.com/questions/4466091/split-string-using-regular-expression-in-go
func regSplit(text string, delimeter string) []string {
    reg := regexp.MustCompile(delimeter)
    indexes := reg.FindAllStringIndex(text, -1)
    laststart := 0
    result := make([]string, len(indexes) + 1)
    for i, element := range indexes {
            result[i] = text[laststart:element[0]]
            laststart = element[1]
    }
    result[len(indexes)] = text[laststart:len(text)]
    return result
}
