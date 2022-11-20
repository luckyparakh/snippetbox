package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/showSnippet", app.showSnippet)
	mux.HandleFunc("/createSnippet", app.createSnippet)
	fileServer := http.FileServer(http.Dir("./ui/static"))
	// The http.FileServer serves the file by looking at the request URL.
	// Since our request URL is /static/css/main.css, it will try to look for a file inside root directory with the path /ui/static/static/css/main.css and it doesn’t exist.
	// we need to remove /static part from the URL
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
