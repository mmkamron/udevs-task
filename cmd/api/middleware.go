package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/mmkamron/miniTwitter/internal/pkg/utils"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		userToken, err := utils.VerifyToken(app.config, tokenString.Value)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		sub, err := userToken.Claims.GetSubject()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		id, _ := strconv.ParseInt(sub, 10, 64)

		ctx := context.WithValue(context.TODO(), "userID", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
