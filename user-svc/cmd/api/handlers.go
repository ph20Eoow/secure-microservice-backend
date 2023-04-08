package main

import (
	"errors"
	"net/http"
)

func (app *Config) authBasic(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// query user by email
	user, err := app.Models.User.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		return
	}
	// validate password
	app.writeJSON(w, http.StatusAccepted, user)
}

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
