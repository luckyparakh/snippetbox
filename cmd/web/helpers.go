package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (a *application) serveError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.errLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (a *application) notFound(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}

func (a *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := a.templateCache[name]
	if !ok {
		a.serveError(w, fmt.Errorf("the template %s not found", name))
		return
	}
	buf := new(bytes.Buffer)
	if err := ts.Execute(buf, a.addDeafultData(td, r)); err != nil {
		a.serveError(w, err)
	}
	buf.WriteTo(w)
}

func (a *application) addDeafultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		return &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	return td
}
