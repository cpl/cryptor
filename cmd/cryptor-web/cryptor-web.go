package main

import (
	"html/template"
	"net/http"
)

func main() {
	templates := template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "base.html", "")
	})

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("templates/css"))))

	http.ListenAndServe(":9000", nil)
}
