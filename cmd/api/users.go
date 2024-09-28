package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mmkamron/miniTwitter/internal/data"
	"github.com/mmkamron/miniTwitter/internal/pkg/utils"
)

// @Summary Register a new user
// @Description Create a new user account
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param user body data.User true "User Info"
// @Success 202 {string} string "Created"
// @Failure 422 {string} string "Error"
// @Router /v1/signup [post]
func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		Username string
		Password string
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "could not process your request", http.StatusBadRequest)
		return
	}

	user := &data.User{
		Name:     input.Name,
		Username: input.Username,
	}

	if err := user.Password.Set(input.Password); err != nil {
		http.Error(w, "could not process your request", http.StatusBadRequest)
	}

	err := app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateUsername):
			http.Error(w, "a user with this username already exists", http.StatusUnprocessableEntity)
			return
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.wg.Add(1)

	err = app.writeJSON(w, http.StatusAccepted, user, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// signIn handles user login.
// @Summary Sign in a user
// @Description Authenticates a user with username and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body object{username=string, password=string} true "Login credentials"
// @Success 200 {string} string "Authentication successful"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 401 {object} string "Invalid credentials"
// @Failure 500 {object} string "Internal server error"
// @Router /v1/signin [post]
func (app *application) signIn(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "could not process your request", http.StatusBadRequest)
		return
	}

	user, err := app.models.Users.GetByUsername(input.Username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	tokenString, err := utils.CreateToken(app.config, user.ID)
	if err != nil {
		w.Write([]byte(err.Error()))
		log.Println("creattoken error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		MaxAge:   3600,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	})

	w.Write([]byte(tokenString))
}

// @Summary Logs out the user
// @Description Clears the authentication token by setting an expired cookie
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param token header string true "JWT token required for authentication"
// @Success 200 {string} string "you are logged out"
// @Router /v1/logout [get]
func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	})

	w.Write([]byte("you are logged out"))
}
