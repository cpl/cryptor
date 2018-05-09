package main

type page struct {
	Sidebar sidebar
	Blocks  []block
	Header  header
}

type sidebar struct {
	IP    string
	Port  int
	Color string

	// TODO: Replace side.Sections[]..Subsections[] with hash maps
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

type block struct {
	Name  string
	Icon  string
	Value string
	Color string
}

type header struct {
	Title string
	Icon  string
}
