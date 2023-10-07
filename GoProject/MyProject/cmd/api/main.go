package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
	game   *Game // Define your game struct here
}

// Game represents your educational board game.
type Game struct {
	// Define game-specific fields and logic here
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	game := &Game{}

	app := &application{
		config: cfg,
		logger: logger,
		game:   game,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/start-game", app.StartGameHandler)
	mux.HandleFunc("/answer-question", app.AnswerQuestionHandler)
	mux.HandleFunc("/check-progress", app.CheckProgressHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}

func (app *application) StartGameHandler(w http.ResponseWriter, r *http.Request) {
}

func (app *application) AnswerQuestionHandler(w http.ResponseWriter, r *http.Request) {
}

func (app *application) CheckProgressHandler(w http.ResponseWriter, r *http.Request) {
}
