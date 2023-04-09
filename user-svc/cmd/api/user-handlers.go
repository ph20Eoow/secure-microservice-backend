package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ph20Eoow/auth-svc/data"
)

func (app *Config) createUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, errors.New("Invalid requests"), http.StatusBadRequest)
		return
	}
	// validate email format
	validEmail, err := app.Models.User.ValidateEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("Invalid email format"), http.StatusBadRequest)
		return
	}
	// validate password
	validPassword, err := app.Models.User.ValidatePassword(requestPayload.Password)
	if err != nil {
		app.errorJSON(w, errors.New("Fail to comply password complexity"), http.StatusBadRequest)
		return
	}
	// Passing validation
	if validEmail && validPassword {
		id, err := app.Models.User.InsertUser(requestPayload.Email, requestPayload.Password)
		if err != nil {
			app.errorJSON(w, err, http.StatusBadRequest)
		}

		type Data struct {
			UserId int `json:"userId"`
		}
		payload := jsonResponse{
			Error:   false,
			Message: "Created user",
			Data: Data{
				UserId: id,
			},
		}
		app.writeJSON(w, http.StatusAccepted, payload)
	}
}

/*
 * update user password
 */
func (app *Config) updatePassword(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*data.User)
	var requestPayload struct {
		OldPassword string `json:"oldPassword"`
		Password    string `json:"password"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, errors.New("Invalid request"), http.StatusBadRequest)
		return
	}
	// validate user password
	passwordMatched, err := user.PasswordMatched(requestPayload.OldPassword)
	if err != nil || !passwordMatched {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate password complexity
	validPassword, err := app.Models.User.ValidatePassword(requestPayload.Password)
	if err != nil || !validPassword {
		app.errorJSON(w, errors.New("Fail to comply password complexity"), http.StatusBadRequest)
		return
	}

	// update user password
	if validPassword {
		ok, err := user.UpdatePassword(requestPayload.OldPassword, requestPayload.Password)
		if err != nil {
			app.errorJSON(w, err, http.StatusBadRequest)
			return
		}
		var payload jsonResponse
		payload.Error = false
		payload.Message = "Success"
		payload.Data = nil
		if ok {
			app.writeJSON(w, http.StatusAccepted, &payload)
			return
		}
	}
}

func (app *Config) getUserProfile(w http.ResponseWriter, r *http.Request) {
	// query user by email
	if userId := chi.URLParam(r, "id"); userId != "" {
		user, err := app.Models.User.GetUserById(userId)
		if err != nil {
			app.errorJSON(w, errors.New("UserID not found"), http.StatusBadRequest)
			return
		}
		app.writeJSON(w, http.StatusAccepted, user)
	} else {
		app.errorJSON(w, errors.New("UserID not found"), http.StatusBadRequest)
		return
	}
}
