package main

import (
	"html/template"
	"net/http"
)

type Page struct {
	Content template.Template
}

func main() {
	templates := template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "default.html", "")
	})

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("templates/css"))))

	http.ListenAndServe(":9000", nil)
}
