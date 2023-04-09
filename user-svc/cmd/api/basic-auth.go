package main

import (
	"context"
	"errors"
	"net/http"
)

// middleware for handling
func (app *Config) isAuthByBasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			app.errorJSON(w, errors.New("Unauthorized, Invalid Header"), http.StatusUnauthorized)
			return
		}
		// query user by email
		user, err := app.Models.User.GetUserByEmail(username)
		if err != nil {
			app.errorJSON(w, errors.New("Unauthorized"), http.StatusBadRequest)
			return
		}
		// validate password
		passMatched, err := user.PasswordMatched(password)
		if err != nil {
			// masking error message: "crypto/bcrypt: hashedPassword is not the hash of the given password"
			app.errorJSON(w, errors.New("Unauthorized"), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		if passMatched {
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
