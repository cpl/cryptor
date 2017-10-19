package main

import (
	"html/template"
	"net/http"
)

// Parse templates
var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {

	page := &Page{
		Sidebar: Sidebar{
			IP:   "192.168.0.103",
			Port: 9000,
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "base.html", page)
	})

	// Handle assets
	http.Handle("/css/", http.StripPrefix(
		"/css/", http.FileServer(http.Dir("templates/css"))))

	http.ListenAndServe(":9000", nil)
}

func exection(t template.Template) template.HTML {
	b := bytes.Buffer{}
    t.ExecuteTemplate(&b, params.Name, nil)
	return template.HTML(b.String())
}

// Cache ...
type Cache struct {
	Path       string
	Size       int
	MaxSize    int
	Percent    float32
	ChunkCount int
}

// Sidebar ...
type Sidebar struct {
	IP   string
	Port int
}

// Page ...
type Page struct {
	Sidebar Sidebar
	Caches  []Cache
}
