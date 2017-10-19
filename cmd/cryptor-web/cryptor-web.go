package main

import (
	"bytes"
	"html/template"
	"net/http"
)

// Parse templates
var templates = template.Must(template.ParseGlob("templates/*.html"))

func redirectRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/overview", 301)
}

func main() {

	// Handle pages
	http.HandleFunc("/", redirectRoot)
	http.HandleFunc("/overview", handleOverview)
	http.HandleFunc("/caches", handleCaches)

	// Handle assets
	http.Handle("/css/", http.StripPrefix(
		"/css/", http.FileServer(http.Dir("templates/css"))))

	http.ListenAndServe(":9000", nil)
}

func execute(t *template.Template, name string, data interface{}) template.HTML {
	b := bytes.Buffer{}
	t.ExecuteTemplate(&b, name, data)
	return template.HTML(b.String())
}

type sidebar struct {
	IP   string
	Port int
}

type content struct {
	Data template.HTML
}

type page struct {
	Content content
	Sidebar sidebar
}

func getSidebar() sidebar {
	return sidebar{
		IP:   "192.168.0.103",
		Port: 9000,
	}
}

func handleOverview(w http.ResponseWriter, r *http.Request) {
	overview := overviewContent{
		Caches: cacheContent{},
	}

	localPage := page{
		Content: content{
			Data: execute(templates, "overview.html", overview),
		},
		Sidebar: getSidebar(),
	}

	templates.ExecuteTemplate(w, "base.html", localPage)
}

type cache struct {
	Path       string
	Size       int
	MaxSize    int
	Percent    float32
	ChunkCount int
}

type cacheContent struct {
	CacheCount      int
	Caches          []cache
	TotalSize       int
	TotalMaxSize    int
	TotalPercent    float32
	TotalChunkCount int
}

type overviewContent struct {
	Caches    cacheContent
	Uptime    int
	PeerCount int
}

func (c *cacheContent) update() {
	c.CacheCount = len(c.Caches)
	for _, cache := range c.Caches {
		c.TotalSize += cache.Size
		c.TotalMaxSize += cache.MaxSize
		c.TotalChunkCount += cache.ChunkCount
	}
}

func handleCaches(w http.ResponseWriter, r *http.Request) {
	caches := cacheContent{}
	caches.update()

	localPage := page{
		Content: content{
			Data: execute(templates, "caches.html", caches),
		},
		Sidebar: getSidebar(),
	}

	templates.ExecuteTemplate(w, "base.html", localPage)
}
