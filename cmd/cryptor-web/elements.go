package main

// TODO: Add functions to store/load JSON

var staticSidebar = sidebar{
	IP:    "localhost",
	Port:  9000,
	Color: "33aa33",
	Sections: []section{
		section{Name: "Dashboard", SubSections: []subsection{
			subsection{
				Name:   "Overview",
				Icon:   "eye",
				Link:   "overview",
				Active: false},
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
}

var cachesBlocks = []block{
	block{
		Name:  "Caches",
		Value: "3",
		Color: "teal",
		Icon:  "database",
	},
	block{
		Name:  "Chunks",
		Value: "482",
		Color: "red",
		Icon:  "cubes",
	},
	block{
		Name:  "Size",
		Value: "1055 MB",
		Color: "green",
		Icon:  "th",
	},
	block{
		Name:  "Free",
		Value: "145 MB",
		Color: "blue",
		Icon:  "th-large",
	},
}

var overviewBlocks = []block{
	// Row 0
	block{
		Name:  "Caches",
		Value: "3",
		Color: "teal",
		Icon:  "database",
	},
	block{
		Name:  "Chunks",
		Value: "482",
		Color: "red",
		Icon:  "cubes",
	},
	block{
		Name:  "Size",
		Value: "1055 MB",
		Color: "green",
		Icon:  "th",
	},
	block{
		Name:  "Free",
		Value: "145 MB",
		Color: "red",
		Icon:  "th-large",
	},
	// Row 1
	block{
		Name:  "Peers",
		Value: "20",
		Color: "teal",
		Icon:  "users",
	},
	block{
		Name:  "UP",
		Value: "2.5 GB",
		Color: "blue",
		Icon:  "upload",
	},
	block{
		Name:  "DOWN",
		Value: "1.2 GB",
		Color: "grey",
		Icon:  "download",
	},
	block{
		Name:  "Active Requests",
		Value: "2",
		Color: "blue",
		Icon:  "tasks",
	},
}
