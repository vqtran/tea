package engines

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
	file, err := readFile(filepath)
	if err != nil {
		return nil, err
	}
	included, err := insertIncludes(file, filepath)
	if err != nil {
		return nil, err
	}
	return template.New(filepath).Parse(included)
}

var Html = html{}

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

func insertIncludes(s string, basepath string) (string, error) {
	includes := getIncludes(s)
	removeIncludes := regexp.MustCompile(`{{[\s]*include(.+?)}}`)
	inserted := removeIncludes.ReplaceAllString(s, "")
	for _, filepath := range includes {
		fullpath := path.Join(path.Dir(basepath), filepath)
		included, err := readFile(fullpath)
		if err != nil {
			return "", err
		}
		inserted += included
	}
	return inserted, nil
}

func getIncludes(s string) []string {
	match := regexp.MustCompile(`{{.+?}}`)
	removeBrackets := regexp.MustCompile(`{|}`)
	removeQuotes := regexp.MustCompile(`"|'`)

	cleaned := make([]string, 0, 5)
	for _, str := range match.FindAllString(s, -1) {
		str = removeBrackets.ReplaceAllString(str, "")
		filepaths := strings.Fields(str)
		if filepaths[0] == "include" {
			for i, filepath := range filepaths {
				if i != 0 {
					filepath = removeQuotes.ReplaceAllString(filepath, "")
					cleaned = append(cleaned, filepath)
				}
			}
		}
	}
	return cleaned
}

