package main

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/signup", app.signUp)
	mux.HandleFunc("POST /v1/signin", app.signIn)
	mux.HandleFunc("GET /v1/logout", app.logout)

	mux.HandleFunc("POST /v1/tweets", app.authMiddleware(app.postTweet))

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	return app.recoverPanic(mux)
}
