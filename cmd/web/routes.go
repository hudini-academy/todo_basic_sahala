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

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))
    mux.Get("/user/special", dynamicMiddleware.ThenFunc(app.specialForm))
	mux.Post("/user/special", dynamicMiddleware.ThenFunc(app.specialadd))
	mux.Get("/user/special", dynamicMiddleware.ThenFunc(app.specialForm))
	mux.Post("/user/special", dynamicMiddleware.ThenFunc(app.specialdelete))
	mux.Get("/user/special", dynamicMiddleware.ThenFunc(app.specialdelete))

	fileServer := http.FileServer(http.Dir(app.config.StaticDir))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
	
}
