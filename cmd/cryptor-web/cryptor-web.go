package main

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseGlob("web/templates/*.html"))

func redirectRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/overview", 301)
}

func handleCaches(w http.ResponseWriter, r *http.Request) {

	localPage := page{
		Sidebar: staticSidebar,
		Blocks:  cachesBlocks,
		Header: header{
			Title: "Cache info",
			Icon:  "database",
		},
	}

	templates.ExecuteTemplate(w, "base", localPage)
}

func handleOverview(w http.ResponseWriter, r *http.Request) {

	localPage := page{
		Sidebar: staticSidebar,
		Blocks:  overviewBlocks,
		Header: header{
			Title: "Statistics",
			Icon:  "pie-chart",
		},
	}

	templates.ExecuteTemplate(w, "base", localPage)
}

func main() {

	// Handle pages
	http.HandleFunc("/", redirectRoot)
	http.HandleFunc("/overview", handleOverview)
	http.HandleFunc("/caches", handleCaches)

	// Handle assets
	http.Handle("/assets/", http.StripPrefix(
		"/assets/", http.FileServer(http.Dir("web/assets/"))))

	http.ListenAndServe(":8080", nil)
}
