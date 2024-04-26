package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /ping", app.Ping)
	mux.HandleFunc("POST /link/new", app.LinkPost)

	mux.HandleFunc("GET /signup", app.SignUpForm)
	mux.HandleFunc("POST /user/add", app.SignUpPost)
	mux.HandleFunc("GET /login", app.LoginForm)

	mux.HandleFunc("GET /", app.LinkRedirect)

	return app.logRequest(mux)
}
