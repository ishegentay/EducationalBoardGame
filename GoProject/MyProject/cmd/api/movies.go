package main

import (
	"EBG.IssataySheg.net/internal/data"
	"fmt"
	"net/http"
	"time"
)

func (app *application) createGameHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Created a new educational board game")
}

func (app *application) showGameInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	game := data.Game{
		ID:          id,
		CreatedAt:   time.Now(),
		Title:       "Math Challenge",
		Description: "Challenge yourself in math!!!",
		Score:       0,
		Version:     1,
	}

	err = app.writeJSON(w, http.StatusOK, game, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
