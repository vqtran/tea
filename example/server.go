package main

import (
	"github.com/vqtran/tea"
	"net/http"
	"log"
)

func main() {
	// Mess around with these settings, and see that all you have to change
	// is what engine you're using and where to look!
	tea.SetEngine("amber")
	tea.MustCompile("templates/amber", tea.Options{".amber", true})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{"Name" : "World"}
		err := tea.Render(w, "index", data)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/layout", func(w http.ResponseWriter, r *http.Request) {
		tea.Render(w, "layouts/layout", nil)
	})

	http.ListenAndServe(":8080", nil)
}
