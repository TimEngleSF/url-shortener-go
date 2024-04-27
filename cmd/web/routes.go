package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /ping", dynamic.ThenFunc(app.Ping))
	mux.Handle("POST /link/new", dynamic.ThenFunc(app.LinkPost))

	mux.Handle("GET /signup", dynamic.ThenFunc(app.SignUpForm))
	mux.Handle("POST /user/add", dynamic.ThenFunc(app.SignUpPost))

	mux.Handle("GET /login", dynamic.ThenFunc(app.LoginForm))

	mux.Handle("GET /", dynamic.ThenFunc(app.LinkRedirect))

	return standard.Then(mux)
}
