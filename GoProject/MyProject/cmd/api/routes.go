package main

import (
	"EBG.IssataySheg.net/internal/data"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func (app *application) CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	game := data.Game{
		ID:          data.GenerateGameID(),
		CreatedAt:   time.Now(),
		Title:       "Math Challenge",
		Description: "Challenge yourself in math!!!",
		Score:       0,
		Version:     1,
	}

	err := app.writeJSON(w, http.StatusCreated, game, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

func (app *application) GetGameInfoHandler(w http.ResponseWriter, r *http.Request) {

	gameID := app.readGameIDParam(r)
	if gameID == "" {
		http.NotFound(w, r)
		return
	}

	game := data.Game{
		ID:          gameID,
		CreatedAt:   time.Now(),
		Title:       "Math Challenge",
		Description: "Challenge yourself in math!!!",
		Score:       100,
		Version:     1,
	}

	err := app.writeJSON(w, http.StatusOK, game, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

func (app *application) readGameIDParam(r *http.Request) string {
	return r.URL.Query().Get("gameID")
}

func (app *application) writeJSON(w http.ResponseWriter, statusCode int, data interface{}, headers map[string]string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = w.Write(jsonData)
	return err
}

func main() {
	app := &application{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}

	http.HandleFunc("/games/create", app.CreateGameHandler)
	http.HandleFunc("/games/info", app.GetGameInfoHandler)

	port := 4000 // Your desired port number
	serverAddr := fmt.Sprintf(":%d", port)

	app.logger.Printf("Starting the educational board game server on %s\n", serverAddr)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		app.logger.Fatal(err)
	}
}
