package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte

	//Request Body
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	//Don't allow multiple json bodies
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("Body must have only a single JSON value")
	}

	return nil
}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	//adding json to out
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	//if there are headers
	if len(headers) > 0 {
		for key, value := range headers[0] {
			// type Header map[string][]string
			// A Header represents the key-value pairs in an HTTP header.
			// The keys should be in canonical form, as returned by CanonicalHeaderKey.
			w.Header()[key] = value
		}
	}
	//Setting content-type to header via responsewriter
	w.Header().Set("Content-Type", "application/json")
	//I.E. http.StatusOK = 200
	//call before w.Write()
	//Writing status code to connection
	w.WriteHeader(status)

	//Writing json to the connection via responsewriter
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	//If the status is present set status code to status[0]
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	//writing statuscode and error to connection/header
	return app.writeJSON(w, statusCode, payload)
}
