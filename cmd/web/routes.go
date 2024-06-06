package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.form))
	mux.Post("/addtask", dynamicMiddleware.ThenFunc(app.addtask))
	mux.Post("/deletetask", dynamicMiddleware.ThenFunc(app.deletetask))
	mux.Post("/update", dynamicMiddleware.ThenFunc(app.update))

	fileServer := http.FileServer(http.Dir(app.config.StaticDir))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
