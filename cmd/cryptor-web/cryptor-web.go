package main

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseGlob("web/templates/*.html"))

func handleRoot(w http.ResponseWriter, r *http.Request) {

	localPage := page{
		Sidebar: sidebar{
			IP:    "localhost",
			Port:  9000,
			Color: "33aa33",
			Sections: []section{
				section{Name: "Dashboard", SubSections: []subsection{
					subsection{
						Name:   "Overview",
						Icon:   "eye",
						Link:   "overview",
						Active: true},
					subsection{
						Name:   "Caches",
						Icon:   "database",
						Link:   "caches",
						Active: false},
					subsection{
						Name:   "Chunks",
						Icon:   "cubes",
						Link:   "chunks",
						Active: false},
					subsection{
						Name:   "Peers",
						Icon:   "users",
						Link:   "peers",
						Active: false},
					subsection{
						Name:   "Settings",
						Icon:   "cog",
						Link:   "settings",
						Active: false},
				}},
				section{Name: "Actions", SubSections: []subsection{
					subsection{
						Name:   "Manage packages",
						Icon:   "archive",
						Link:   "request",
						Active: false},
				}},
			},
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

type page struct {
	Sidebar sidebar
}

type sidebar struct {
	IP       string
	Port     int
	Color    string
	Sections []section
}

type section struct {
	Name        string
	SubSections []subsection
}

type subsection struct {
	Name   string
	Icon   string
	Link   string
	Active bool
}
