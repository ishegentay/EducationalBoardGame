package main

import (
	"EBG.IssataySheg.net/internal/data"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (app *application) createGameHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be in the
	// HTTP request body (note that the field names and types in the struct are a subset
	// of the Movie struct that we created earlier). This struct will be our *target
	// decode destination*.
	var input struct {
		Title string   `json:"title"`
		Score int32    `json:"score"`
		Games []string `json:"games"`
	}
	// Initialize a new json.Decoder instance which reads from the request body, and
	// then use the Decode() method to decode the body contents into the input struct.
	// Importantly, notice that when we call Decode() we pass a *pointer* to the input
	// struct as the target decode destination. If there was an error during decoding,
	// we also use our generic errorResponse() helper to send the client a 400 Bad
	// Request response containing the error message.
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	// Dump the contents of the input struct in a HTTP response.
	fmt.Fprintf(w, "%+v\n", input)
}
func (app *application) showGameHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	game := data.Game{
		ID:          id,
		CreatedAt:   time.Now(),
		Title:       "Educational Board Game",
		Description: "Learn for fun!",
		Games:       []string{"History Board Games", "Geography Board Games", "Science Board Games", "Math Board Games", "Language Learning Board Games"},
		Score:       58,
		Version:     1,
	}
	// Create an envelope{"movie": movie} instance and pass it to writeJSON(), instead
	// of passing the plain movie struct.
	err = app.writeJSON(w, http.StatusOK, envelope{"game": game}, nil)
	if err != nil {
		// Use the new serverErrorResponse() helper.
		app.serverErrorResponse(w, r, err)
	}
}
