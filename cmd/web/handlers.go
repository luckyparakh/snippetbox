package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/luckyparakh/snippetbox/pkg/models"
)

func (a *application) home(w http.ResponseWriter, r *http.Request) {
	// To restrict default catch all behaviour of Servemux
	if r.URL.Path != "/" {
		a.notFound(w)
		return
	}
	snippets, err := a.snippets.Latest()
	if err != nil {
		a.serveError(w, err)
		return
	}
	data := &templateData{Snippets: snippets}
	a.render(w, r, "home.page.tmpl", data)
}
func (a *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("allow", http.MethodPost)
		// following function should be called once, calling it multiple time will give warning
		w.WriteHeader(405)
		w.Write([]byte("method not allowed"))
		return
	}
	title := "First Post"
	content := "Enjoying this work"
	id, err := a.snippets.Insert(title, content, "7")
	if err != nil {
		a.serveError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/showSnippet?id=%d", id), http.StatusSeeOther)
}
func (a *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		a.notFound(w)
		return
	}
	s, err := a.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			a.notFound(w)
		} else {
			a.serveError(w, err)
		}
		return
	}
	data := &templateData{
		Snippet: s,
	}
	a.render(w, r, "show.page.tmpl", data)
}
