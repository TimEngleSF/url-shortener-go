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

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /ping", dynamic.ThenFunc(app.Ping))
	mux.Handle("POST /link/new", dynamic.ThenFunc(app.LinkPost))

	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle(("POST /user/logout"), protected.ThenFunc(app.userLogoutPost))

	mux.Handle("GET /", dynamic.ThenFunc(app.LinkRedirect))

	return standard.Then(mux)
}
