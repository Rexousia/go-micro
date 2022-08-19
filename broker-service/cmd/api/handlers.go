package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
	// out, _ := json.MarshalIndent(payload, "", "\t")
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusAccepted)
	// w.Write(out)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	//reading the JSON inside of requestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	//if the value at Action is...
	switch requestPayload.Action {
	case "auth":
		//authenticating the credentials from the client
		app.authenticate(w, requestPayload.Auth)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	//create some json & send off to auth mircorservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	//calling service
	//name of folder
	//url you want to hit
	//creating a new request, buffering the bytes into json, post the json to the authentication service
	//POST  HTTP/1.1 1 1 map[]{}
	request, err := http.NewRequest("POST", "http://authentication-service:80/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	//creating a value of type &http.Client{}
	client := &http.Client{}
	//returning a response to the client.
	response, err := client.Do(request)
	// fmt.Println("Response:", response)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	//close at the end
	defer response.Body.Close()
	// fmt.Println(response.Body)
	//checking status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	//decode json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	//writing the status to the header and passing data
	app.writeJSON(w, http.StatusAccepted, payload)
}
