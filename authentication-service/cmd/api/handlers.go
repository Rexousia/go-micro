package main

import (
	"errors"
	"fmt"
	"net/http"
)

//authenticate is a reciever function of type *Config
//Authenticate is a handlerfunction
func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	//creating a data strutucure to decode json into
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//calling app.ReadJSON because readJSON has a receiver of type *Config app. Is of type *Config
	//decoding the json being passed through the request
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// app is of type config
	// Models share the Config struct and is of type Config
	// User is a field inside of the Models struct
	// GetByEmail is a function that takes a receiver of type User
	// and passes a string
	// requestPayload. Email access the data in the Email field being passed.
	// validate the user against the database
	// return value of type *data.User
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}
	//passing the password in the requestPayload
	//checking if the entered password associated with the email is correct
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	//if it made it this far log in successful
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	//writing and converting data from payload to the connection via json status code 202
	app.writeJSON(w, http.StatusAccepted, payload)
}
