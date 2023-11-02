package main

import (
	"EBG.IssataySheg.net/internal/data"
	"EBG.IssataySheg.net/internal/validator"
	"fmt"
	"net/http"
	"time"
)

func (app *application) createGameHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string     `json:"title"`
		Score data.Score `json:"score"`
		Games []string   `json:"games"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	movie := &data.Game{
		Title: input.Title,
		Score: input.Score,
		Games: input.Games,
	}
	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
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
	err = app.writeJSON(w, http.StatusOK, envelope{"game": game}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
