package main

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseGlob("web/templates/*.html"))

func handleRoot(w http.ResponseWriter, r *http.Request) {

	localPage := page{
		Sidebar: staticSidebar,
		Blocks:  staticBlocks,
		Header: header{
			Title: "Cache info",
			Icon:  "database",
		},
	}

	templates.ExecuteTemplate(w, "base", localPage)
}

func main() {

	// Handle pages
	http.HandleFunc("/", handleRoot)

	// Handle assets
	http.Handle("/assets/", http.StripPrefix(
		"/assets/", http.FileServer(http.Dir("web/assets/"))))

	http.ListenAndServe(":8080", nil)
}
