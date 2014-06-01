package main

import (
	"github.com/vqtran/tea"
	"net/http"
)

func main() {
	tea.MustCompile("templates/html", tea.Options{".html", true})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, ok := tea.Get("index")
		if ok {
			data := map[string]string{"Name" : "World"}
			tmpl.Execute(w, data)
		}
	})

	http.HandleFunc("/layout", func(w http.ResponseWriter, r *http.Request) {
		tmpl, ok := tea.Get("index")
		if ok {
			tmpl.Execute(w, nil)
		}
	})

	http.ListenAndServe(":8080", nil)
}
